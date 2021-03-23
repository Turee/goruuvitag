#!/bin/bash -x

#echo "Fetching dependencies"
#go install
#
#echo "Building binary"
#rm -f goruuvitag
#go build

echo "Extracting binaries"
tar xvvzf binaries.tar.gz

echo "Stop current service"
sudo systemctl stop goruuvitag

echo "Copying binary"
sudo cp output/linux/arm/goruuvitag /usr/local/bin
# give capabilities to read bluetooth
sudo setcap 'cap_net_raw,cap_net_admin=eip' /usr/local/bin/goruuvitag

echo "Copying config"
sudo cp goruuvitag.json /usr/local/etc

echo "Copying systemd service file"
sudo cp goruuvitag.service /etc/systemd/system/goruuvitag.service

echo "Enabling systemd service"
sudo systemctl enable goruuvitag
sudo systemctl restart goruuvitag
