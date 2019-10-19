package epub

import (
	"archive/zip"
	"crypto/sha1"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/geek1011/BookBrowser/booklist"
	"github.com/geek1011/BookBrowser/formats"

	"time"

	"github.com/beevik/etree"
	"github.com/moraes/isbn"
	"github.com/pkg/errors"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

type epub struct {
	hascover  bool
	book      *booklist.Book
	coverpath *string
}

func (e *epub) Book() *booklist.Book {
	return e.book
}

func (e *epub) HasCover() bool {
	return e.coverpath != nil
}

func (e *epub) GetCover() (i image.Image, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("panic while decoding cover image")
		}
	}()

	zr, err := zip.OpenReader(e.book.FilePath)
	if err != nil {
		return nil, errors.Wrap(err, "error opening epub as zip")
	}
	defer zr.Close()

	zfs := zipfs.New(zr, "epub")

	cr, err := zfs.Open(*e.coverpath)
	if err != nil {
		return nil, errors.Wrapf(err, "could not open cover '%s'", *e.coverpath)
	}
	defer cr.Close()

	i, _, err = image.Decode(cr)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding image")
	}

	return i, nil
}

func load(filename string) (formats.BookInfo, error) {
	e := &epub{book: &booklist.Book{}, hascover: false}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, errors.Wrapf(err, "could not stat book")
	}
	e.book.FilePath = filename
	e.book.FileSize = fi.Size()
	e.book.ModTime = fi.ModTime()

	s := sha1.New()
	i, err := io.Copy(s, f)
	if err == nil && i != fi.Size() {
		err = errors.New("could not read whole file")
	}
	if err != nil {
		f.Close()
		return nil, errors.Wrap(err, "could not hash book")
	}
	e.book.Hash = fmt.Sprintf("%x", s.Sum(nil))

	f.Close()

	zr, err := zip.OpenReader(filename)
	if err != nil {
		return nil, errors.Wrap(err, "error opening epub as zip")
	}
	defer zr.Close()

	zfs := zipfs.New(zr, "epub")

	rsk, err := zfs.Open("/META-INF/container.xml")
	if err != nil {
		return nil, errors.Wrap(err, "error reading container.xml")
	}
	defer rsk.Close()

	container := etree.NewDocument()
	_, err = container.ReadFrom(rsk)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing container.xml")
	}

	rootfile := ""
	for _, e := range container.FindElements("//rootfiles/rootfile[@full-path]") {
		rootfile = e.SelectAttrValue("full-path", "")
	}

	if rootfile == "" {
		return nil, errors.Wrap(err, "could not find rootfile in container.xml")
	}

	opfdir := filepath.Dir(rootfile)

	rrsk, err := zfs.Open("/" + rootfile)
	if err != nil {
		return nil, errors.Wrap(err, "error reading rootfile")
	}
	defer rrsk.Close()

	opf := etree.NewDocument()
	_, err = opf.ReadFrom(rrsk)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing rootfile")
	}

	e.book.Title = filepath.Base(e.book.FilePath)
	for _, el := range opf.FindElements("//title") {
		e.book.Title = el.Text()
		break
	}
	for _, el := range opf.FindElements("//creator") {
		e.book.Author = el.Text()
		break
	}
	for _, el := range opf.FindElements("//publisher") {
		e.book.Publisher = el.Text()
		break
	}
	for _, el := range opf.FindElements("//description") {
		e.book.Description = el.Text()
		break
	}

	isbnTags := []string{
		"//source",
		"//identifier",
	}

findISBN:
	for _, tag := range isbnTags {
		for _, el := range opf.FindElements(tag) {
			val := el.Text()
			if len(val) < 10 {
				continue
			}
			if val[0:9] == "urn:isbn:" {
				val = val[9:]
			}

			if isbn.Validate(val) {
				e.book.ISBN = val
				break findISBN
			}
		}
	}

	pubDate := ""
	for _, el := range opf.FindElements("//date") {
		event := el.SelectAttrValue("opf:event", "")
		if event == "original-publication" || event == "published" || event == "publication" {
			pubDate = el.Text()
			// found a concrete publication date; we're done
			break
		} else if event == "" {
			pubDate = el.Text()
			// keep searching in case we can find a date specifically tagged as a publication date
		}
	}

	e.book.PublishDate = parsePublishDate(pubDate)

	for _, el := range opf.FindElements("//meta[@name='cover']") {
		coverid := el.SelectAttrValue("content", "")
		if coverid != "" {
			for _, f := range opf.FindElements("//[@id='" + coverid + "']") {
				coverPath := f.SelectAttrValue("href", "")
				if coverPath != "" {
					coverPath = "/" + opfdir + "/" + coverPath
					e.coverpath = &coverPath
				}
			}
			break
		}
	}

	// Calibre series metadata
	if el := opf.FindElement("//meta[@name='calibre:series']"); el != nil {
		e.book.Series = el.SelectAttrValue("content", "")

		if el := opf.FindElement("//meta[@name='calibre:series_index']"); el != nil {
			e.book.SeriesIndex, _ = strconv.ParseFloat(el.SelectAttrValue("content", "0"), 64)
		}
	}

	// EPUB3 series metadata
	if e.book.Series == "" {
		if el := opf.FindElement("//meta[@property='belongs-to-collection']"); el != nil {
			e.book.Series = strings.TrimSpace(el.Text())

			var ctype string
			if id := el.SelectAttrValue("id", ""); id != "" {
				for _, el := range opf.FindElements("//meta[@refines='#" + id + "']") {
					val := strings.TrimSpace(el.Text())
					switch el.SelectAttrValue("property", "") {
					case "collection-type":
						ctype = val
					case "group-position":
						e.book.SeriesIndex, _ = strconv.ParseFloat(val, 64)
					}
				}
			}

			if ctype != "" && ctype != "series" {
				e.book.Series, e.book.SeriesIndex = "", 0
			}
		}
	}

	return e, nil
}

func init() {
	formats.Register("epub", load)
}

func parsePublishDate(s string) time.Time {
	// handle the various dumb decisions people make when encoding dates
	format := ""
	switch len(s) {
	case 32:
		//2012-02-13T20:20:58.175203+00:00
		format = "2006-01-02T15:04:05.000000-07:00"
	case 25:
		//2000-10-31 00:00:00-06:00
		//2009-04-19T22:00:00+00:00
		format = "2006-01-02" + string(s[10]) + "15:04:05-07:00"
	case 20:
		//2016-08-11T14:09:25Z
		format = "2006-01-02T15:04:05Z"
	case 19:
		//2008-01-28T07:00:00
		//2000-10-31 00:00:00
		format = "2006-01-02" + string(s[10]) + "15:04:05"
	case 10:
		//1998-07-01
		format = "2006-01-02"
	default:
		return time.Time{}
	}

	t, err := time.Parse(format, s)
	if err != nil {
		t = time.Time{}
	}
	return t
}
