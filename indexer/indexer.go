package indexer

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/geek1011/BookBrowser/booklist"
	"github.com/geek1011/BookBrowser/formats"

	zglob "github.com/mattn/go-zglob"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
)

type Indexer struct {
	Verbose   bool
	Progress  float64
	coverpath *string
	paths     []string
	exts      []string
	booklist  booklist.BookList
	mu        sync.Mutex
	indMu     sync.Mutex
}

func New(paths []string, coverpath *string, exts []string) (*Indexer, error) {
	for i := range paths {
		p, err := filepath.Abs(paths[i])
		if err != nil {
			return nil, errors.Wrap(err, "error resolving path")
		}
		paths[i] = p
	}

	cp := (*string)(nil)
	if coverpath != nil {
		p, err := filepath.Abs(*coverpath)
		if err != nil {
			return nil, errors.Wrap(err, "error resolving cover path")
		}
		cp = &p
	}

	return &Indexer{paths: paths, coverpath: cp, exts: exts}, nil
}

func (i *Indexer) Refresh() ([]error, error) {
	i.indMu.Lock()
	defer i.indMu.Unlock()

	defer func() {
		i.Progress = 0
	}()

	errs := []error{}

	if len(i.paths) < 1 {
		return errs, errors.New("no paths to index")
	}

	booklist := booklist.BookList{}
	seen := map[string]bool{}

	filenames := []string{}
	for _, path := range i.paths {
		for _, ext := range i.exts {
			l, err := zglob.Glob(filepath.Join(path, "**", fmt.Sprintf("*.%s", ext)))
			if l != nil {
				filenames = append(filenames, l...)
			}
			if err != nil {
				errs = append(errs, errors.Wrapf(err, "error scanning '%s' for type '%s'", path, ext))
				if i.Verbose {
					log.Printf("Error: %v", errs[len(errs)-1])
				}
			}
		}
	}

	for fi, filepath := range filenames {
		if i.Verbose {
			log.Printf("Indexing %s", filepath)
		}

		book, err := i.getBook(filepath)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "error reading book '%s'", filepath))
			if i.Verbose {
				log.Printf("--> Error: %v", errs[len(errs)-1])
			}
			continue
		}
		if !seen[book.ID()] {
			booklist = append(booklist, book)
			seen[book.ID()] = true
		}

		i.Progress = float64(fi+1) / float64(len(filenames))
	}

	i.mu.Lock()
	i.booklist = booklist
	i.mu.Unlock()

	return errs, nil
}

func (i *Indexer) BookList() booklist.BookList {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.booklist
}

func (i *Indexer) getBook(filename string) (*booklist.Book, error) {
	// TODO: caching
	bi, err := formats.Load(filename)
	if err != nil {
		return nil, errors.Wrap(err, "error loading book")
	}

	b := bi.Book()
	b.HasCover = false
	if i.coverpath != nil && bi.HasCover() {
		coverpath := filepath.Join(*i.coverpath, fmt.Sprintf("%s.jpg", b.ID()))
		thumbpath := filepath.Join(*i.coverpath, fmt.Sprintf("%s_thumb.jpg", b.ID()))

		_, err := os.Stat(coverpath)
		_, errt := os.Stat(thumbpath)
		if err != nil || errt != nil {
			i, err := bi.GetCover()
			if err != nil {
				return nil, errors.Wrap(err, "error getting cover")
			}

			f, err := os.Create(coverpath)
			if err != nil {
				return nil, errors.Wrap(err, "could not create cover file")
			}
			defer f.Close()

			err = jpeg.Encode(f, i, nil)
			if err != nil {
				os.Remove(coverpath)
				return nil, errors.Wrap(err, "could not write cover file")
			}

			ti := resize.Thumbnail(400, 400, i, resize.Bicubic)

			tf, err := os.Create(thumbpath)
			if err != nil {
				return nil, errors.Wrap(err, "could not create cover thumbnail file")
			}
			defer tf.Close()

			err = jpeg.Encode(tf, ti, nil)
			if err != nil {
				os.Remove(coverpath)
				return nil, errors.Wrap(err, "could not write cover thumbnail file")
			}
		}

		b.HasCover = true
	}

	return b, nil
}
