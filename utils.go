package main

import (
	"encoding/binary"
	"errors"
	"log"
)

// SensorData to be posted
type SensorData struct {
	Temperature         *float64
	Humidity            *float64
	Pressure            *uint32
	BatteryVoltageMv    *int
	TxPower             *int
	AccelerationX       *float64
	AccelerationY       *float64
	AccelerationZ       *float64
	Acceleration        *float64
	MovementCounter     *uint8
	MeasurementSequence *uint16
	MAC                 *string
}

// IsRuuviTag A helper to check if the manufacturer id of a ble advertisement matches Ruuvi's
func IsRuuviTag(data []byte) bool {
	return len(data) > 2 && binary.LittleEndian.Uint16(data[0:2]) == 0x0499
}

//ParseRuuviData parses ruuvidata
func ParseRuuviData(data []byte, macAddress string) (result *SensorData, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			err = errors.New("exception while parsing a potential RuuviTag")
		}
	}()

	// first two bytes are for manufacturer id
	sensorFormat := data[2]
	log.Printf("Ruuvi data with sensor format %d\n", sensorFormat)
	switch sensorFormat {
	case 3:
		// MAC is included v5's payload but not in v3's
		return ParseSensorFormat3(data, macAddress), nil
	case 5:
		return ParseSensorFormat5(data[2:]), nil
	default:
		log.Printf("Unknown sensor format %d", sensorFormat)
	}

	return nil, errors.New("could not parse tag")
}
