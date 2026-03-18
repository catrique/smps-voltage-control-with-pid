package engine

import "math"

func RouthHurwitz(p Plant) (float64, string) {

	zeta := (p.R / 2) * math.Sqrt(p.C/p.L)
	if zeta < 1 {
		return zeta, "Subamortecido"
	} else if zeta == 1 {
		return zeta, "Criticamente amortecido"
	} else if zeta > 1 {
		return zeta, "Superamortecido"
	}
	return 0, "Valor inválido"
}
