package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/paypal/gatt"
	"github.com/paypal/gatt/examples/option"
)

var isPoweredOn = false
var scanMutex = sync.Mutex{}

func beginScan(d gatt.Device) {
	scanMutex.Lock()
	for isPoweredOn {
		d.Scan(nil, true) //Scan for five seconds and then restart
		time.Sleep(5 * time.Second)
		d.StopScanning()
	}
	scanMutex.Unlock()
}

func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("scanning...")
		isPoweredOn = true
		go beginScan(d)
		return
	case gatt.StatePoweredOff:
		log.Println("REINIT ON POWER OFF")
		isPoweredOn = false
		d.Init(onStateChanged)
	default:
		log.Println("WARN: unhandled state: ", fmt.Sprint(s))
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	if !IsRuuviTag(a.ManufacturerData) {
		return
	}

	fmt.Printf("\nPeripheral ID:%s, NAME:(%s)\n", p.ID(), p.Name())
	// fmt.Println("  TX Power Level    =", a.TxPowerLevel)
	sensorData, err := ParseRuuviData(a.ManufacturerData, p.ID())
	if err != nil {
		log.Fatal(err)
		return
	}

	WriteData(sensorData)
}

func createSysInfoSender() chan struct{} {
	SendSysInfo()
	log.Println("Sent system info")

	sysInfoTicker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-sysInfoTicker.C:
				SendSysInfo()
				log.Println("Sent system info")
			case <-quit:
				sysInfoTicker.Stop()
				return
			}
		}
	}()
	return quit
}

func main() {
	InitializeClient()
	d, err := gatt.NewDevice(option.DefaultClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	stopSysInfo := createSysInfoSender()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		<-stopSysInfo
		CleanUp()
		os.Exit(0)
	}()

	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(onPeriphDiscovered))
	d.Init(onStateChanged)
	select {}
}
