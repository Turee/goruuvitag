package main

import (
	"testing"
	
	"encoding/hex"
)

func assertEquals(t *testing.T, first float64, second float64, description string) {
	if (first != second) {
		t.Errorf("%f != %f; %s", first, second, description)
	}
}

// https://github.com/ruuvi/ruuvi-sensor-protocols/blob/master/dataformat_05.md#case-valid-data
func TestParseSensorFormat5(t *testing.T) {
	data, _ := hex.DecodeString("0512FC5394C37C0004FFFC040CAC364200CDCBB8334C884F")
	
	results := parseSensorFormat5(data)
	
	
	
}

func TestParseTemperature(t *testing.T) {
	assertEquals(t, 0.0, parseTemperature(128, 0), "Negative zero is zero")
	assertEquals(t, -2.0, parseTemperature(130, 0), "-2")
	assertEquals(t, 2.0, parseTemperature(2, 0), "2")
	assertEquals(t, 2.2, parseTemperature(2, 20), "fraction is ok")
	assertEquals(t, -2.99, parseTemperature(130, 99), "fraction is ok negative")
}
