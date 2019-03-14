package booklist

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

type cachedID struct {
	id string
	value string
}

type Book struct {
	Hash     string
	FilePath string
	FileSize int64
	ModTime  time.Time

	HasCover    bool
	Title       string
	Author      string
	Description string
	Series      string
	SeriesIndex float64
	Publisher   string

	seriesid	cachedID
	authorid	cachedID
}

func (b *Book) ID() string {
	return b.Hash[:10]
}

func (b *Book) cachedIDStr(cached *cachedID, value string) string {
	if len(cached.id)==0 || value != cached.value {
		cached.value = value
		cached.id = fmt.Sprintf("%x", sha1.Sum([]byte(value)))[:10]
	}

	return cached.id
}

func (b *Book) AuthorID() string {
	return b.cachedIDStr(&b.authorid,b.Author)
}

func (b *Book) SeriesID() string {
	return b.cachedIDStr(&b.seriesid,b.Series)
}

func (b *Book) FileType() string {
	return strings.Replace(strings.ToLower(filepath.Ext(b.FilePath)), ".", "", -1)
}
