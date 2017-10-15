package pdf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEPUBMetadata(t *testing.T) {
	book, cover, err := indexer("pdf_test.pdf")
	assert.Nil(t, err, "should not error when loading book")

	assert.NotNil(t, book, "book should not be nil")
	assert.Nil(t, cover, "should not return a cover")

	// TODO: Finish rest of tests
}
