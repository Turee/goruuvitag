# Description

This program listens for [RuuviTag](https://tag.ruuvi.com/) devices and posts JSON data to HTTP endpoint specified by HTTP_URL environment variable.

Currently support [ruuvi sensor protocol 3](https://github.com/ruuvi/ruuvi-sensor-protocols).

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

# ElasticSearch quickstart

Create index
```
PUT ruuvi/_mapping
{
  "type": {
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

# TODO

- [ ] [InfluxDB support](https://docs.influxdata.com/influxdb/v1.4/guides/writing_data/)
- [ ] Precompiled binaries

# License

Copyright (c) Turkka Mannila

Licensed under MIT license
