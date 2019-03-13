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
		"encoding/binary"
)

type mobi struct {
	book *booklist.Book
	coverstart int64
	coverend int64
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

	if i, _, err = image.Decode(f); err != nil {
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

	for _, rec := range r.Exth.Records {
		switch(rec.RecordType) {
		case mobirdr.EXTH_COVEROFFSET:
			v := binary.BigEndian.Uint32([]byte(rec.Value))
			coverPDBOffset := r.Header.FirstImageIndex + v

			n := int(coverPDBOffset)
			if n <= int(r.Pdf.RecordsNum)-1 {
				if n+1 < int(r.Pdf.RecordsNum) {
					m.coverend = int64(r.Offsets[n+1].Offset)
				} else {
					m.coverend = -1
				}
				m.coverstart = int64(r.Offsets[n].Offset)
			}

		case mobirdr.EXTH_TITLE:
			if len(m.book.Title)==0 {
				m.book.Title = string(rec.Value)
			}
		case mobirdr.EXTH_UPDATEDTITLE:
			if len(m.book.Title)==0 {
				m.book.Title = string(rec.Value)
			}
		case mobirdr.EXTH_AUTHOR:
			if len(m.book.Author) == 0 {
				m.book.Author = string(rec.Value)
			}
		case mobirdr.EXTH_DESCRIPTION:
			if len(m.book.Description) == 0 {
				m.book.Description = string(rec.Value)
			}
		case mobirdr.EXTH_PUBLISHER:
			if len(m.book.Publisher) == 0 {
				m.book.Publisher = string(rec.Value)
			}
/* // uncomment after merging ISBN/publishing date pull request :)
		case mobirdr.EXTH_ISBN:
			if len(m.book.ISBN) == 0 {
				m.book.ISBN = string(rec.Value)
			}
		case mobirdr.EXTH_PUBLISHINGDATE:
			if m.book.PublishDate.IsZero() {
				m.book.PublishDate = parsePublishDate(string(rec.Value))
			}
*/
		}
	}

	if len(m.book.Title)==0 {
		m.book.Title = filepath.Base(filename)
	}

	return m, nil
}

func init() {
	formats.Register("mobi", load)
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