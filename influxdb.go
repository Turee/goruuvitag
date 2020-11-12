package main

import (
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

var client influxdb2.Client
var writeAPI api.WriteAPI

// InitializeClient is required to open the connection
func InitializeClient() {
	// Create a new client using an InfluxDB server base URL and an authentication token
	// and set batch size to 20
	client = influxdb2.NewClientWithOptions(os.Getenv("INFLUXDB_HOST"), os.Getenv("INFLUXDB_TOKEN"),
		influxdb2.DefaultOptions().SetBatchSize(20))
	// Get non-blocking write client
	writeAPI = client.WriteAPI(os.Getenv("INFLUXDB_ORG"), os.Getenv("INFLUXDB_BUCKET"))
}

// WriteData writes a single point to InfluxDB. Because the client is batched, the writes may not happen
// immediately. Currently no error handling if the client dies for some reason or so. :)
func WriteData(sensorData *SensorData) {
	// create point
	p := influxdb2.NewPoint(
		"system",
		map[string]string{
			"address": *sensorData.MAC,
		},
		map[string]interface{}{
			"temperature":         *sensorData.Temperature,
			"humidity":            *sensorData.Humidity,
			"pressure":            *sensorData.Pressure,
			"txPower":             *sensorData.TxPower,
			"accelerationX":       *sensorData.AccelerationX,
			"accelerationY":       *sensorData.AccelerationY,
			"accelerationZ":       *sensorData.AccelerationZ,
			"movementCounter":     *sensorData.MovementCounter,
			"measurementSequence": *sensorData.MeasurementSequence,
			"acceleration":        *sensorData.Acceleration,
		},
		time.Now())
	// write asynchronously
	writeAPI.WritePoint(p)
}

// CleanUp flushes writes and closes the connection
func CleanUp() {
	// Force all unwritten data to be sent
	writeAPI.Flush()
	// Ensures background processes finishes
	client.Close()
}
