package booklist

import "sort"

type SeriesList []Series

type Series struct {
	Name, ID string
}

func (bl BookList) Series() *SeriesList {
	series := SeriesList{}
	done := map[string]bool{}
	for _, b := range bl {
		if b.Series == "" {
			continue
		}

		if done[b.SeriesID()] {
			continue
		}
		series = append(series, Series{b.Series, b.SeriesID()})
		done[b.SeriesID()] = true
	}
	return &series
}

func (bl BookList) BySeries(s Series) BookList {
	return bl.Filtered(func(b *Book) bool {
		return b.SeriesID() == s.ID
	})
}

func (sl SeriesList) Sorted(less func(a, b Series) bool) SeriesList {
	nsl := sl[:]
	sort.SliceStable(nsl, func(i, j int) bool {
		return less(sl[i], sl[j])
	})
	return nsl
}
