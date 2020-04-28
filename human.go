package covid19

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
)

// A Human represents a human being.
type Human struct {
	Name   string `json:"name"`
	Gender Gender `json:"gender"`
	Age    int    `json:"age"`

	Virus  *Virus `json:"virus"`
	Health int    `json:"health"` // between 0 and 100
}

// Infected returns true if the Human has been infected by a Virus.
func (h *Human) Infected() bool { return h.Virus != nil }

// Dead returns true if the Human is dead.
func (h *Human) Dead() bool { return h.Health == 0 }

// Status describes the overall condition of the Human.
func (h *Human) Status() string {
	switch {
	case h.Health == 0:
		return "☠️"
	case h.Health < 25:
		return "AHHH MY VITAL ORGANS!!1!"
	case h.Health < 50:
		return "i CANT FEEL my lUNGS"
	case h.Health < 75:
		return "ouch ouch oweee"
	default:
		return "im jus' chillin"
	}
}

// Suffer causes the Human to lose the given amount of health.
func (h *Human) Suffer(amount int) {
	h.Health -= amount
	if h.Health < 0 {
		h.Health = 0
	}
}

// A Birther creates Humans.
type Birther struct {
	client *http.Client
}

// NewBirther creates a Birther that makes HTTP requests using client.
//
// If client is nil, http.DefaultClient will be used.
func NewBirther(client *http.Client) *Birther {
	if client == nil {
		client = http.DefaultClient
	}
	return &Birther{client: client}
}

const (
	randomUserEndpoint = "https://randomuser.me/api/"
	maxSpawnCount      = 5000
)

// SpawnMany spawns n random uninfected Humans.
func (b *Birther) SpawnMany(n int) ([]*Human, error) {
	if n > maxSpawnCount {
		return nil, errors.Newf(
			"covid19: cannot spawn more than %d Humans",
			maxSpawnCount,
		)
	}

	url := fmt.Sprintf("%s?results=%d", randomUserEndpoint, n)
	res, err := b.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var data struct {
		Results []struct {
			Gender Gender `json:"gender"`
			Name   struct {
				First string `json:"first"`
				Last  string `json:"last"`
			} `json:"name"`
			DOB struct {
				Age int `json:"age"`
			} `json:"dob"`
		}
	}
	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, errors.Wrap(err, "decode response")
	}

	humans := make([]*Human, len(data.Results))
	for i, r := range data.Results {
		name := r.Name
		humans[i] = &Human{
			Name:   fmt.Sprintf("%s %s", name.First, name.Last),
			Gender: r.Gender,
			Age:    r.DOB.Age,
			Health: 100,
		}
	}
	return humans, nil
}

// Spawn spawns a random uninfected Human.
func (b *Birther) Spawn() (*Human, error) {
	humans, err := b.SpawnMany(1)
	if err != nil {
		return nil, err
	}
	return humans[0], nil
}
