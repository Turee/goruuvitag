package sysinfoticker

import (
	"log"
	"time"

	"github.com/elastic/go-sysinfo"
	"github.com/joelmertanen/goruuvitag/internal/payloadtype"
)

type InfluxClient interface {
	SendSysInfo(p payloadtype.Payload)
}

func getSysInfo() payloadtype.Payload {
	host, err := sysinfo.Host()

	if err != nil {
		log.Fatal("Could not get host information")
	}

	process, err := sysinfo.Self()

	if err != nil {
		log.Fatal("Could not get own process data")
	}

	info, err := process.Info()

	if err != nil {
		log.Fatal("Could not get own process information")
	}

	payload := map[string]any{}
	payload["process_uptime"] = time.Since(info.StartTime)

	payload["uptime"] = host.Info().Uptime()

	memoryInfo, err := host.Memory()
	if err == nil {
		payload["total_memory"] = memoryInfo.Total
		payload["used_memory"] = memoryInfo.Used
	}
	return payload
}

func Start(client InfluxClient) chan bool {
	client.SendSysInfo(getSysInfo())
	log.Println("Sent system info")

	sysInfoTicker := time.NewTicker(10 * time.Second)
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-sysInfoTicker.C:
				client.SendSysInfo(getSysInfo())
				log.Println("Sent system info")
			case <-quit:
				sysInfoTicker.Stop()
				return
			}
		}
	}()
	return quit
}
