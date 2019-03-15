package booklist

import "sort"

type AuthorList []Author

type Author struct {
	Name, ID string
}

func (bl BookList) Authors() *AuthorList {
	authors := AuthorList{}
	done := map[string]bool{}
	for _, b := range bl {
		if b.Author == "" {
			continue
		}

		if done[b.AuthorID()] {
			continue
		}
		authors = append(authors, Author{b.Author, b.AuthorID()})
		done[b.AuthorID()] = true
	}
	return &authors
}

func (bl BookList) ByAuthor(a Author) BookList {
	return bl.Filtered(func(b *Book) bool {
		return b.AuthorID() == a.ID
	})
}

func (al AuthorList) Sorted(less func(a, b Author) bool) AuthorList {
	nal := al[:]
	sort.SliceStable(nal, func(i, j int) bool {
		return less(al[i], al[j])
	})
	return nal
}
