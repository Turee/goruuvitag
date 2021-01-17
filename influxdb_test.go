package main

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestHandleEmptyPayload(t *testing.T) {
	testPayload := SensorData{}
	_, _, err := getPayload(&testPayload)
	assert.EqualError(t, err, "no MAC data in payload")
}

func TestHandleMinimalPayload(t *testing.T) {
	// By minimum only MAC is expected
	mac := new(string)
	*mac = "00:01:AA:BB:CC:DD"
	testPayload := SensorData{MAC: mac}
	result, label, _ := getPayload(&testPayload)

	expected := map[string]interface{}{"address": *mac}

	assert.Equal(t, result, expected)
	assert.Equal(t, label, *mac)
}

func TestHandleSensiblePayload(t *testing.T) {
	mac := new(string)
	*mac = "00:01:AA:BB:CC:DD"
	temperature := new(float64)
	*temperature = 14.4
	humidity := new(float64)
	*humidity = 64.12

	testPayload := SensorData{MAC: mac, Temperature: temperature, Humidity: humidity}
	result, label, _ := getPayload(&testPayload)

	expected := map[string]interface{}{"address": *mac, "temperature": *temperature, "humidity": *humidity}

	assert.Equal(t, result, expected)
	assert.Equal(t, label, *mac)
}

func TestLabelHandling(t *testing.T) {
	mac := new(string)
	*mac = "00:01:aa:bb:cc:dd"
	viper.Set("ruuvitag-labels", map[string]interface{}{"00:01:AA:BB:CC:DD": "my kingdom"})

	testPayload := SensorData{MAC: mac}
	_, label, _ := getPayload(&testPayload)

	assert.Equal(t, label, "my kingdom")
}
