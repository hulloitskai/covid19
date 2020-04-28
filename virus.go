package covid19

import (
	"math/rand"
	"time"

	"github.com/cockroachdb/errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// A Virus represents a class of SARS-CoV-2 virus.
type Virus struct {
	// Virus strain name, and parent strain name.
	Strain string `json:"strain"`
	Parent string `json:"parent"`

	// Lethality is a value between 0 and 100 that describes how deadly a Virus
	// is.
	Lethality int `json:"lethality"`

	// Virulence is a value between 0 and 100 that describes how infectious a
	// Virus is.
	Virulence int `json:"virulence"`
}

// Mutate creates a new Virus that has mutated from v, with a maximum deviation
// of maxDeviation.
func (v *Virus) Mutate(maxDeviation int) *Virus {
	if err := validation.Validate(maxDeviation, validation.Max(100)); err != nil {
		panic(errors.Wrap(err, "invalid max deviation"))
	}

	var (
		n         = maxDeviation*2 + 1
		lethality = prand.Intn(n) - maxDeviation
		virulence = prand.Intn(n) - maxDeviation
	)

	// Check bounds on lethality and virulence.
	bounded := []*int{&lethality, &virulence}
	for _, n := range bounded {
		switch {
		case *n > 100:
			*n = 100
		case *n < 0:
			*n = 0
		}
	}

	return &Virus{
		Strain:    newStrain(),
		Parent:    v.Strain,
		Lethality: lethality,
		Virulence: virulence,
	}
}

var prand = rand.New(rand.NewSource(time.Now().UnixNano()))

// NewVirus creates a Virus with the given lethality and virulence.
func NewVirus(lethality, virulence int) *Virus {
	values := map[string]int{
		"lethality": lethality,
		"virulence": virulence,
	}
	for name, value := range values {
		if err := validation.Validate(
			value,
			validation.Min(0),
			validation.Max(100),
		); err != nil {
			panic(errors.Wrapf(err, "invalid %s", name))
		}
	}
	return &Virus{
		Strain:    newStrain(),
		Lethality: lethality,
		Virulence: virulence,
	}
}

func newStrain() string { return primitive.NewObjectID().Hex() }
