package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joelmertanen/goruuvitag/internal/repository"
	"github.com/joelmertanen/goruuvitag/internal/service"
	"github.com/joelmertanen/goruuvitag/internal/sysinfoticker"

	"github.com/spf13/viper"
)

func readConfig() repository.Config {
	viper.SetConfigName("goruuvitag")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/usr/local/etc")
	if err := viper.ReadInConfig(); err != nil {
		// no point to continue without a config
		panic(err)
	}

	return repository.Config{
		Host:         viper.GetString("influxdb.host"),
		Token:        viper.GetString("influxdb.token"),
		Bucket:       viper.GetString("influxdb.bucket"),
		Organisation: viper.GetString("influxdb.org"),
		Labels:       viper.GetStringMapString("ruuvitag-labels"),
	}
}

type Scan interface {
	Start() error
}

func main() {
	config := readConfig()
	influxClient := repository.New(config)
	influxClient.Open()

	stopSysInfo := sysinfoticker.Start(influxClient)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs

		log.Println("Shutting down...")
		stopSysInfo <- true
		if err := influxClient.Close(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}()

	svc := service.New(config.Labels, influxClient)
	if err := svc.Start(); err != nil {
		panic(err)
	}

	// run until os.Exit gets called in the signal handler
	select {}
}
