package booklist

import "sort"

type AuthorList []struct{ Name, ID string }

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
		authors = append(authors, struct{ Name, ID string }{b.Author, b.AuthorID()})
		done[b.AuthorID()] = true
	}
	return &authors
}

func (al AuthorList) Sorted(less func(a, b struct{ Name, ID string }) bool) AuthorList {
	nal := al[:]
	sort.SliceStable(nal, func(i, j int) bool {
		return less(al[i], al[j])
	})
	return nal
}
