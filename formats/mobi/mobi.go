package mobi

import (
	"crypto/sha1"
	"fmt"
	"image"
	"io"
	"os"
	"path/filepath"

	"github.com/geek1011/BookBrowser/booklist"
	"github.com/geek1011/BookBrowser/formats"
	"github.com/pkg/errors"
)

type mobi struct {
	book *booklist.Book
}

func (e *mobi) Book() *booklist.Book {
	return e.book
}

func (e *mobi) HasCover() bool {
	return false
}

func (e *mobi) GetCover() (i image.Image, err error) {
	return nil, errors.New("no cover")
}

func load(filename string) (bi formats.BookInfo, ferr error) {
	defer func() {
		if r := recover(); r != nil {
			bi = nil
			ferr = fmt.Errorf("unknown error: %s", r)
		}
	}()

	m := &mobi{book: &booklist.Book{}}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, errors.Wrapf(err, "could not stat book")
	}
	m.book.FilePath = filename
	m.book.FileSize = fi.Size()
	m.book.ModTime = fi.ModTime()

	s := sha1.New()
	i, err := io.Copy(s, f)
	if err == nil && i != fi.Size() {
		err = errors.New("could not read whole file")
	}
	if err != nil {
		f.Close()
		return nil, errors.Wrap(err, "could not hash book")
	}
	m.book.Hash = fmt.Sprintf("%x", s.Sum(nil))

	f.Close()

	m.book.Title = filepath.Base(filename)

	return m, nil
}

func init() {
	formats.Register("mobi", load)
}
