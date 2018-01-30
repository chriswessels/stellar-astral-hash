package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/paulbellamy/ratecounter"
	"github.com/stellar/go/keypair"
)

var version = "develop"

var numWorkersFlag = flag.Int("workers", 0, "Number of concurrent workers - defaults to number of CPU cores detected")
var suffixFlag = flag.String("suffix", "", "Desired vanity suffix for Address")

func main() {
	log.Printf("astral-hash %s", version)

	flag.Parse()
	numWorkers := *numWorkersFlag
	suffix := *suffixFlag

	if suffix == "" {
		log.Fatalln("You must set your desired address suffix. See --help")
	}
	if numWorkers == 0 {
		numWorkers = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(numWorkers)

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
