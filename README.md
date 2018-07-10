# Universal Reverse Tunnel Gateway

# Server

# Client
## Multiple Port Client
```
./proxy-server </path/to/config.json> [--systemd]
```

will send sd_notify when `--systemd` is set.

### Config File Structure
```json
{
	// required, where to connect, server ip or hostname. (port is optional)
	"server": "1.2.3.4:5678",
	// required, array of ports to open, see below
	"openPorts": [{ /*...*/ }],
	// optional, use a socks5/http proxy to handle all data transfer, FALSE will also ignore any environment variables.
	"proxy": "socks5://127.0.0.1:1080"
}
```

### Config Section Example:
On remote machine, open port `3306` listen only `lo` interface. Any connection to it will passthru to `/var/run/mysql.sock` on `current` machine.
```json
{
	"listen": {
		"type": "tcp",
		"address": "127.0.0.1:3306"
	},
	"connect": {
		"type": "unix",
		"address": "/var/run/mysql.sock"
	}
}
```
On remote machine, create an abstract unix socket named `/tmp/.X11-unix/X3`, and any program using `DIPLAY=:3` will show window on current machine's `:0` display. 
```json
{
	"listen": {
		"type": "unix",
		"address": "@/tmp/.X11-unix/X3"
	},
	"connect": {
		"type": "unix",
		"address": "@/tmp/.X11-unix/X0"
	}
}
```

## Client Work With 
