package epub

import (
	"archive/zip"
	"errors"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"golang.org/x/tools/godoc/vfs/zipfs"

	"github.com/beevik/etree"
	"github.com/geek1011/BookBrowser/formats"
	"github.com/geek1011/BookBrowser/models"
)

func indexer(filename string) (book *models.Book, cover *models.Cover, err error) {
	defer func() {
		if r := recover(); r != nil {
			book = nil
			cover = nil
			err = fmt.Errorf("Unknown error parsing book. Skipping. Error: %s", r)
		}
	}()

	var title string
	var author string
	var seriesName string
	var seriesIndex float64
	var publisher string
	var description string
	var hasCover bool
	var modtime time.Time
	var coverTmp models.Cover

	if file, err := os.Stat(filename); err == nil {
		modtime = file.ModTime()
	}

	zr, err := zip.OpenReader(filename)
	if err != nil {
		return nil, nil, err
	}

	zfs := zipfs.New(zr, "epub")

	rsk, err := zfs.Open("/META-INF/container.xml")
	if err != nil {
		return nil, nil, err
	}
	defer rsk.Close()

	container := etree.NewDocument()
	_, err = container.ReadFrom(rsk)
	if err != nil {
		return nil, nil, err
	}

	rootfile := ""
	for _, e := range container.FindElements("//rootfiles/rootfile[@full-path]") {
		rootfile = e.SelectAttrValue("full-path", "")
	}

	if rootfile == "" {
		return nil, nil, errors.New("Cannot parse container")
	}

	rrsk, err := zfs.Open("/" + rootfile)
	if err != nil {
		return nil, nil, err
	}
	defer rrsk.Close()

	opfdir := filepath.Dir(rootfile)

	opf := etree.NewDocument()
	_, err = opf.ReadFrom(rrsk)
	if err != nil {
		return nil, nil, err
	}

	title = filepath.Base(filename)
	for _, e := range opf.FindElements("//title") {
		title = e.Text()
		break
	}
	for _, e := range opf.FindElements("//creator") {
		author = e.Text()
		break
	}
	for _, e := range opf.FindElements("//publisher") {
		publisher = e.Text()
		break
	}
	for _, e := range opf.FindElements("//description") {
		description = e.Text()
		break
	}
	for _, e := range opf.FindElements("//meta[@name='calibre:series']") {
		seriesName = e.SelectAttrValue("content", "")
		break
	}
	for _, e := range opf.FindElements("//meta[@name='calibre:series_index']") {
		i, err := strconv.ParseFloat(e.SelectAttrValue("content", "0"), 64)
		if err == nil {
			seriesIndex = i
			break
		}
	}

	for _, e := range opf.FindElements("//meta[@name='cover']") {
		coverid := e.SelectAttrValue("content", "")
		if coverid != "" {
			for _, f := range opf.FindElements("//[@id='" + coverid + "']") {
				coverPath := f.SelectAttrValue("href", "")
				if coverPath != "" {
					cr, err := zfs.Open("/" + opfdir + "/" + coverPath)
					if err != nil {
						continue
					}
					defer cr.Close()

					ext := filepath.Ext(coverPath)
					if ext == ".jpeg" {
						ext = ".jpg"
					}

					switch ext {
					case ".jpg":
						coverTmp, err = jpeg.Decode(cr)
						if err != nil {
							continue
						}
					case ".gif":
						coverTmp, err = gif.Decode(cr)
						if err != nil {
							continue
						}
					case ".png":
						coverTmp, err = png.Decode(cr)
						if err != nil {
							continue
						}
					}

					hasCover = true
				}
			}
			break
		}
	}

	return models.NewBook(title, author, publisher, seriesName, seriesIndex, description, filename, hasCover, modtime, "epub"), &coverTmp, nil
}

func init() {
	formats.RegisterFormat(&formats.Format{
		Glob:      "**/*.epub",
		Extension: ".epub",
		Indexer:   indexer,
	})
}
