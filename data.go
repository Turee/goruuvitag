package main

import (
	"bytes"

	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	httpURL = os.Getenv("HTTP_URL")
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

//SensorData to be posted
type SensorDataFormat5 struct {
	Temperature         *float64
	Humidity            *float64
	Pressure            *uint32
	BatteryVoltageMv    *int
	txPower             *int
	AccelerationX       *float64
	AccelerationY       *float64
	AccelerationZ       *float64
	MovementCounter     *uint8
	MeasurementSequence *uint16
	MAC                 *string
}

//SensorFormat5 https://github.com/ruuvi/ruuvi-sensor-protocols/blob/master/dataformat_05.md#data-format-5-protocol-specification-rawv2
type SensorFormat5 struct {
	DataFormat          uint8
	Temperature         uint16
	Humidity            uint16
	Pressure            uint16
	AccelerationX       uint16
	AccelerationY       uint16
	AccelerationZ       uint16
	PowerInfo           uint16
	MovementCounter     uint8
	MeasurementSequence uint16
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
func parseSensorFormat3(data []byte) *SensorData {
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

func readTemperature(val uint16) *float64 {
	if val == 0x8000 {
		return nil
	}
	result := float64(fromTwosComplement(val, 16)) / 200
	return &result
}

func fromTwosComplement(bytes uint16, bits uint16) int {
	mask := uint16(math.Pow(2, float64(bits-1)))
	result := int(bytes & ^mask)

	signBit := -1 * int((bytes & mask))
	return signBit + result
}

func readPowerInfo(powerInfo uint16) (*int, *int) {
	// Read only first 11 bits out of 16 for battery voltage
	// Make a copy of it first because we need powerInfo later
	batteryVoltage := new(int)
	*batteryVoltage = int(powerInfo) >> 5
	if *batteryVoltage == 2047 {
		batteryVoltage = nil
	} else {
		*batteryVoltage += 1600
	}

	txPower := new(int)
	// Read last five bits for TX Power in increments of two
	*txPower = (int(powerInfo) & 0b11111)

	if *txPower == 31 {
		txPower = nil
	} else {
		*txPower = *txPower*2 - 40
	}

	return batteryVoltage, txPower
}

func readMAC(data []byte) *string {
	isValid := false
	for i := 0; i < len(data); i++ {
		if data[i] != 0xFF {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil
	}

	result := make([]string, len(data))

	for i := 0; i < len(data); i++ {
		result[i] = fmt.Sprintf("%X", data[i])
	}

	val := strings.Join(result, ":")
	return &val
}

func getAcceleration(val uint16) *float64 {
	if val == 0x8000 {
		return nil
	}
	result := float64(fromTwosComplement(val, 16)) / 1000
	return &result
}

func readHumidity(val uint16) *float64 {
	if val == 0xFFFF {
		return nil
	}
	result := float64(val) / 400
	return &result
}

func readPressure(val uint16) *uint32 {
	if val == 0xFFFF {
		return nil
	}
	result := uint32(val) + 50000
	return &result
}

func readMeasurementSequence(val uint16) *uint16 {
	if val == 0xFFFF {
		return nil
	}
	return &val
}

func readMovementCounter(val uint8) *uint8 {
	if val == 0xFF {
		return nil
	}
	return &val
}

func bToUint(first byte, second byte) uint16 {
	return uint16(first)<<8 + uint16(second)
}

// https://github.com/ruuvi/ruuvi-sensor-protocols
func parseSensorFormat5(data []byte) *SensorDataFormat5 {
	reader := bytes.NewReader(data)
	result := SensorFormat5{}
	err := binary.Read(reader, binary.BigEndian, &result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", result)
	sensorData := SensorDataFormat5{}
	sensorData.Temperature = readTemperature(bToUint(data[1], data[2]))

	sensorData.Humidity = readHumidity(bToUint(data[3], data[4]))
	sensorData.Pressure = readPressure(bToUint(data[5], data[6]))

	batteryVoltage, txPower := readPowerInfo(result.PowerInfo)
	sensorData.BatteryVoltageMv = batteryVoltage
	sensorData.txPower = txPower
	sensorData.MeasurementSequence = readMeasurementSequence(bToUint(data[16], data[17]))
	sensorData.MovementCounter = readMovementCounter(uint8(data[15]))
	sensorData.AccelerationX = getAcceleration(bToUint(data[7], data[8]))
	sensorData.AccelerationY = getAcceleration(bToUint(data[9], data[10]))
	sensorData.AccelerationZ = getAcceleration(bToUint(data[11], data[12]))

	sensorData.MAC = readMAC(data[18:24])
	fmt.Printf("%+v\n", sensorData)
	fmt.Printf("%+v\n", result)
	return &sensorData
}

func IsRuuviTag(data []byte) bool {
	return binary.LittleEndian.Uint16(data[0:2]) == 0x0499
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

	sensorFormat := data[2]
	fmt.Printf("Ruuvi data with sensor format %d\n", sensorFormat)
	switch sensorFormat {
	case 3:
		sendData(parseSensorFormat3(data))
	case 5:
		parseSensorFormat5(data)
	default:
		fmt.Printf("Unknown sensor format %d", sensorFormat)
	}
}
