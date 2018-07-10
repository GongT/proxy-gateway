package systemd

import (
	"log"
	"github.com/coreos/go-systemd/daemon"
	"flag"
)

var systemdEnabled bool

func ConfigEnabled() bool {
	return systemdEnabled
}

func init() {
	flag.BoolVar(&systemdEnabled, "systemd", false, "send sd_notify when run from systemdEnabled.")
	flag.BoolVar(&systemdEnabled, "S", false, "send sd_notify when run from systemdEnabled. (short hand)")
}

func Status(status string) {
	if systemdEnabled {
		ok, err := daemon.SdNotify(false, "STATUS="+status)
		if !ok {
			if err != nil {
				log.Fatal("failed to notify systemd: ", err)
			} else {
				log.Println("systemd not used, or service Type!=notify")
			}
		}
	}
}

func Ready(status string) {
	Status(status)
	if systemdEnabled {
		ok, err := daemon.SdNotify(false, daemon.SdNotifyReady)
		if !ok {
			if err != nil {
				log.Fatal("failed to notify systemd: ", err)
			} else {
				log.Println("systemd not used, or service Type!=notify")
			}
		}
	}
}
