package repository

import (
	"log"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/joelmertanen/goruuvitag/internal/payloadtype"
)

type influxClient struct {
	config         Config
	clientInstance influxdb2.Client
	writeAPI       api.WriteAPI
}

type Config struct {
	Host         string
	Token        string
	Bucket       string
	Organisation string
	Labels       map[string]string
}

func New(config Config) *influxClient {
	return &influxClient{
		config: config,
	}
}

// Open opens the connection
func (cl *influxClient) Open() {
	// Create a new influxClient using an InfluxDB server base URL and an authentication token
	// and set batch size to 20
	cl.clientInstance = influxdb2.NewClientWithOptions(
		cl.config.Host,
		cl.config.Token,
		influxdb2.DefaultOptions().SetBatchSize(20),
	)
	// Get non-blocking write influxClient
	cl.writeAPI = cl.clientInstance.WriteAPI(cl.config.Organisation, cl.config.Bucket)

	// async write errors appear here
	errorsCh := cl.writeAPI.Errors()
	go func() {
		for err := range errorsCh {
			log.Fatalf("Error when writing to influxdb: %s\n", err.Error())
		}
	}()
}

// CleanUp flushes writes and closes the connection
func (cl *influxClient) CleanUp() {
	// Force all unwritten data to be sent
	cl.writeAPI.Flush()
	// Ensures background processes finishes
	cl.clientInstance.Close()
}

func (cl *influxClient) SendSysInfo(payload payloadtype.Payload) {
	point := influxdb2.NewPoint(
		"system",
		map[string]string{},
		payload,
		time.Now())

	// write asynchronously
	cl.writeAPI.WritePoint(point)
}

func (cl *influxClient) Store(label string, payload payloadtype.Payload) {
	point := influxdb2.NewPoint(
		"ruuvitag",
		map[string]string{
			"label": label,
		},
		payload,
		time.Now())

	// write asynchronously
	cl.writeAPI.WritePoint(point)
}
