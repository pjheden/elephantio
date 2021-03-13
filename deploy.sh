#!/bin/bash

# Help us not have to start over when connected to the wrong wifi
IP=${1:-192.168.1.79}
echo "pinging $IP"
if ! ping -c1 $IP &>/dev/null; then
	echo "maybe switch wifi networks; trying once more in 5s (^C now to give up)"
	sleep 5
	if ! ping -c1 $IP &>/dev/null; then
		echo "not likely to connect on $IP"
		exit 2
	fi;
fi;
echo "PI is reachable"

# Build all the binaries
prg=${2:-elephantio-runner}
GOOS=linux GOARCH=arm GOARM=5 go build -o $prg ./cmd/$prg/

# Move all the binaries and make them executable
chmod +x $prg
scp $prg pi@$IP:

# Remove the local binaries
rm $prg