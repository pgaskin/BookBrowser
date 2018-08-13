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

	"github.com/geek1011/BookBrowser/booklist"
	"github.com/geek1011/BookBrowser/formats"

	"github.com/beevik/etree"
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
	for _, el := range opf.FindElements("//meta[@name='calibre:series']") {
		s := el.SelectAttrValue("content", "")
		e.book.Series = s
	}

	if e.book.Series != "" {
		for _, el := range opf.FindElements("//meta[@name='calibre:series_index']") {
			i, _ := strconv.ParseFloat(el.SelectAttrValue("content", "0"), 64)
			e.book.SeriesIndex = i
			break
		}
	}

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

	return e, nil
}

func init() {
	formats.Register("epub", load)
}
