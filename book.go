package main

import (
	"archive/zip"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"time"

	"strings"

	"github.com/beevik/etree"
	zglob "github.com/mattn/go-zglob"
	"github.com/nfnt/resize"
	"golang.org/x/tools/godoc/vfs/zipfs"
)

// Series represents a book series
type Series struct {
	Name  string  `json:"name,omitempty"`
	ID    string  `json:"id,omitempty"`
	Index float64 `json:"index,omitempty"`
}

// Book represents a book
type Book struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author,omitempty"`
	AuthorID    string    `json:"authorid"`
	Publisher   string    `json:"publisher,omitempty"`
	Description string    `json:"description,omitempty"`
	Series      Series    `json:"series,omitempty"`
	Filepath    string    `json:"filepath"`
	HasCover    bool      `json:"hascover"`
	ModTime     time.Time `json:"modtime,omitempty"`
	FileType    string    `json:"filetype,omitempty"`
}

// NewBookFromFile creates a book object from a file
func NewBookFromFile(path, coverpath string) (bk *Book, err error) {
	defer func() {
		if r := recover(); r != nil {
			bk = nil
			err = fmt.Errorf("Unknown error parsing book. Skipping. Error: %s", r)
		}
	}()

	book := new(Book)
	book.Title = filepath.Base(path)
	book.Filepath = path
	book.FileType = strings.ToLower(strings.Replace(filepath.Ext(path), ".", "", -1))

	if file, err := os.Stat(path); err == nil {
		book.ModTime = file.ModTime()
	}

	switch ft := book.FileType; ft {
	case "pdf":
		book.Title = filepath.Base(path)

		m, err := GetPDFMeta(path)
		if err == nil {
			book.Title = m.Title
			book.Author = m.Author
		}

		id := sha1.New()
		io.WriteString(id, book.Author)
		book.AuthorID = hex.EncodeToString(id.Sum(nil))[:10]
		io.WriteString(id, book.Series.Name)
		io.WriteString(id, book.Title)
		book.ID = hex.EncodeToString(id.Sum(nil))[:10]
	case "epub":
		zr, err := zip.OpenReader(path)
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
		book.Title = filepath.Base(path)
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
						cpath := filepath.Join(coverpath, book.ID+".jpg")
						thumbpath := filepath.Join(coverpath, book.ID+"_thumb"+".jpg")

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
	default:
		return nil, fmt.Errorf("Unknown filetype: %s", book.FileType)
	}

	return book, nil
}

// BookList is a slice of books
type BookList []Book

// NewBookListFromDir creates a BookList from the books in a dir. It will still return if there are errors indexing some of the books.
func NewBookListFromDir(path, coverdir string, printlog bool) (*BookList, error) {
	matches, err := zglob.Glob(filepath.Join(path, "/**/*.epub"))
	if err != nil {
		return nil, err
	}

	pdfmatches, err := zglob.Glob(filepath.Join(path, "/**/*.pdf"))
	if err != nil {
		return nil, err
	}
	matches = append(matches, pdfmatches...)

	var books BookList
	for i, filename := range matches {
		if printlog {
			log.Printf("%.f%% Indexing %s\n", float64(i)/float64(len(matches))*100, filename)
		}
		book, err := NewBookFromFile(filename, coverdir)
		if err != nil {
			if printlog {
				log.Printf("Error indexing %s: %s\n", filename, err)
			}
			continue
		}
		books = append(books, *book)
	}
	debug.FreeOSMemory()
	return &books, nil
}
