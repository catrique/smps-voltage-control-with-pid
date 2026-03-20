package engine

import "math"

type PID struct {
	Kp, Ki, Kd    float64
	PreviousError float64
	Integral      float64
	LastD         float64
	TauD          float64
}

func (c *PID) ComputeControl(targetV float64, outV float64, dt float64) float64 {
	err := targetV - outV

	P := c.Kp * err

	rawD := -(outV - c.PreviousError) / dt
	c.PreviousError = outV

	tauD := c.TauD
	if tauD <= 0 && c.Kp > 0 && c.Kd > 0 {
		tauD = c.Kd / (c.Kp * 10.0)
	}
	var alphaF float64
	if tauD > 0 {
		alphaF = dt / (tauD + dt)
	} else {
		alphaF = 1.0
	}
	c.LastD = c.LastD + alphaF*(rawD-c.LastD)
	D := c.Kd * c.LastD

	c.Integral += c.Ki * err * dt
	I := c.Integral

	controlOutput := P + I + D
	outputSat := controlOutput
	if controlOutput > 1.0 {
		outputSat = 1.0
	} else if controlOutput < 0.0 {
		outputSat = 0.0
	}

	if outputSat != controlOutput {
		kb := math.Sqrt(c.Ki * c.Kp)
		if kb < c.Ki {
			kb = c.Ki
		}
		c.Integral += kb * (outputSat - controlOutput) * dt
	}

	return outputSat
}

func (c *PID) Tune(p Plant, alpha float64) {
	LC := p.L * p.C
	if LC <= 0 {
		return
	}

	c.Kp = 2*alpha*LC - 1
	if c.Kp < 0.05 {
		c.Kp = 0.05
	}

	c.Ki = alpha * alpha * LC
	if c.Ki < 0.5 {
		c.Ki = 0.5
	}

	rawKd := (2*alpha*p.L - p.R) / (alpha * alpha * LC)
	if rawKd < 0 {
		rawKd = 0
	}
	maxKd := c.Kp / (5 * alpha)
	if rawKd > maxKd {
		rawKd = maxKd
	}
	c.Kd = rawKd

	if c.Kp > 0 && c.Kd > 0 {
		c.TauD = c.Kd / (c.Kp * 10.0)
	}
}

func (c *PID) AutoTune(p Plant) float64 {
	low := 50.0
	high := 8000.0
	bestAlpha := low

	for i := 0; i < 50; i++ {
		testAlpha := (low + high) / 2
		testPlant := p
		testPlant.Vout = 0
		testPlant.iL = 0
		tempPID := PID{}
		tempPID.Tune(testPlant, testAlpha)
		hasOvershoot := false
		dt := 0.0001

		for t := 0.0; t < 0.1; t += dt {
			control := tempPID.ComputeControl(testPlant.VTarget, testPlant.Vout, dt)
			testPlant.Vin = control * p.Vin
			testPlant.ComputeStep(dt)

			if testPlant.Vout > testPlant.VTarget*1.05 {
				hasOvershoot = true
				break
			}
		}

		if hasOvershoot {
			high = testAlpha
		} else {
			bestAlpha = testAlpha
			low = testAlpha
		}
	}

	c.Tune(p, bestAlpha)
	return bestAlpha
}
