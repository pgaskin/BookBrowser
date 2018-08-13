package pdf

import (
	"crypto/sha1"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/beevik/etree"
	"github.com/geek1011/BookBrowser/booklist"
	"github.com/geek1011/BookBrowser/formats"
	"github.com/geek1011/BookBrowser/modules/util"
	"github.com/pkg/errors"
)

type pdf struct {
	book *booklist.Book
}

func (e *pdf) Book() *booklist.Book {
	return e.book
}

func (e *pdf) HasCover() bool {
	return false
}

func (e *pdf) GetCover() (i image.Image, err error) {
	return nil, errors.New("no cover")
}

func load(filename string) (bi formats.BookInfo, ferr error) {
	defer func() {
		if r := recover(); r != nil {
			bi = nil
			ferr = fmt.Errorf("unknown error: %s", r)
		}
	}()

	p := &pdf{book: &booklist.Book{}}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, errors.Wrapf(err, "could not stat book")
	}
	p.book.FilePath = filename
	p.book.FileSize = fi.Size()
	p.book.ModTime = fi.ModTime()

	s := sha1.New()
	i, err := io.Copy(s, f)
	if err == nil && i != fi.Size() {
		err = errors.New("could not read whole file")
	}
	if err != nil {
		f.Close()
		return nil, errors.Wrap(err, "could not hash book")
	}
	p.book.Hash = fmt.Sprintf("%x", s.Sum(nil))

	f.Close()

	c, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	str := string(c)
	c = []byte{}

	str = util.StringBetween(str, "<?xpacket begin", "</x:xmpmeta>")
	str = util.StringAfter(str, ">")

	xmp := etree.NewDocument()
	err = xmp.ReadFromString(str)
	if err != nil {
		return nil, err
	}

	p.book.Title = filepath.Base(filename)

	for _, e := range xmp.FindElements("//format") {
		// Make sure it is a pdf, not another piece of embedded RDF metadata
		if e.Text() != "application/pdf" {
			str = ""
			debug.FreeOSMemory()
			return p, nil
		}
		break
	}

	for _, e := range xmp.FindElements("//title/Alt/li") {
		p.book.Title = e.Text()
		break
	}

	for _, e := range xmp.FindElements("//creator/Seq/li") {
		p.book.Author = e.Text()
		break
	}

	str = ""
	debug.FreeOSMemory()

	return p, nil
}

func init() {
	formats.Register("pdf", load)
}
