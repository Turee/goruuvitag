package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"encoding/hex"
)

func TestIsRuuviTagValidRuuviTag(t *testing.T) {
	data, _ := hex.DecodeString("9904058000FFFFFFFF800080008000FFFFFFFFFFFFFFFFFFFFFF")

	assert.True(t, IsRuuviTag(data))
}

func TestIsRuuviTagNonRuuviTag(t *testing.T) {
	data, _ := hex.DecodeString("058000FFFFFFFF800080008000FFFFFFFFFFFFFFFFFFFFFF")

	assert.False(t, IsRuuviTag(data))
}

func TestIsRuuviTagEmptyInput(t *testing.T) {
	assert.False(t, IsRuuviTag(make([]byte, 0)))
}

func TestParseRuuviDataMissingFields(t *testing.T) {
	data, _ := hex.DecodeString("0000058000FFFFFFFF8")

	assert.NotPanics(t, func() { ParseRuuviData(data, "AB:CD:EF:01:23:45") })
}
