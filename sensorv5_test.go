package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"encoding/hex"
)

// https://github.com/ruuvi/ruuvi-sensor-protocols/blob/master/dataformat_05.md#case-valid-data
func TestParseSensorFormat5(t *testing.T) {
	data, _ := hex.DecodeString("0512FC5394C37C0004FFFC040CAC364200CDCBB8334C884F")

	results := ParseSensorFormat5(data)
	assert.Equal(t, 24.3, *results.Temperature, "Temperature fail")
	assert.Equal(t, uint32(100044), *results.Pressure, "Pressure fail")
	assert.Equal(t, 53.49, *results.Humidity, "Humidity fail")
	assert.Equal(t, 0.004, *results.AccelerationX, "Acceleration X fail")
	assert.Equal(t, -0.004, *results.AccelerationY, "Acceleration Y fail")
	assert.Equal(t, 1.036, *results.AccelerationZ, "Acceleration Z fail")
	assert.Equal(t, 2977, *results.BatteryVoltageMv, "BatteryVoltage fail")
	assert.Equal(t, 4, *results.TxPower, "TX Power fail")
	assert.Equal(t, uint16(205), *results.MeasurementSequence, "Measurement Sequence fail")
	assert.Equal(t, uint8(66), *results.MovementCounter, "Movement Counter fail")
	assert.Equal(t, "CB:B8:33:4C:88:4F", *results.MAC, "MAC fail")
}

func TestParseSensorFormat5MaxValues(t *testing.T) {
	data, _ := hex.DecodeString("057FFFFFFEFFFE7FFF7FFF7FFFFFDEFEFFFECBB8334C884F")

	results := ParseSensorFormat5(data)
	assert.Equal(t, 163.835, *results.Temperature, "Temperature fail")
	assert.Equal(t, uint32(115534), *results.Pressure, "Pressure fail")
	assert.Equal(t, 163.8350, *results.Humidity, "Humidity fail")
	assert.Equal(t, 32.767, *results.AccelerationX, "Acceleration X fail")
	assert.Equal(t, 32.767, *results.AccelerationY, "Acceleration Y fail")
	assert.Equal(t, 32.767, *results.AccelerationZ, "Acceleration Z fail")
	assert.Equal(t, 3646, *results.BatteryVoltageMv, "BatteryVoltage fail")
	assert.Equal(t, 20, *results.TxPower, "TX Power fail")
	assert.Equal(t, uint16(65534), *results.MeasurementSequence, "Measurement Sequence fail")
	assert.Equal(t, uint8(254), *results.MovementCounter, "Movement Counter fail")
	assert.Equal(t, "CB:B8:33:4C:88:4F", *results.MAC, "MAC fail")
}

func TestParseSensorFormat5MinValues(t *testing.T) {
	data, _ := hex.DecodeString("058001000000008001800180010000000000CBB8334C884F")

	results := ParseSensorFormat5(data)
	assert.Equal(t, -163.835, *results.Temperature, "Temperature fail")
	assert.Equal(t, uint32(50000), *results.Pressure, "Pressure fail")
	assert.Equal(t, 0.0, *results.Humidity, "Humidity fail")
	assert.Equal(t, -32.767, *results.AccelerationX, "Acceleration X fail")
	assert.Equal(t, -32.767, *results.AccelerationY, "Acceleration Y fail")
	assert.Equal(t, -32.767, *results.AccelerationZ, "Acceleration Z fail")
	assert.Equal(t, 1600, *results.BatteryVoltageMv, "BatteryVoltage fail")
	assert.Equal(t, -40, *results.TxPower, "TX Power fail")
	assert.Equal(t, uint16(0), *results.MeasurementSequence, "Measurement Sequence fail")
	assert.Equal(t, uint8(0), *results.MovementCounter, "Movement Counter fail")
	assert.Equal(t, "CB:B8:33:4C:88:4F", *results.MAC, "MAC fail")
}

func TestParseSensorFormat5InvalidValues(t *testing.T) {
	data, _ := hex.DecodeString("058000FFFFFFFF800080008000FFFFFFFFFFFFFFFFFFFFFF")

	results := ParseSensorFormat5(data)
	assert.Nil(t, results.Temperature)
	assert.Nil(t, results.Pressure)
	assert.Nil(t, results.Humidity)
	assert.Nil(t, results.AccelerationX)
	assert.Nil(t, results.AccelerationY)
	assert.Nil(t, results.AccelerationZ)
	assert.Nil(t, results.BatteryVoltageMv)
	assert.Nil(t, results.TxPower)
	assert.Nil(t, results.MeasurementSequence)
	assert.Nil(t, results.MovementCounter)
	assert.Nil(t, results.MAC)
}
