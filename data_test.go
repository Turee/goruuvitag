package main

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"encoding/hex"
	"fmt"
)

// https://github.com/ruuvi/ruuvi-sensor-protocols/blob/master/dataformat_05.md#case-valid-data
func TestParseSensorFormat5(t *testing.T) {
	data, _ := hex.DecodeString("0512FC5394C37C0004FFFC040CAC364200CDCBB8334C884F")

	fmt.Println(data)
	fmt.Printf("%b\n", data)

	results := parseSensorFormat5(data)
	assert.Equal(t, 24.3, results.Temperature, "Temperature fail")
	assert.Equal(t, uint32(100044), results.Pressure, "Pressure fail")
	assert.Equal(t, 53.49, results.Humidity, "Humidity fail")
	assert.Equal(t, 0.004, results.AccelerationX, "Acceleration X fail")
	assert.Equal(t, -0.004, results.AccelerationY, "Acceleration Y fail")
	assert.Equal(t, 1.036, results.AccelerationZ, "Acceleration Z fail")
	assert.Equal(t, 2977, results.BatteryVoltageMv, "BatteryVoltage fail")
	assert.Equal(t, 4, results.txPower, "TX Power fail")
	assert.Equal(t, uint16(205), results.MeasurementSequence, "Measurement Sequence fail")
	assert.Equal(t, uint8(66), results.MovementCounter, "Movement Counter fail")
	assert.Equal(t, "CB:B8:33:4C:88:4F", results.MAC, "MAC fail")
}

func TestParseFormat3Temperature(t *testing.T) {
	assert.Equal(t, 0.0, parseFormat3Temperature(128, 0), "Negative zero is zero")
	assert.Equal(t, -2.0, parseFormat3Temperature(130, 0), "-2")
	assert.Equal(t, 2.0, parseFormat3Temperature(2, 0), "2")
	assert.Equal(t, 2.2, parseFormat3Temperature(2, 20), "fraction is ok")
	assert.Equal(t, -2.99, parseFormat3Temperature(130, 99), "fraction is ok negative")
}
