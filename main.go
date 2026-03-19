package main

import (
	"fmt"
	"smps-voltage-control-with-pid/engine"
)

const TimeStep float64 = 0.0001

func main() {
	var plant engine.Plant

	fmt.Println("Informe a tensão de entrada (Vin):")
	fmt.Scan(&plant.Vin)
	originalVin := plant.Vin

	fmt.Println("Informe a tensão de saída desejada (Vtarget):")
	fmt.Scan(&plant.VTarget)

	fmt.Println("Informe a indutância (L) em Henrys:")
	fmt.Scan(&plant.L)

	fmt.Println("Informe a capacitância (C) em Farads:")
	fmt.Scan(&plant.C)

	fmt.Println("Informe a resistência (R) em Ohms:")
	fmt.Scan(&plant.R)

	if plant.Vin == 0 || plant.VTarget == 0 || plant.R == 0 || plant.L == 0 || plant.C == 0 {
		fmt.Println("ERRO: L, C e R devem ser maiores que zero para o simulador funcionar!")
		return
	}

	z, status := engine.RouthHurwitz(plant)

	fmt.Printf("Zeta: %f, %s \n\n", z, status)

	if z == 0 {
		fmt.Printf("%s", status)
	}

	var pid engine.PID
	idealAlpha := pid.AutoTune(plant)
	pid.Tune(plant, idealAlpha)

	fmt.Printf("Sintonia Automática Concluída! Alfa escolhido: %.2f\n", idealAlpha)
	fmt.Printf("Ganhos Calculados -> Kp: %.2f | Ki: %.2f | Kd: %.2f\n", pid.Kp, pid.Ki, pid.Kd)

	fmt.Println("\nIniciando Simulação...")

	powerSupply := engine.NewPowerSupply(originalVin, 0.05, 0.01)

	for i := 0.0; i < 50.0; i += TimeStep {
		controlOutput := pid.ComputeControl(plant.VTarget, plant.Vout, TimeStep)
		currentGridVin := powerSupply.GetVoltage()
		plant.Vin = currentGridVin * controlOutput
		plant.ComputeStep(TimeStep)
		if int(i/TimeStep)%500 == 0 {
			fmt.Printf("Tempo: %.3fs | Vin Rede: %.2fV | Vout: %.2fV | Duty Cycle: %.1f%%\n",
				i, currentGridVin, plant.Vout, controlOutput*100)
		}
	}
}
