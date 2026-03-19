package engine

import (
	"math"
	"time"
)

type PowerSupply struct {
	NominalVin     float64
	RipplePercent  float64
	NoiseGenerator *Noise
	StartTime      time.Time
}

func NewPowerSupply(vin float64, ripple float64, noisePercent float64) *PowerSupply {
	return &PowerSupply{
		NominalVin:     vin,
		RipplePercent:  ripple,
		NoiseGenerator: NewNoise(noisePercent),
		StartTime:      time.Now(),
	}
}

func (ps *PowerSupply) GetVoltage() float64 {
	t := time.Since(ps.StartTime).Seconds()
	ripple := ps.NominalVin * ps.RipplePercent * math.Sin(2*math.Pi*120*t)
	noiseFactor := ps.NoiseGenerator.GenerateNoise()
	return (ps.NominalVin + ripple) * noiseFactor
}
