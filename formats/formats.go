package formats

import (
	"image"
	"path/filepath"
	"strings"

	"github.com/geek1011/BookBrowser/booklist"

	"github.com/pkg/errors"
)

var formats = map[string]func(filename string) (BookInfo, error){}

type BookInfo interface {
	Book() *booklist.Book
	HasCover() bool
	GetCover() (image.Image, error)
}

func Register(ext string, load func(filename string) (BookInfo, error)) {
	ext = strings.ToLower(ext)
	if _, ok := formats[ext]; ok {
		panic("attempted to register existing format " + ext)
	}
	formats[ext] = load
}

func Load(filename string) (BookInfo, error) {
	ext := strings.Replace(filepath.Ext(filename), ".", "", 1)
	load, ok := formats[strings.ToLower(ext)]
	if !ok {
		return nil, errors.Errorf("could not load format %s", ext)
	}
	return load(filename)
}

func GetExts() []string {
	exts := []string{}
	for ext := range formats {
		exts = append(exts, ext)
	}
	return exts
}
