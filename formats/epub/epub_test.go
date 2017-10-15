package epub

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEPUBMetadata(t *testing.T) {
	book, cover, err := indexer("epub_test.epub")
	assert.Nil(t, err, "should not error when loading book")

	assert.NotNil(t, book, "book should not be nil")
	assert.NotNil(t, cover, "should return a cover")
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
