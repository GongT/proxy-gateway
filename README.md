# Universal Reverse Tunnel Gateway

## With "normal" tunnel (eg: ssh -L), you can do:
```text
                          +----------------------+                  +---------------------+       +------------------+
                          |Local computer        |     Internet     |Server with public IP|       |                  |
                          |                      |      Tunnel      |                     |       | Target server    |
+----------------+        |     Tunnel client  +-+------------------+->  Tunnel server    |       |                  |
| Application    |        |                                                               |       |                  |
|                +------->O-------------------------------------------------------------->O+----->O
+----------------+        |                                                               |       |
                          |                    +-+------------------+->                   |       |
                          |                      |                  |                     |       |                  |
                          +----------------------+                  +---------------------+       +------------------+
```

## With "REVERSE" tunnel (eg: ssh -R), You can do: 
```text
+------------------+        +----------------------+                  +---------------------+
|                  |        |Local computer        |     Internet     |Server with public IP|
| Target           |        |                      |      Tunnel      |                     |
| Application      |        |     Tunnel client  +-+------------------+->  Tunnel server    |       +----------------+
|                  |        |                                                               |       | Application    |
|                  O<------+O<-------------------------------------------------------------+O<------+                |
|                  |        |                                                               |       +----------------+
|                  |        |                    +-+------------------+->                   |
|                  |        |                      |                  |                     |
+------------------+        +----------------------+                  +---------------------+
```

If you use `ssh -R` or `ssh -X`, you will got an super secure, and **super slow** connection.    
I promise you do not need that secure.    
There is so many "normal" proxy server available, but no reverse one.   
So I write one.

# Build & install
dependencies must install before compile: // TODO

```bash
./scripts/build.sh
./scripts/install.sh
```

# Server
get help:
> proxy-server -h

server will listen on port 55600, and not able to change.    
currently you can use systemd socket activation or firewall rule to change port.

# Client
```
./proxy-client -c /path/to/config.json
```

get help:
> proxy-client -h

## Multiple Port Client
```
./proxy-client -c /path/to/config.json
```

will send sd_notify when `--systemd` is set.

### Config File Structure
Client and server use same config file. (but the server only use "kcptun" config)

```json
{
	// optional
	//   * if used, connection will encrypted with "enough secure" 3-DES algorithm, and kcp will active. 
	//   * if not used, transfer will be plaintext, kcp will disable.
	// I strongly recommend use this, it's really super fast! (but wast some bandwidth), see xtaci/kcptun for more.
	"kcptun": "your-password",
	// required, where to connect, server ip or hostname. (port is optional)
	"server": "1.2.3.4",
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
