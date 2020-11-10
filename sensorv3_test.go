package goruuvitag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFormat3Temperature(t *testing.T) {
	assert.Equal(t, 0.0, parseFormat3Temperature(128, 0), "Negative zero is zero")
	assert.Equal(t, -2.0, parseFormat3Temperature(130, 0), "-2")
	assert.Equal(t, 2.0, parseFormat3Temperature(2, 0), "2")
	assert.Equal(t, 2.2, parseFormat3Temperature(2, 20), "fraction is ok")
	assert.Equal(t, -2.99, parseFormat3Temperature(130, 99), "fraction is ok negative")
}
