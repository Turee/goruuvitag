package main

import (
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/spf13/viper"
)

var client influxdb2.Client
var writeAPI api.WriteAPI

// InitializeClient is required to open the connection
func InitializeClient() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		// no point to continue without a config
		panic(err)
	}

	// Create a new client using an InfluxDB server base URL and an authentication token
	// and set batch size to 20
	client = influxdb2.NewClientWithOptions(viper.GetString("influxdb.host"), viper.GetString("influxdb.token"),
		influxdb2.DefaultOptions().SetBatchSize(20))
	// Get non-blocking write client
	writeAPI = client.WriteAPI(viper.GetString("influxdb.org"), viper.GetString("influxdb.bucket"))
}

// WriteData writes a single point to InfluxDB. Because the client is batched, the writes may not happen
// immediately. Currently no error handling if the client dies for some reason or so. :)
func WriteData(sensorData *SensorData) {
	labels := viper.GetStringMapString("ruuvitag-labels")
	// everything is lower cased for viper configs
	label, exists := labels[strings.ToLower(*sensorData.MAC)]
	if !exists {
		label = *sensorData.MAC
	}
	// create point
	p := influxdb2.NewPoint(
		"system",
		map[string]string{
			"label": label,
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
			"address":             *sensorData.MAC,
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
