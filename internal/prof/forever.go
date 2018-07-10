package prof

import (
	"os"
	"syscall"
	"log"
	"os/signal"
)

func RunForever() {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigs
	log.Println("Will quit with signal:", sig)
	os.Exit(1)
}

func Debug() {
	log.Println("debug with SIGUSR1: kill -SIGUSR1", os.Getpid())
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGUSR1)

	for range sigs {
		SnapshotHeap()
	}
}
