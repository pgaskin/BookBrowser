package booklist

import (
	"sort"
	"strings"
)

type BookList []*Book

func (bl BookList) Sorted(less func(a, b *Book) bool) BookList {
	nbl := bl[:]
	sort.SliceStable(nbl, func(i, j int) bool {
		return less(bl[i], bl[j])
	})
	return nbl
}

func (bl BookList) Filtered(fn func(b *Book) bool) BookList {
	nbl := BookList{}
	for _, b := range bl {
		if fn(b) {
			nbl = append(nbl, b)
		}
	}
	return nbl
}

func (bl BookList) Skip(n int) BookList {
	if n >= len(bl) {
		return BookList{}
	}
	return bl[n:]
}

func (bl BookList) Take(n int) BookList {
	if n > len(bl) {
		return bl
	}
	return bl[:n]
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
func (l BookList) SortBy(sort string) (nl BookList, sorted bool) {
	sort = strings.ToLower(sort)

	nb := l[:]

	switch sort {
	case "author-asc":
		nb = nb.Sorted(func(a, b *Book) bool {
			if a.Author != "" && b.Author != "" {
				return a.Author < b.Author
			}
			return false
		})
		break
	case "author-desc":
		nb = nb.Sorted(func(a, b *Book) bool {
			if a.Author != "" && b.Author != "" {
				return a.Author > b.Author
			}
			return false
		})
		break
	case "title-asc":
		nb = nb.Sorted(func(a, b *Book) bool {
			return a.Title < b.Title
		})
		break
	case "title-desc":
		nb = nb.Sorted(func(a, b *Book) bool {
			return a.Title > b.Title
		})
		break
	case "series-asc":
		nb = nb.Sorted(func(a, b *Book) bool {
			if a.Series != "" && b.Series != "" {
				return a.Series < b.Series
			}
			return false
		})
		break
	case "series-desc":
		nb = nb.Sorted(func(a, b *Book) bool {
			if a.Series != "" && b.Series != "" {
				return a.Series > b.Series
			}
			return false
		})
		break
	case "seriesindex-asc":
		nb = nb.Sorted(func(a, b *Book) bool {
			if a.Series != "" && b.Series != "" {
				return a.SeriesIndex < b.SeriesIndex
			}
			return false
		})
		break
	case "seriesindex-desc":
		nb = nb.Sorted(func(a, b *Book) bool {
			if a.Series != "" && b.Series != "" {
				return a.SeriesIndex > b.SeriesIndex
			}
			return false
		})
		break
	case "modified-desc":
		nb = nb.Sorted(func(a, b *Book) bool {
			return a.ModTime.Unix() > b.ModTime.Unix()
		})
		break
	default:
		return nb, false
	}

	return nb, true
}
