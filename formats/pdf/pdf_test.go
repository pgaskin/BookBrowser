package pdf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEPUBMetadata(t *testing.T) {
	pdf, err := load("pdf_test.pdf")
	assert.Nil(t, err, "should not error when loading book")

	assert.NotNil(t, pdf, "pdf should not be nil")

	// TODO: Finish rest of tests
}
