# Description

Originally a fork of <https://github.com/Turee/goruuvitag>, but ended up rewriting almost everything.
Now with support for RuuviTag's protocol v5 and influxdb2.

This program listens for [RuuviTag](https://tag.ruuvi.com/) devices and posts JSON data to HTTP endpoint specified by HTTP_URL environment variable.

Bluetooth implementation relies on [Paypal's GATT library](https://github.com/paypal/gatt). See [setup](https://github.com/paypal/gatt#setup).

# Requirements

- Recent enough Go (> 1.19) due to module support and tests
- Bluetooth module works at least on Raspbian
- InfluxDB2 - One option is to sign up for the free plan in Influx Cloud

# Setup

## For a local test

```
cp goruuvitag.json.example goruuvitag.json # You might want to check the insides too
sudo hciconfig hci0 down
make local-dev
```

## QA

```
make test
make lint
```

## Install a systemd service

```
./install.sh
```

Licensed under MIT license
