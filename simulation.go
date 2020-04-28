package covid19

import (
	"time"

	"github.com/cockroachdb/errors"
)

// Epoch marks the moment that the COVID-19 pandemic first started.
var Epoch = time.Date(2019, time.November, 17, 0, 0, 0, 0, time.UTC)

// A Simulation can simulate a medical scenario with Humans and Viruses.
type Simulation struct {
	humans  []*Human
	viruses []*Virus
	birther *Birther

	date              time.Time
	pendingInfections []infection

	sociability       int // number of humans that one human will be in close contact with every day
	mutationFrequency time.Duration
}

type infection struct {
	Spreader *Human
	Receiver *Human
}

// NewSimulation creates a Simulation.
func NewSimulation(b *Birther, start time.Time) *Simulation {
	return &Simulation{
		birther:           b,
		sociability:       defaultSociability,
		mutationFrequency: defaultMutationFrequency,
		date:              start,
	}
}

const (
	defaultSociability       = 3
	defaultMutationFrequency = 7 * howLongIsADay
)

func (sim *Simulation) addHuman(h *Human) { sim.humans = append(sim.humans, h) }
func (sim *Simulation) addVirus(v *Virus) { sim.viruses = append(sim.viruses, v) }

// Progenerate begins an outbreak with an initial infected human (a "patient
// zero").
func (sim *Simulation) Progenerate(v *Virus) error {
	sim.addVirus(v)

	h, err := sim.birther.Spawn()
	if err != nil {
		return errors.Wrap(err, "spawn Human")
	}
	h.Virus = v

	// Add the infected human into the simulation.
	sim.addHuman(h)
	return nil
}

// Prepopulate prepopulates the simulated environment with n humans.
func (sim *Simulation) Prepopulate(n int) error {
	humans, err := sim.birther.SpawnMany(n)
	if err != nil {
		return err
	}
	sim.humans = append(sim.humans, humans...)
	return nil
}

const howLongIsADay = 24 * time.Hour

// Tick simulates a single day.
func (sim *Simulation) Tick() {
	// Clear pending infections, simulate each human.
	for _, h := range sim.humans {
		sim.tickHuman(h)
	}

	// Apply pending infections.
	for _, infection := range sim.pendingInfections {
		if infection.Receiver.Infected() {
			continue
		}
		infection.Receiver.Virus = infection.Spreader.Virus
	}

	// Increment date, clear pending infections.
	sim.date = sim.date.Add(howLongIsADay)
	sim.pendingInfections = nil
}

func (sim *Simulation) tickHuman(h *Human) {
	if !h.Infected() || h.Dead() {
		return
	}

	// Update human health.
	h.Health -= h.Virus.Lethality
	if h.Health < 0 {
		h.Health = 0
	}

	// Spread the virus!
	sim.spreadVirusOnBehalfOf(h)
}

func (sim *Simulation) spreadVirusOnBehalfOf(h *Human) {
	for i := 0; i < sim.sociability; i++ {
		picked := sim.pickRandomHumanOtherThan(h)
		if picked.Infected() {
			continue
		}

		// If luckySpin is less than the virulence, then that individual should be
		// infected.
		luckySpin := prand.Intn(100) + 1
		if luckySpin <= h.Virus.Virulence {
			sim.pendingInfections = append(sim.pendingInfections, infection{
				Spreader: h,
				Receiver: picked,
			})
		}
	}
}

func (sim *Simulation) pickRandomHumanOtherThan(h *Human) *Human {
	if len(sim.humans) == 1 {
		panic(errors.New("covid19: only one human"))
	}

Top:
	picked := sim.humans[prand.Intn(len(sim.humans))]
	if picked == h {
		goto Top
	}
	return picked
}

// Humans returns the humans present in the simulation.
func (sim *Simulation) Humans() []*Human { return sim.humans }

// Date returns the current date of the simulation.
func (sim *Simulation) Date() time.Time { return sim.date }
