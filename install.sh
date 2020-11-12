#!/bin/bash -x

echo "Fetching dependencies"
go install

echo "Building binary"
go build

echo "Copying binary"
sudo cp goruuvitag /usr/local/bin
# give capabilities to read bluetooth
sudo setcap 'cap_net_raw,cap_net_admin=eip' /usr/local/bin/goruuvitag

echo "Copying config"
sudo cp goruuvitag.json /usr/local/etc

echo "Copying systemd service file"
sudo cp goruuvitag.service /etc/systemd/system/goruuvitag.service

echo "Enabling systemd service"
systemctl enable goruuvitag
