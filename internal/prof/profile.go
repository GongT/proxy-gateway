package prof

import (
	"log"
	"runtime"
	"os"
	"runtime/pprof"
	"time"
	"fmt"
)

func Snapshot(point string) {
	runtime.GC()
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	log.Printf("|%s| goroutine: %d, memory: %d bytes (%d obj)", point, runtime.NumGoroutine(), ms.Alloc, ms.HeapObjects)
}

func SnapshotHeap() {
	pp := pprof.Lookup("goroutine")
	pp.WriteTo(os.Stderr, 1)
	f, err := os.Create(fmt.Sprintf("%d.txt", time.Now().Unix()))
	if err == nil {
		defer f.Close()
		pp.WriteTo(f, 1)
		f.WriteString(fmt.Sprintf("goroutine: %d", runtime.NumGoroutine()))
		log.Printf("write to: %d\n", f.Name())
	}
	log.Printf("goroutine: %d", runtime.NumGoroutine())
}
