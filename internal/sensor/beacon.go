package sensor

import (
	"errors"
	"log"
)

type Beacon interface {
	Parse(data []byte, macAddress string) (err error)
	GetData() BeaconData
}

type BeaconData struct {
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

func New() Beacon {
	return &BeaconData{}
}

func (b *BeaconData) GetData() BeaconData {
	return *b
}

func (b *BeaconData) Parse(data []byte, macAddress string) (err error) {
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
		b.ReadV3(data, macAddress)
	case 5:
		b.ReadV5(data[2:])
	default:
		log.Printf("Unknown sensor format %d", sensorFormat)
	}

	return nil
}
