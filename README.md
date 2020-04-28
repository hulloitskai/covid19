# covid19

_Just me and @OlyaJaworsky fucking around in the midst of a pandemic, nothing
to see here._

<img src="./docs/simulation.gif" width="600px" />

## Usage

```golang
// Create simulation.
birther := covid19.NewBirther(nil)
simulation := covid19.NewSimulation(birther, covid19.Epoch)

// Prepopulate simulation with 25 humans.
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
  fmt.Printf("[%s]\n", simulation.Date().Format("2006-01-02"))
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
```

See [`./cmd/simulation`](./cmd/simulation/main.go) for an example simulation.
