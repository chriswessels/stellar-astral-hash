package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/paulbellamy/ratecounter"
	"github.com/stellar/go/keypair"
)

var version = "develop"
var desc = "A simple, high-throughput vanity address scanner for Stellar Accounts."

var suffixFlag = flag.String("suffix", "", "Desired vanity suffix for Address - required flag")
var numWorkersFlag = flag.Int("workers", 0, "Number of concurrent workers - defaults to number of CPU cores detected")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "astral-hash (%s)\n%s\n\nUsage:\n", version, desc)
		flag.PrintDefaults()
	}
	flag.Parse()
	numWorkers := *numWorkersFlag
	suffix := strings.ToUpper(*suffixFlag)

	if suffix == "" {
		flag.Usage()
		os.Exit(2)
	}
	if numWorkers == 0 {
		numWorkers = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(numWorkers)

	log.Printf("astral-hash (%s)", version)
	result := make(chan *keypair.Full)
	rateCounter := ratecounter.NewRateCounter(1 * time.Second)
	for i := 0; i < numWorkers; i++ {
		go func(i int) {
			for true {
				kp, err := keypair.Random()
				rateCounter.Incr(1)
				if err != nil {
					log.Fatalf("%s", err)
				}
				address := kp.Address()
				if strings.HasSuffix(address, suffix) {
					fmt.Printf("\r")
					log.Printf("Worker %d has found a matching pair", i)
					result <- kp
					break
				}
			}
		}(i)
	}
	log.Printf("Started %d workers to find match for:", numWorkers)
	log.Printf("Address suffix: %s\n\n", suffix)

	tick := time.Tick(1 * time.Second)
	for {
		select {
		case <-tick:
			fmt.Printf("\rScanning %d key pairs per second across %d workers...", rateCounter.Rate(), numWorkers)
			// return
		case match := <-result:
			log.Printf("Public Address: %s", match.Address())
			log.Printf("Private Seed: %s", match.Seed())
			log.Printf("Please keep your seed safe!")
			return
		}
	}
}
