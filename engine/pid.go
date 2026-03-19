package engine

import "math"

type PID struct {
	Kp, Ki, Kd    float64
	PreviousError float64
	Integral      float64
}

func (c *PID) ComputeControl(targetV float64, outV float64, dt float64) float64 {
	err := targetV - outV

	P := c.Kp * err
	I := c.Ki * c.Integral
	D := c.Kd * (err - c.PreviousError) / dt
	controlOutput := P + I + D

	isSaturatedHigh := controlOutput > 1.0 && err > 0
	isSaturatedLow := controlOutput < 0.0 && err < 0

	if !isSaturatedHigh && !isSaturatedLow {
		c.Integral += err * dt
	}
	c.PreviousError = err

	if controlOutput > 1.0 {
		return 1.0
	} else if controlOutput < 0.0 {
		return 0.0
	} else {
		return controlOutput
	}
}

func (c *PID) Tune(p Plant, alpha float64) {
	c.Kd = (3 * alpha * p.L) - p.R
	if c.Kd < 0 {
		c.Kd = 0
	}

	c.Kp = (3 * (math.Pow(alpha, 2) * p.L * p.C)) - 1
	if c.Kp < 0.1 {
		c.Kp = 0.1
	}

	c.Ki = math.Pow(alpha, 3) * p.L * p.C
	if c.Ki < 0.1 {
		c.Ki = 0.1
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
	return bestAlpha
}
