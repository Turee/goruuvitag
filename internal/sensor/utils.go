package sensor

import (
	"encoding/binary"
)

// IsRuuviTag A helper to check if the manufacturer id of a ble advertisement matches Ruuvi's
func IsRuuviTag(data []byte) bool {
	return len(data) > 2 && binary.LittleEndian.Uint16(data[0:2]) == 0x0499
}
