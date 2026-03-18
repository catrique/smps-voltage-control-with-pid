package engine

type Plant struct {
	L        float64
	R        float64
	C        float64
	V_in     float64
	V_target float64
	Vout     float64
	iL       float64
}

func (p *Plant) CalcularPasso(timeStep float64) {
	d_iL := ((p.V_in - p.Vout) / p.L) * timeStep
	p.iL += d_iL
	if p.iL < 0 {
		p.iL = 0
	}
	d_Vout := ((p.iL - (p.Vout / p.R)) / p.C) * timeStep
	p.Vout += d_Vout
	if p.Vout < 0 {
		p.Vout = 0
	}
}
