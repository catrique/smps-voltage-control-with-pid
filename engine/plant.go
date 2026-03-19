package engine

type Plant struct {
	L       float64
	R       float64
	C       float64
	Vin     float64
	VTarget float64
	Vout    float64
	iL      float64
}

func (p *Plant) ComputeStep(timeStep float64) {
	diL := ((p.Vin - p.Vout) / p.L) * timeStep
	p.iL += diL
	if p.iL < 0 {
		p.iL = 0
	}

	dVout := ((p.iL - (p.Vout / p.R)) / p.C) * timeStep
	p.Vout += dVout
	if p.Vout < 0 {
		p.Vout = 0
	}
}
