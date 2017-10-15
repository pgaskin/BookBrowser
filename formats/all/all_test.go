package all

import (
	"testing"

	"github.com/geek1011/BookBrowser/formats"
	"github.com/stretchr/testify/assert"
)

func TestFormats(t *testing.T) {
	assert.Equal(t, 2, len(formats.Formats), "should be two formats registered")
}
