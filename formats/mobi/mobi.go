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

	mobirdr "github.com/sblinch/mobi"
)

type mobi struct {
	book *booklist.Book
	coverstart int64
	coverlength int64
}

func (e *mobi) Book() *booklist.Book {
	return e.book
}

func (e *mobi) HasCover() bool {
	return e.coverstart > 0
}

func (e *mobi) GetCover() (i image.Image, err error) {
	if !e.HasCover() {
		return nil, errors.New("no cover")
	}

	f, err := os.Open(e.book.FilePath)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open book file")
	}
	defer f.Close()

	if _, err := f.Seek(e.coverstart, 0); err != nil {
		return nil, errors.Wrap(err, "unable to see to cover offset")
	}

	ltd := io.LimitReader(f,e.coverlength)
	if i, _, err = image.Decode(ltd); err != nil {
		return nil, errors.Wrap(err, "unable to decode book cover")
	}

	return i, nil
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

	r, err := mobirdr.NewReader(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	m.coverstart, m.coverlength = r.CoverOffsetLength()

	m.book.Title = r.BestTitle()

	authors := r.Authors()
	if len(authors)>0 {
		m.book.Author = authors[0]
	}

	m.book.Description = r.Description()
	m.book.Publisher = r.Publisher()

	/* // uncomment and import github.com/moraes/isbn after merging ISBN/publishing date pull request :)
	isbnStr := r.Isbn()
	if len(isbnStr)>0 && isbn.Validate(isbnStr) {
		m.book.ISBN = isbnStr
	}

	m.book.PublishDate = parsePublishDate(r.PublishingDate())
	*/

	if len(m.book.Title)==0 {
		m.book.Title = filepath.Base(filename)
	}

	return m, nil
}

func init() {
	formats.Register("mobi", load)
	formats.Register("azw", load)
	formats.Register("azw3", load)
}

/* // uncomment after merging ISBN/publishing date pull request :)
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

	t, err := time.Parse(format,s)
	if err != nil {
		t = time.Time{}
	}
	return t
}
*/