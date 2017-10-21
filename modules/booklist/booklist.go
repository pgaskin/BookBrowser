package booklist

import (
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/geek1011/BookBrowser/formats"
	"github.com/geek1011/BookBrowser/models"
	zglob "github.com/mattn/go-zglob"
	"github.com/nfnt/resize"
)

// BookList represents a list of Books
type BookList []*models.Book

// IndexerError represents a indexer error.
type IndexerError struct {
	Filename string
	Error    error
}

// NewBookListFromDir creates a BookList from a directory of books.
func NewBookListFromDir(dir, coverOutDir string, verbose bool) (*BookList, []*IndexerError) {
	errors := []*IndexerError{}
	books := BookList{}

	filenames := map[string][]string{}
	for _, format := range formats.Formats {
		matches, err := zglob.Glob(filepath.Join(dir, format.Glob))
		if err != nil {
			errors = append(errors, &IndexerError{
				Filename: format.Glob,
				Error:    fmt.Errorf("error getting list of matched filenames for format %s: %v", format.Extension, err),
			})

			if verbose {
				log.Printf("Error getting matches for %s: %v", format.Glob, err)
			}

			continue
		}

		filenames[format.Extension] = matches
	}

	total := 0
	for _, i := range filenames {
		total += len(i)
	}

	current := 0
	for _, format := range formats.Formats {
		for _, filename := range filenames[format.Extension] {
			current++

			if verbose {
				log.Printf("[%v/%v] Indexing %s", current, total, filename)
			}

			book, cover, err := format.Indexer(filename)
			if err != nil {
				errors = append(errors, &IndexerError{
					Filename: filename,
					Error:    err,
				})

				if verbose {
					log.Printf("[%v/%v] Error indexing %s: %v", current, total, filename, err)
				}

				continue
			}

			if cover != nil {
				coverPath := filepath.Join(coverOutDir, book.ID+".jpg")
				thumbPath := filepath.Join(coverOutDir, book.ID+"_thumb.jpg")

				coverFile, err := os.Create(coverPath)
				if err != nil {
					continue
				}
				defer coverFile.Close()

				err = jpeg.Encode(coverFile, cover, nil)
				if err != nil {
					continue
				}

				// Better quality: thumb := resize.Resize(200, 0, img, resize.Lanczos2)
				thumb := resize.Resize(200, 0, cover, resize.Bicubic)
				thumbFile, err := os.Create(thumbPath)
				if err != nil {
					continue
				}
				defer thumbFile.Close()

				err = jpeg.Encode(thumbFile, thumb, nil)
				if err != nil {
					continue
				}

				book.HasCover = true
			}

			books = append(books, book)
		}
	}

	debug.FreeOSMemory()
	return &books, errors
}

// Sorted returns a copy of the BookList sorted by the function
func (l *BookList) Sorted(sorter func(a, b *models.Book) bool) BookList {
	// Make a copy
	sorted := make(BookList, len(*l))
	copy(sorted, *l)
	// Sort the copy
	sort.Slice(sorted, func(i, j int) bool {
		return sorter(sorted[i], sorted[j])
	})
	return sorted
}

// Filtered returns a copy of the BookList filtered by the function
func (l *BookList) Filtered(filterer func(a *models.Book) bool) *BookList {
	filtered := BookList{}
	for _, a := range *l {
		if filterer(a) {
			filtered = append(filtered, a)
		}
	}

	return &filtered
}

// AuthorList is a list of authors
type AuthorList []*models.Author

// SeriesList is a list of series
type SeriesList []*models.Series

// GetAuthors gets the authors in a BookList
func (l *BookList) GetAuthors() *AuthorList {
	authors := AuthorList{}
	done := map[string]bool{}
	for _, b := range *l {
		if b.Author == nil {
			continue
		}

		if done[b.Author.ID] {
			continue
		}
		authors = append(authors, b.Author)
		done[b.Author.ID] = true
	}
	return &authors
}

// Sorted returns a copy of the AuthorList sorted by the function
func (l *AuthorList) Sorted(sorter func(a, b *models.Author) bool) *AuthorList {
	// Make a copy
	sorted := make(AuthorList, len(*l))
	copy(sorted, *l)
	// Sort the copy
	sort.Slice(sorted, func(i, j int) bool {
		return sorter(sorted[i], sorted[j])
	})
	return &sorted
}

// Filtered returns a copy of the AuthorList filtered by the function
func (l *AuthorList) Filtered(filterer func(a *models.Author) bool) *AuthorList {
	filtered := AuthorList{}
	for _, a := range *l {
		if filterer(a) {
			filtered = append(filtered, a)
		}
	}

	return &filtered
}

// GetSeries gets the series in a BookList
func (l *BookList) GetSeries() *SeriesList {
	series := SeriesList{}
	done := map[string]bool{}
	for _, b := range *l {
		if b.Series == nil {
			continue
		}

		if done[b.Series.ID] {
			continue
		}
		series = append(series, b.Series)
		done[b.Series.ID] = true
	}

	return &series
}

// Sorted returns a copy of the SeriesList sorted by the function
func (l *SeriesList) Sorted(sorter func(a, b *models.Series) bool) *SeriesList {
	// Make a copy
	sorted := make(SeriesList, len(*l))
	copy(sorted, *l)
	// Sort the copy
	sort.Slice(sorted, func(i, j int) bool {
		return sorter(sorted[i], sorted[j])
	})
	return &sorted
}

// Filtered returns a copy of the SeriesList filtered by the function
func (l *SeriesList) Filtered(filterer func(a *models.Series) bool) *SeriesList {
	filtered := SeriesList{}
	for _, a := range *l {
		if filterer(a) {
			filtered = append(filtered, a)
		}
	}

	return &filtered
}

// HasBook checks whether a book with an id exists
func (l *BookList) HasBook(id string) bool {
	exists := false
	for _, b := range *l {
		if b.ID == id {
			exists = true
		}
	}
	return exists
}

// HasAuthor checks whether an author with an id exists
func (l *BookList) HasAuthor(id string) bool {
	exists := false
	for _, b := range *l {
		if b.Author == nil {
			continue
		}

		if b.Author.ID == id {
			exists = true
		}
	}
	return exists
}

// HasSeries checks whether a series with an id exists
func (l *BookList) HasSeries(id string) bool {
	exists := false
	for _, b := range *l {
		if b.Series == nil {
			continue
		}
		if b.Series.ID == id {
			exists = true
		}
	}
	return exists
}

// SortBy sorts by sort, and returns a sorted copy. If sorter is invalid, it returns the original list.
//
// sort can be:
// - author-asc
// - author-desc
// - title-asc
// - title-desc
// - series-asc
// - series-desc
// - seriesindex-asc
// - seriesindex-desc
// - modified-desc
func (l *BookList) SortBy(sort string) (nl *BookList, sorted bool) {
	sort = strings.ToLower(sort)

	nb := *l

	switch sort {
	case "author-asc":
		nb = nb.Sorted(func(a, b *models.Book) bool {
			if a.Author != nil && b.Author != nil {
				return a.Author.Name < b.Author.Name
			}
			return false
		})
		break
	case "author-desc":
		nb = nb.Sorted(func(a, b *models.Book) bool {
			if a.Author != nil && b.Author != nil {
				return a.Author.Name > b.Author.Name
			}
			return false
		})
		break
	case "title-asc":
		nb = nb.Sorted(func(a, b *models.Book) bool {
			return a.Title < b.Title
		})
		break
	case "title-desc":
		nb = nb.Sorted(func(a, b *models.Book) bool {
			return a.Title > b.Title
		})
		break
	case "series-asc":
		nb = nb.Sorted(func(a, b *models.Book) bool {
			if a.Series != nil && b.Series != nil {
				return a.Series.Name < b.Series.Name
			}
			return false
		})
		break
	case "series-desc":
		nb = nb.Sorted(func(a, b *models.Book) bool {
			if a.Series != nil && b.Series != nil {
				return a.Series.Name > b.Series.Name
			}
			return false
		})
		break
	case "seriesindex-asc":
		nb = nb.Sorted(func(a, b *models.Book) bool {
			if a.Series != nil && b.Series != nil {
				return a.Series.Index < b.Series.Index
			}
			return false
		})
		break
	case "seriesindex-desc":
		nb = nb.Sorted(func(a, b *models.Book) bool {
			if a.Series != nil && b.Series != nil {
				return a.Series.Index > b.Series.Index
			}
			return false
		})
		break
	case "modified-desc":
		nb = nb.Sorted(func(a, b *models.Book) bool {
			return a.ModTime.Unix() > b.ModTime.Unix()
		})
		break
	default:
		return &nb, false
	}

	return &nb, true
}
