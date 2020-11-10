package goruuvitag

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var (
	httpURL = os.Getenv("HTTP_URL")
)

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

// IsRuuviTag A helper to check if the manufacturer id of a ble advertisement matches Ruuvi's
func IsRuuviTag(data []byte) bool {
	return binary.LittleEndian.Uint16(data[0:2]) == 0x0499
}

//ParseRuuviData parses ruuvidata
func ParseRuuviData(data []byte, a string) (*SensorData, *SensorData) {
	/*
		sendData := func(sensorData *SensorData) {
			sensorData.Address = a
			sensorData.TimeStamp = time.Now()
			if httpURL != "" {
				sendSensorData(sensorData, httpURL)
			}
		}
	*/

	sensorFormat := data[2]
	fmt.Printf("Ruuvi data with sensor format %d\n", sensorFormat)
	switch sensorFormat {
	case 3:
		ParseSensorFormat3(data)
	case 5:
		ParseSensorFormat5(data)
	default:
		fmt.Printf("Unknown sensor format %d", sensorFormat)
	}

	return nil, nil
}
