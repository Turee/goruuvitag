package sensor

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

func getTotalAcceleration(accX, accY, accZ float64) *float64 {
	acc := math.Sqrt(accX*accX + accY*accY + accZ*accZ)
	return &acc
}

// ParseSensorFormat5 parses the given input and returns a struct which contains pointers to the read values. If a value cannot be read,
// its pointer will be nil.
// Pick fields according to
// https://github.com/ruuvi/ruuvi-sensor-protocols/blob/master/dataformat_05.md#data-format-5-protocol-specification-rawv2
func (b *BeaconData) ReadV5(data []byte) {
	b.Temperature = readTemperature(bToUint(data[1], data[2]))

	b.Humidity = readHumidity(bToUint(data[3], data[4]))
	b.Pressure = readPressure(bToUint(data[5], data[6]))

	b.AccelerationX = getAcceleration(bToUint(data[7], data[8]))
	b.AccelerationY = getAcceleration(bToUint(data[9], data[10]))
	b.AccelerationZ = getAcceleration(bToUint(data[11], data[12]))
	if b.AccelerationX != nil && b.AccelerationY != nil && b.AccelerationZ != nil {
		b.Acceleration = getTotalAcceleration(*b.AccelerationX, *b.AccelerationY, *b.AccelerationZ)
	}

	batteryVoltage, txPower := readPowerInfo(bToUint(data[13], data[14]))
	b.BatteryVoltageMv = batteryVoltage
	b.TxPower = txPower
	b.MovementCounter = readMovementCounter(uint8(data[15]))
	b.MeasurementSequence = readMeasurementSequence(bToUint(data[16], data[17]))

	b.MAC = readMAC(data[18:24])
}
