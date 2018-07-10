#!/usr/bin/env bash

set -e

function rmif() {
	if [ -e "$1" ]; then
		rm -v "$1"
	fi
}
rmif /usr/local/bin/proxy-server
rmif /usr/local/bin/proxy-client

systemctl disable --now proxy-server.service proxy-client.service || true

rmif /etc/systemd/system/proxy-server.service
rmif /etc/systemd/system/proxy-client.service

echo "success... files in /etc/proxy-ports not delete..."
