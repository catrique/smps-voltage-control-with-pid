package engine

import (
	"math/rand/v2"
	"time"
)

type Noise struct {
	ErrorMargin float64
	Generator   *rand.Rand
}

func NewNoise(margin float64) *Noise {
	now := time.Now().UnixNano()
	pcg := rand.NewPCG(uint64(now), uint64(now>>32))

	return &Noise{
		ErrorMargin: margin,
		Generator:   rand.New(pcg),
	}
}

func (n *Noise) GenerateNoise() float64 {
	randomFactor := (n.Generator.Float64() * 2) - 1
	return 1.0 + (n.ErrorMargin * randomFactor)
}
