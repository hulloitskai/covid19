package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	"go.stevenxie.me/covid19"
)

func main() {
	if err := exec(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %+v\n", err)
		os.Exit(1)
	}
}

func exec(args []string) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = args

	// Create simulation.
	birther := covid19.NewBirther(nil)
	simulation := covid19.NewSimulation(birther, covid19.Epoch)

	// Prepopulate simulation with uninfected humans.
	if err := simulation.Prepopulate(25); err != nil {
		return errors.Wrap(err, "prepopulate")
	}

	// Define virus and infect patient zero.
	virus := covid19.NewVirus(18, 70)
	if err := simulation.Progenerate(virus); err != nil {
		return errors.Wrap(err, "progenerate")
	}

	// Play simulation, day-by-day.
	for {
		fmt.Printf("[ %s ]\n", simulation.Date().Format("2006-01-02"))
		for _, h := range simulation.Humans() {
			fmt.Printf(
				"%s (🤮: %t, ❤️: %d): %s\n",
				h.Name, h.Infected(), h.Health, h.Status(),
			)
		}
		fmt.Println()
		time.Sleep(5 * time.Second)
		simulation.Tick()
	}
}

func mustEncode(enc *json.Encoder, v interface{}) {
	if err := enc.Encode(v); err != nil {
		panic(err)
	}
}
