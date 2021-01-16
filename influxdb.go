package main

import (
	"errors"
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
	viper.SetConfigName("goruuvitag")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/usr/local/etc")
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

func getPayload(sensorData *SensorData) (map[string]interface{}, string, error) {
	labels := viper.GetStringMapString("ruuvitag-labels")
	if sensorData.MAC == nil {
		return map[string]interface{}{}, "", errors.New("No MAC data in payload")
	}

	// everything is lower cased for viper configs
	label, exists := labels[strings.ToLower(*sensorData.MAC)]
	if !exists {
		label = *sensorData.MAC
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

	return readValues, label, nil
}

// WriteData writes a single point to InfluxDB. Because the client is batched, the writes may not happen
// immediately. Currently no error handling if the client dies for some reason or so. :)
func WriteData(sensorData *SensorData) {
	// get payload
	payload, label, err := getPayload(sensorData)
	if err != nil {
		return
	}

	point := influxdb2.NewPoint(
		"system",
		map[string]string{
			"label": label,
		},
		payload,
		time.Now())

	// write asynchronously
	writeAPI.WritePoint(point)
}

// CleanUp flushes writes and closes the connection
func CleanUp() {
	// Force all unwritten data to be sent
	writeAPI.Flush()
	// Ensures background processes finishes
	client.Close()
}
