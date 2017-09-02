package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEPUBMetadata(t *testing.T) {
	td, err := ioutil.TempDir("", "bookbrowser")
	assert.Nil(t, err, "should not error when creating temp dir")
	defer os.RemoveAll(td)

	book, err := NewBookFromFile("testdata/books/test1.epub", td)
	assert.Nil(t, err, "should not error when loading book")

	coverfile := filepath.Join(td, book.ID+".jpg")
	if _, err := os.Stat(coverfile); err != nil {
		assert.Nil(t, err, "cover file should exist")
	}

	coverthumbfile := filepath.Join(td, book.ID+"_thumb.jpg")
	if _, err := os.Stat(coverthumbfile); err != nil {
		if os.IsNotExist(err) {
			assert.Nil(t, err, "cover thumbnail file should exist")
		}
	}

	assert.Equal(t, "BookBrowser Test Book 1", book.Title, "title")
	assert.Equal(t, "Patrick G", book.Author.Name, "author")
	assert.Equal(t, "Patrick G", book.Publisher, "publisher")
	assert.Equal(t, "<p>This is a test book for <i>BookBrowser</i>, a ebook content server.</p>", book.Description, "description")
	assert.Equal(t, "epub", book.FileType, "filetype")
	assert.True(t, book.HasCover, "should have a cover")
	assert.Equal(t, "Test Series", book.Series.Name, "series name")
	assert.Equal(t, float64(1), book.Series.Index, "series index")
	assert.Equal(t, "a611744562", book.ID, "book id")
}
