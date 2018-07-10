#!/usr/bin/env bash

set -e

./scripts/build.sh

cp -v dist/proxy-server /usr/local/bin/proxy-server
cp -v dist/proxy-client /usr/local/bin/proxy-client

if [ ! -e /etc/proxy-ports/default.json ] ; then
	mkdir -vp /etc/proxy-ports
	cp configs/default.json /etc/proxy-ports/
	echo '{}' > /etc/proxy-ports/default.json
fi
if [ ! -e /etc/proxy-ports/server.json ] ; then
	cp configs/server.json /etc/proxy-ports/
fi

cp -v init/proxy-client.service init/proxy-server.service /etc/systemd/system/

echo "Install success...
 * you need edit config in \`/etc/proxy-ports\`
 * you need run \`systemctl enable proxy-<client or server>.service\`
"
