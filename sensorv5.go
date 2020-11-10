package main

import (
	"fmt"
	"math"
	"strings"
)

func readTemperature(val uint16) *float64 {
	if val == 0x8000 {
		return nil
	}
	result := float64(fromTwosComplement(val, 16)) / 200
	return &result
}

func readPowerInfo(powerInfo uint16) (*int, *int) {
	// Read only first 11 bits out of 16 for battery voltage
	// Make a copy of it first because we need powerInfo later
	batteryVoltage := new(int)
	*batteryVoltage = int(powerInfo) >> 5

	// Error value
	if *batteryVoltage == 2047 {
		batteryVoltage = nil
	} else {
		*batteryVoltage += 1600
	}

	txPower := new(int)
	// Read last five bits for TX Power
	*txPower = (int(powerInfo) & 0b11111)

	// Error value
	if *txPower == 31 {
		txPower = nil
	} else {
		// TX Power's values are in increments of two, with offset of -40
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

func fromTwosComplement(bytes uint16, bits uint16) int {
	mask := uint16(math.Pow(2, float64(bits-1)))
	result := int(bytes & ^mask)

	signBit := -1 * int((bytes & mask))
	return signBit + result
}

// ParseSensorFormat5 parses the given input and returns a struct which contains pointers to the read values. If a value cannot be read,
// its pointer will be nil.
// Pick fields according to
// https://github.com/ruuvi/ruuvi-sensor-protocols/blob/master/dataformat_05.md#data-format-5-protocol-specification-rawv2
func ParseSensorFormat5(data []byte) *SensorData {
	sensorData := SensorData{}
	sensorData.Temperature = readTemperature(bToUint(data[1], data[2]))

	sensorData.Humidity = readHumidity(bToUint(data[3], data[4]))
	sensorData.Pressure = readPressure(bToUint(data[5], data[6]))

	sensorData.AccelerationX = getAcceleration(bToUint(data[7], data[8]))
	sensorData.AccelerationY = getAcceleration(bToUint(data[9], data[10]))
	sensorData.AccelerationZ = getAcceleration(bToUint(data[11], data[12]))

	batteryVoltage, txPower := readPowerInfo(bToUint(data[13], data[14]))
	sensorData.BatteryVoltageMv = batteryVoltage
	sensorData.txPower = txPower
	sensorData.MovementCounter = readMovementCounter(uint8(data[15]))
	sensorData.MeasurementSequence = readMeasurementSequence(bToUint(data[16], data[17]))

	sensorData.MAC = readMAC(data[18:24])

	return &sensorData
}
