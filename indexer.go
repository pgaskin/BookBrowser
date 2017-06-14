package main

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/beevik/etree"
	zglob "github.com/mattn/go-zglob"
	"github.com/nfnt/resize"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

func getMetadata(epub string) (*Book, error) {
	zr, err := zip.OpenReader(epub)
	if err != nil {
		return nil, err
	}

	zfs := zipfs.New(zr, "epub")

	rsk, err := zfs.Open("/META-INF/container.xml")
	if err != nil {
		return nil, err
	}
	defer rsk.Close()
	container := etree.NewDocument()
	_, err = container.ReadFrom(rsk)
	if err != nil {
		return nil, err
	}
	rootfile := ""
	for _, e := range container.FindElements("//rootfiles/rootfile[@full-path]") {
		rootfile = e.SelectAttrValue("full-path", "")
	}
	if rootfile == "" {
		return nil, errors.New("Cannot parse container")
	}

	book := new(Book)

	rrsk, err := zfs.Open("/" + rootfile)
	if err != nil {
		return nil, err
	}
	defer rrsk.Close()
	opfdir := filepath.Dir(rootfile)
	opf := etree.NewDocument()
	_, err = opf.ReadFrom(rrsk)
	if err != nil {
		return nil, err
	}
	book.Title = filepath.Base(epub)
	for _, e := range opf.FindElements("//title") {
		book.Title = e.Text()
		break
	}
	for _, e := range opf.FindElements("//creator") {
		book.Author = e.Text()
		break
	}
	for _, e := range opf.FindElements("//publisher") {
		book.Publisher = e.Text()
		break
	}
	for _, e := range opf.FindElements("//description") {
		book.Description = e.Text()
		break
	}
	for _, e := range opf.FindElements("//meta[@name='calibre:series']") {
		book.Series.Name = e.SelectAttrValue("content", "")
		seriesid := sha1.New()
		io.WriteString(seriesid, book.Series.Name)
		book.Series.ID = hex.EncodeToString(seriesid.Sum(nil))[:10]
		break
	}
	for _, e := range opf.FindElements("//meta[@name='calibre:series_index']") {
		i, err := strconv.ParseFloat(e.SelectAttrValue("content", "0"), 64)
		if err == nil {
			book.Series.Index = i
			break
		}
	}

	id := sha1.New()
	io.WriteString(id, book.Author)
	book.AuthorID = hex.EncodeToString(id.Sum(nil))[:10]
	io.WriteString(id, book.Series.Name)
	io.WriteString(id, book.Title)
	book.ID = hex.EncodeToString(id.Sum(nil))[:10]

	for _, e := range opf.FindElements("//meta[@name='cover']") {
		coverid := e.SelectAttrValue("content", "")
		if coverid != "" {
			for _, f := range opf.FindElements("//[@id='" + coverid + "']") {
				cover := f.SelectAttrValue("href", "")
				if cover != "" {
					cr, err := zfs.Open("/" + opfdir + "/" + cover)
					if err != nil {
						continue
					}
					defer cr.Close()

					ext := filepath.Ext(cover)
					if ext == ".jpeg" {
						ext = ".jpg"
					}
					cpath := filepath.Join(*tempdir, book.ID+".jpg")
					thumbpath := filepath.Join(*tempdir, book.ID+"_thumb"+".jpg")

					var img image.Image

					switch ext {
					case ".jpg":
						img, err = jpeg.Decode(cr)
						if err != nil {
							continue
						}
					case ".gif":
						img, err = gif.Decode(cr)
						if err != nil {
							continue
						}
					case ".png":
						img, err = png.Decode(cr)
						if err != nil {
							continue
						}
					}

					coverfile, err := os.Create(cpath)
					if err != nil {
						continue
					}
					defer coverfile.Close()
					err = jpeg.Encode(coverfile, img, nil)
					if err != nil {
						continue
					}

					// Better quality: thumb := resize.Resize(200, 0, img, resize.Lanczos2)
					thumb := resize.Resize(200, 0, img, resize.Bicubic)
					thumbfile, err := os.Create(thumbpath)
					if err != nil {
						continue
					}
					defer thumbfile.Close()
					err = jpeg.Encode(thumbfile, thumb, nil)
					if err != nil {
						continue
					}
					book.HasCover = true
					break
				}
			}
			break
		}
	}

	book.Filepath = epub

	return book, nil
}

func indexBooks() ([]Book, error) {
	matches, err := zglob.Glob(filepath.Join(*bookdir, "/**/*.epub"))
	if err != nil {
		return nil, err
	}

	books := []Book{}
	for i, filename := range matches {
		log.Printf("%.f%% Indexing %s\n", float64(i)/float64(len(matches))*100, filename)
		book, err := getMetadata(filename)
		if err != nil {
			log.Printf("Error indexing %s: %s\n", filename, err)
			continue
		}
		books = append(books, *book)
	}
	return books, nil
}
