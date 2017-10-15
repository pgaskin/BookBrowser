package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringBetween(t *testing.T) {
	assert.Equal(t, "23", StringBetween("1234", "1", "4"))
	assert.Equal(t, "23", StringBetween("12341234", "1", "4"))
	assert.Equal(t, "", StringBetween("asdasd", "1", "4"))
}

func TestStringAfter(t *testing.T) {
	assert.Equal(t, "34", StringAfter("1234", "2"))
	assert.Equal(t, "341234", StringAfter("12341234", "2"))
	assert.Equal(t, "", StringAfter("asdasd", "2"))
}

func TestFixString(t *testing.T) {
	assert.Equal(t, `""""''`, FixString(`“‹”›‘’`))
}
