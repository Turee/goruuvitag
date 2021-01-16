package main

import (
	"bytes"
	"encoding/binary"
)

//SensorFormat3 RuuviData
type SensorFormat3 struct {
	ManufacturerID      uint16
	DataFormat          uint8
	Humidity            uint8
	Temperature         uint8
	TemperatureFraction uint8
	Pressure            uint16
	AccelerationX       int16
	AccelerationY       int16
	AccelerationZ       int16
	BatteryVoltageMv    uint16
}

func parseFormat3Temperature(t uint8, f uint8) float64 {
	var mask uint8
	mask = (1 << 7)
	isNegative := (t & mask) > 0
	temp := float64(t&^mask) + float64(f)/100.0
	if isNegative {
		temp *= -1
	}
	return temp
}

// ParseSensorFormat3 parses according to https://github.com/ruuvi/ruuvi-sensor-protocols
func ParseSensorFormat3(data []byte, macAddress string) *SensorData {
	reader := bytes.NewReader(data)
	result := SensorFormat3{}
	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		panic(err)
	}
	sensorData := SensorData{}
	temperature := parseFormat3Temperature(result.Temperature, result.TemperatureFraction)
	humidity := float64(result.Humidity) / 2.0
	pressure := uint32(result.Pressure) + 50000
	batteryVoltage := int(result.BatteryVoltageMv)
	accelerationX := float64(result.AccelerationX) / 1000
	accelerationY := float64(result.AccelerationY) / 1000
	accelerationZ := float64(result.AccelerationZ) / 1000

	sensorData.Temperature = &temperature
	sensorData.Humidity = &humidity
	sensorData.Pressure = &pressure
	sensorData.BatteryVoltageMv = &batteryVoltage
	sensorData.AccelerationX = &accelerationX
	sensorData.AccelerationY = &accelerationY
	sensorData.AccelerationZ = &accelerationZ
	sensorData.MAC = &macAddress
	return &sensorData
}
