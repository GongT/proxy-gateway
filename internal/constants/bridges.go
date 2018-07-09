package constants

var OpenPortMap = map[string]string{
	// listen to tcp port
	":12123": "127.0.0.1:22",
	":6013": "shabao-work:6000",
	":7070": "10.0.6.2:1080",
}

var ClosePortMap = map[string]string{
	// listen to (abstract) unix-socket
	"@/tmp/.X11-unix/X13": "shabao-work:6000",
	"/tmp/.X11-unix/X13":  "shabao-work:6000",
}

