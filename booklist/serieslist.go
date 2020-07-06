package booklist

import "sort"

type SeriesList []struct{ Name, ID string }

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
		series = append(series, struct{ Name, ID string }{b.Series, b.SeriesID()})
		done[b.SeriesID()] = true
	}
	return &series
}

func (sl SeriesList) Sorted(less func(a, b struct{ Name, ID string }) bool) SeriesList {
	nsl := sl[:]
	sort.SliceStable(nsl, func(i, j int) bool {
		return less(sl[i], sl[j])
	})
	return nsl
}

func (sl SeriesList) Skip(n int) SeriesList {
	if n >= len(sl) {
		return SeriesList{}
	}
	return sl[n:]
}

func (sl SeriesList) Take(n int) SeriesList {
	if n > len(sl) {
		return sl
	}
	return sl[:n]
}