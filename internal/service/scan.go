package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/joelmertanen/goruuvitag/internal/payloadtype"
	"github.com/joelmertanen/goruuvitag/internal/sensor"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

type InfluxClient interface {
	Store(label string, payload payloadtype.Payload)
}

type scan struct {
	isPoweredOn bool
	scanMutex   sync.Mutex
	labels      map[string]string
	repo        InfluxClient
}

func New(labels map[string]string, repo InfluxClient) *scan {
	return &scan{
		labels: labels,
		repo:   repo,
	}
}

func (s *scan) Start() error {
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open bluetooth device, err: %s\n", err)
	}

	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(s.onPeripheralDiscovered))
	if err := d.Init(s.onStateChanged); err != nil {
		return err
	}
	return nil
}

func (s *scan) onPeripheralDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if !sensor.IsRuuviTag(a.ManufacturerData) {
		return
	}

	log.Printf("Peripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())

	sensorData := sensor.New()
	if err := sensorData.Parse(a.ManufacturerData, p.ID()); err != nil {
		log.Fatal(err)
	}

	d := sensorData.GetData()

	storable, err := s.toStorable(d)
	if err != nil {
		panic(err)
	}

	// everything is lower cased for configs
	label, exists := s.labels[strings.ToLower(*d.MAC)]
	if !exists {
		label = *d.MAC
	}

	s.repo.Store(label, *storable)
}

func (s *scan) toStorable(sensorData sensor.BeaconData) (*map[string]any, error) {
	if sensorData.MAC == nil {
		return nil, errors.New("no MAC data in payload")
	}

	readValues := map[string]interface{}{}

	if sensorData.Temperature != nil {
		readValues["temperature"] = *sensorData.Temperature
	}
	if sensorData.BatteryVoltageMv != nil {
		readValues["batteryMv"] = *sensorData.BatteryVoltageMv
	}
	if sensorData.Humidity != nil {
		readValues["humidity"] = *sensorData.Humidity
	}
	if sensorData.Pressure != nil {
		readValues["pressure"] = *sensorData.Pressure
	}
	if sensorData.TxPower != nil {
		readValues["txPower"] = *sensorData.TxPower
	}
	if sensorData.AccelerationX != nil {
		readValues["accelerationX"] = *sensorData.AccelerationX
	}
	if sensorData.AccelerationY != nil {
		readValues["accelerationY"] = *sensorData.AccelerationY
	}
	if sensorData.AccelerationZ != nil {
		readValues["accelerationZ"] = *sensorData.AccelerationZ
	}
	if sensorData.MovementCounter != nil {
		readValues["movementCounter"] = *sensorData.MovementCounter
	}
	if sensorData.MeasurementSequence != nil {
		readValues["measurementSequence"] = *sensorData.MeasurementSequence
	}
	if sensorData.Acceleration != nil {
		readValues["acceleration"] = *sensorData.Acceleration
	}
	if sensorData.MAC != nil {
		readValues["address"] = *sensorData.MAC
	}

	return &readValues, nil
}

func (s *scan) onStateChanged(d gatt.Device, state gatt.State) {
	log.Println("State:", s)
	switch state {
	case gatt.StatePoweredOn:
		log.Println("Scanning...")
		s.isPoweredOn = true
		go s.beginScan(d)
		return
	case gatt.StatePoweredOff:
		log.Println("REINIT ON POWER OFF")
		s.isPoweredOn = false
		if err := d.Init(s.onStateChanged); err != nil {
			panic(err)
		}
	default:
		log.Println("WARN: unhandled state: ", fmt.Sprint(s))
	}
}

func (s *scan) beginScan(d gatt.Device) {
	s.scanMutex.Lock()
	for s.isPoweredOn {
		d.Scan(nil, true) //Scan for five seconds and then restart
		time.Sleep(5 * time.Second)
		d.StopScanning()
	}
	s.scanMutex.Unlock()
}
