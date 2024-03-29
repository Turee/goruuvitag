
# This repository has been deprecated!

Check out this fork instead:

https://github.com/joelmertanen/goruuvitag


[![Build Status](https://travis-ci.org/Turee/goruuvitag.svg?branch=master)](https://travis-ci.org/Turee/goruuvitag)

# Description

This program listens for [RuuviTag](https://tag.ruuvi.com/) devices and posts JSON data to HTTP endpoint specified by HTTP_URL environment variable.

Bluetooth implementation relies on [Paypal's GATT library](https://github.com/paypal/gatt). See [setup](https://github.com/paypal/gatt#setup).

Currently supports [ruuvi sensor protocol 3](https://github.com/ruuvi/ruuvi-sensor-protocols). 

JSON data looks like this:

```
{
	"Temp": 1.25, <- temperature celcius
	"Humidity": 94.5, <- relative humidity (%)
	"Pressure": 98760, <- atmospheric pressure (Pa)
	"Battery": 2845, <- battery voltage (mV)
	"Address": "D3:C7:D0:2E:14:D7", <- unique device address
	"AccelerationX": -74, <-- acceleration (milli g)
	"AccelerationY": -156,
	"AccelerationZ": 976,
	"TimeStamp": "2017-12-30T15:02:44.3560173+02:00"
}
```
I currently use this project to record sensor data from Ruuvi Tags around my apartment. The data is posted straight to ElasticSearch index.


# ElasticSearch + Kibana quickstart

Start ELK stack using docker

```
$ sudo docker run -p 5601:5601 -p 9200:9200  -p 5044:5044 --restart unless-stopped \
    -v elk-data:/var/lib/elasticsearch --name elk sebp/elk
```

Access Kibana on http://localhost:5601 , navigate to developer tools.

Create index by executing following in developer tools:
```
PUT ruuvi
{
  "mappings": {
    "data_point": {
      "properties": {
        "Address": {
          "type": "keyword"
        },
        "TimeStamp": {
          "type": "date"
        },
        "Humidity": {
          "type": "float"
        },
        "Pressure": {
          "type": "float"
        },
        "Battery": {
          "type": "float"
        },
        "Temp": {
          "type": "float"
        },
        "AccelerationX": {
          "type": "integer"
        },
        "AccelerationY": {
          "type": "integer"
        },
        "AccelerationZ": {
          "type": "integer"
        }
      }
    }
  }
}
```
Start posting data.

```
$ export HTTP_URL=http://<elasticsearch ip>:9200/ruuvi/data_point
$ ./goruuvitag
```

Explore your data using Kibana :).

# TODO

- [x] Listen and parse ruuvi data
- [x] ElasticSearch support
- [ ] [InfluxDB support](https://docs.influxdata.com/influxdb/v1.4/guides/writing_data/)
- [x] Precompiled binaries

# License

Copyright (c) Turkka Mannila

Licensed under MIT license
