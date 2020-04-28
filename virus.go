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

	// Check bounds for lethality and virulence.
	if lethality > 100 {
		lethality = 100
	} else if lethality < 0 {
		lethality = 0
	}
	if virulence > 100 {
		virulence = 100
	} else if virulence < 0 {
		virulence = 0
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
	err := validation.Validate(lethality, validation.Min(0), validation.Max(100))
	if err != nil {
		panic(errors.Wrap(err, "invalid lethality"))
	}
	err = validation.Validate(virulence, validation.Min(0), validation.Max(100))
	if err != nil {
		panic(errors.Wrap(err, "invalid virulence"))
	}
	return &Virus{
		Strain:    newStrain(),
		Lethality: lethality,
		Virulence: virulence,
	}
}

func newStrain() string { return primitive.NewObjectID().Hex() }
