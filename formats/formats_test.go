package formats

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/geek1011/BookBrowser/models"
)

func TestRegisterFormat(t *testing.T) {
	l := len(Formats)

	RegisterFormat(&Format{
		Glob:      "**/*.test",
		Extension: "test",
		Indexer: func(filename string) (*models.Book, *models.Cover, error) {
			return nil, nil, nil
		},
	})
	l++

	assert.Equal(t, l, len(Formats), "number of formats should have increased by one")
}
