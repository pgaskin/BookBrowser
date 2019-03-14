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

	"github.com/mattn/go-zglob"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"encoding/json"
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
	seen      *SeenCache
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

	return &Indexer{paths: paths, coverpath: cp, exts: exts, seen: NewSeenCache()}, nil
}

func (i *Indexer) Load() error {
	i.indMu.Lock()
	defer i.indMu.Unlock()

	booklist := booklist.BookList{}

	jsonFilename := filepath.Join(*i.coverpath, "index.json")
	f, err := os.Open(jsonFilename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return errors.Wrap(err, "could not open index cache file")
		}
	}
	dec := json.NewDecoder(f)
	err = dec.Decode(&booklist)
	if err != nil {
		return errors.Wrap(err, "could not decode index cache file")
	}
	seen := NewSeenCache()
	for index, b := range booklist {
		seen.Add(b.FilePath, b.FileSize, b.ModTime, index)
	}

	if i.Verbose {
		log.Printf("Loaded %d items from index cache", len(booklist))
	}

	i.mu.Lock()
	i.booklist = booklist
	i.seen = seen
	i.mu.Unlock()

	return nil
}

func (i *Indexer) Save() error {
	i.indMu.Lock()
	defer i.indMu.Unlock()

	i.mu.Lock()
	booklist := i.booklist
	i.mu.Unlock()

	tmpFilename := filepath.Join(*i.coverpath, ".index.json.tmp")
	jsonFilename := filepath.Join(*i.coverpath, "index.json")
	f, err := os.Create(tmpFilename)
	if err != nil {
		f.Close()
		return errors.Wrap(err, "could not create index cache temporary file")
	}

	enc := json.NewEncoder(f)
	err = enc.Encode(&booklist)
	if err != nil {
		f.Close()
		return errors.Wrap(err, "could not encode index cache file")
	}

	err = os.Rename(tmpFilename, jsonFilename)
	if err != nil {
		return errors.Wrap(err, "could not replace index cache file with temporary file")
	}

	if i.Verbose {
		log.Printf("Saved %d items to index cache", len(booklist))
	}

	return nil
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

	// seenID may be redundant at this point given that SeenCache does essentially the same thing, but
	// seenCache is based on the mtime/size/filename of each book (for performance), whereas seenID is based on
	// the file hash
	seenID := map[string]bool{}
	seen := NewSeenCache()

	i.mu.Lock()
	bl := i.booklist
	i.mu.Unlock()

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

	exists := make([]bool, len(bl), len(bl))

	for fi, filepath := range filenames {
		if i.Verbose {
			log.Printf("Indexing %s", filepath)
		}

		stat, err := os.Stat(filepath)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "cannot stat file '%s'", filepath))
			if i.Verbose {
				log.Printf("--> Error: %v", errs[len(errs)-1])
			}
			continue
		}

		var book *booklist.Book
		hash := i.seen.Hash(filepath, stat.Size(), stat.ModTime())
		haveSeen, blIndex := i.seen.SeenHash(hash)
		if haveSeen {
			exists[blIndex] = true
			seen.AddHash(hash, blIndex)
			if i.Verbose {
				log.Printf("Already seen; not reindexing")
			}
		} else {
			// TODO: pass stat variable to i.getBook() to avoid a duplicate os.Stat() for each book
			book, err = i.getBook(filepath)
			if err != nil {
				errs = append(errs, errors.Wrapf(err, "error reading book '%s'", filepath))
				if i.Verbose {
					log.Printf("--> Error: %v", errs[len(errs)-1])
				}
				continue
			}
			if !seenID[book.ID()] {
				bl = append(bl, book)
				seenID[book.ID()] = true
				blIndex = len(bl) - 1
				seen.AddHash(hash, blIndex)
			}
		}

		i.Progress = float64(fi+1) / float64(len(filenames))
	}

	// remove any books that have disappeared since our last indexing job
	lastEntry := len(bl)-1
	for index, stillExists := range exists {
		if !stillExists {
			bl[index] = bl[lastEntry]
			lastEntry--
		}
	}
	bl = bl[0:lastEntry+1]

	i.mu.Lock()
	i.booklist = bl
	i.seen = seen
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
