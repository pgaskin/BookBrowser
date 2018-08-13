package booklist

import (
	"crypto/sha1"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

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
}

func (b *Book) ID() string {
	return b.Hash[:10]
}

func (b *Book) AuthorID() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(b.Author)))[:10]
}

func (b *Book) SeriesID() string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(b.Series)))[:10]
}

func (b *Book) FileType() string {
	return strings.Replace(strings.ToLower(filepath.Ext(b.FilePath)), ".", "", -1)
}
