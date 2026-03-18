package engine

import "math"

type PID struct {
	Kp, Ki, Kd   float64
	ErroAnterior float64
	Integral     float64
}

func (c *PID) CalcularControlador(v_target float64, v_out float64, dt float64) float64 {
	erro := v_target - v_out

	P := c.Kp * erro
	I := c.Ki * c.Integral
	D := c.Kd * (erro - c.ErroAnterior) / dt
	controlador := P + I + D

	saturadoTopo := controlador > 1.0 && erro > 0
	saturadoFundo := controlador < 0.0 && erro < 0

	if !saturadoTopo && !saturadoFundo {
		c.Integral += erro * dt
	}
	c.ErroAnterior = erro

	if controlador > 1.0 {
		return 1.0
	} else if controlador < 0.0 {
		return 0.0
	} else {
		return controlador
	}
}

func (c *PID) Sintonizar(p Plant, a float64) {
	c.Kd = (3 * a * p.L) - p.R
	if c.Kd < 0 {
		c.Kd = 0
	}

	c.Kp = (3 * (math.Pow(a, 2) * p.L * p.C)) - 1
	if c.Kp < 0.1 {
		c.Kp = 0.1
	}

	c.Ki = math.Pow(a, 3) * p.L * p.C
	if c.Ki < 0.1 {
		c.Ki = 0.1
	}
}
func (c *PID) AutoSintonizar(p Plant) float64 {
	low := 50.0
	high := 8000.0
	melhorAlfa := low

	for i := 0; i < 50; i++ {
		alfaTeste := (low + high) / 2
		testeP := p
		testeP.Vout = 0
		testeP.iL = 0
		tempPID := PID{}
		tempPID.Sintonizar(testeP, alfaTeste)
		overshoot := false
		dt := 0.0001

		for t := 0.0; t < 0.1; t += dt {
			controle := tempPID.CalcularControlador(testeP.V_target, testeP.Vout, dt)
			testeP.V_in = controle * p.V_in
			testeP.CalcularPasso(dt)

			if testeP.Vout > testeP.V_target*1.05 {
				overshoot = true
				break
			}
		}

		if overshoot {
			high = alfaTeste
		} else {
			melhorAlfa = alfaTeste
			low = alfaTeste
		}
	}
	return melhorAlfa
}
