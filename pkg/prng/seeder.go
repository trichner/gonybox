package prng

import (
	"machine"
	"trelligo/pkg/shims/rand"
)

var _ = Seeder(machine.GetRNG)

type Seeder func() (uint32, error)

func New(s Seeder) (*rand.Rand, error) {

	hi, err := s()
	if err != nil {
		return nil, err
	}
	lo, err := s()
	if err != nil {
		return nil, err
	}
	rsrc := rand.NewSource(int64(hi)<<32 | int64(lo))
	return rand.New(rsrc), nil
}

func NewDefault() (*rand.Rand, error) {
	return New(machine.GetRNG)
}
