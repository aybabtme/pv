package main

import (
	"github.com/aybabtme/iocontrol"
	"io"
	"log"
	"os"
	"time"

	"github.com/aybabtme/uniplot/spark"
	"github.com/dustin/go-humanize"
)

func main() {
	log.SetPrefix("pv: ")
	log.Printf("streaming!")

	start := time.Now()
	mrw := iocontrol.NewMeasuredReader(os.Stdin)
	pw, pr, sample := iocontrol.Profile(os.Stdout, mrw)

	go func() {
		for range time.NewTicker(5 * time.Second).C {
			profile := sample()
			log.Printf("%s at %sps (read=%v write=%v total=%v)",
				humanize.IBytes(uint64(mrw.Total())),
				humanize.IBytes(mrw.BytesPerSec()),
				profile.WaitRead,
				profile.WaitWrite,
				profile.Total,
			)
		}
	}()

	n, err := io.Copy(pw, spark.ReaderOut(pr, os.Stderr))
	if err != nil {
		log.Fatalf("error after %s: %v", humanize.IBytes(uint64(n)), err)
	}
	log.Printf("done after %s in %v", humanize.IBytes(uint64(n)), time.Since(start))
}
