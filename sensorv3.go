package goruuvitag

import (
	"bytes"
	"encoding/binary"
	"time"
)

//SensorData to be posted
type SensorData struct {
	Temperature   float64
	Humidity      float64
	Pressure      uint32
	Battery       uint16
	AccelerationX int16
	AccelerationY int16
	AccelerationZ int16
	Address       string
	TimeStamp     time.Time
}

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

// https://github.com/ruuvi/ruuvi-sensor-protocols
func ParseSensorFormat3(data []byte) *SensorData {
	reader := bytes.NewReader(data)
	result := SensorFormat3{}
	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		panic(err)
	}
	sensorData := SensorData{}
	sensorData.Temperature = parseFormat3Temperature(result.Temperature, result.TemperatureFraction)
	sensorData.Humidity = float64(result.Humidity) / 2.0
	sensorData.Pressure = uint32(result.Pressure) + 50000
	sensorData.Battery = result.BatteryVoltageMv
	sensorData.AccelerationX = result.AccelerationX
	sensorData.AccelerationY = result.AccelerationY
	sensorData.AccelerationZ = result.AccelerationZ
	return &sensorData
}
