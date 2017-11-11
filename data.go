package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	httpURL = os.Getenv("HTTP_URL")
)

//SensorData to be posted
type SensorData struct {
	Temp          float64
	Humidity      float64
	Pressure      uint32
	Battery       uint16
	Address       string
	AccelerationX int16
	AccelerationY int16
	AccelerationZ int16
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

func sendSensorData(data *SensorData, url string) {
	s, err := json.Marshal(data)
	if err == nil {
		fmt.Printf("Posting json %s \n", string(s))
		res, err := http.Post(url, "application/json", bytes.NewReader(s))
		if err != nil {
			fmt.Printf("Error making request to elastic %s \n", err)
		} else {
			defer res.Body.Close()
			fmt.Printf("Post status %d", res.StatusCode)
		}

	} else {
		fmt.Printf("Error converting to JSON %s \n", err)
	}
}

func parseTemperature(t uint8, f uint8) float64 {
	var mask uint8
	mask = (1 << 7)
	isNegative := (t & mask) > 0
	temp := float64(t&^mask) + float64(f)/100.0
	if isNegative {
		temp *= -1
	}
	return temp
}

func parseSensorFormat3(data []byte) *SensorData {
	reader := bytes.NewReader(data)
	result := SensorFormat3{}
	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		panic(err)
	}
	sensorData := SensorData{}
	sensorData.Temp = parseTemperature(result.Temperature, result.TemperatureFraction)
	sensorData.Humidity = float64(result.Humidity) / 2.0
	sensorData.Pressure = uint32(result.Pressure) + 50000
	sensorData.Battery = result.BatteryVoltageMv
	sensorData.AccelerationX = result.AccelerationX
	sensorData.AccelerationY = result.AccelerationY
	sensorData.AccelerationZ = result.AccelerationZ
	return &sensorData
}

//ParseRuuviData parses ruuvidata
func ParseRuuviData(data []byte, a string) {

	sendData := func(sensorData *SensorData) {
		sensorData.Address = a
		sensorData.TimeStamp = time.Now()
		if httpURL != "" {
			sendSensorData(sensorData, httpURL)
		}

	}

	if len(data) == 20 && binary.LittleEndian.Uint16(data[0:2]) == 0x0499 {
		sensorFormat := data[2]
		fmt.Printf("Ruuvi data with sensor format %d\n", sensorFormat)
		switch sensorFormat {
		case 3:
			sendData(parseSensorFormat3(data))
		default:
			fmt.Printf("Unknown sensor format %d", sensorFormat)
		}
	} else {
		fmt.Printf("Not a ruuvi device \n")
	}

}
