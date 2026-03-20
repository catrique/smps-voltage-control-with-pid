package main

import (
	"fmt"
	"os"
	"os/exec"
	"smps-voltage-control-with-pid/engine"
	"strconv"
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

	var pid engine.PID
	idealAlpha := pid.AutoTune(plant)
	pid.Tune(plant, idealAlpha)

	fmt.Printf("Sintonia Automática Concluída! Alfa escolhido: %.2f\n", idealAlpha)
	fmt.Printf("Ganhos Calculados -> Kp: %.2f | Ki: %.2f | Kd: %.2f\n", pid.Kp, pid.Ki, pid.Kd)
	fmt.Println("\nIniciando interface gráfica e simulação...")

	cmd := exec.Command("love", "gui",
		strconv.FormatFloat(plant.VTarget, 'f', 4, 64),
	)

	luaStdin, _ := cmd.StdinPipe()
	cmd.Start()

	powerSupply := engine.NewPowerSupply(originalVin, 0.05)

	noiseChannel := make(chan float64, 10)
	go func() {
		n := engine.NewNoise(0.02)
		for {
			noiseChannel <- n.GenerateNoise()
		}
	}()

	dataGraph := make(chan string, 500)

	go func() {
		for msg := range dataGraph {
			fmt.Fprint(luaStdin, msg)
		}
	}()

	for i := 0.0; i < 50.0; i += TimeStep {
		controlOutput := pid.ComputeControl(plant.VTarget, plant.Vout, TimeStep)
		noiseFactor := <-noiseChannel
		currentGridVin := powerSupply.GetVoltage(noiseFactor)
		plant.Vin = currentGridVin * controlOutput
		plant.ComputeStep(TimeStep)

		if int(i/TimeStep)%10 == 0 {
			// fmt.Printf("Tempo: %.3fs | Vin Rede: %.2fV | Vout: %.2fV | Duty Cycle: %.1f%%\n",
			// 	i, currentGridVin, plant.Vout, controlOutput*100)
			line := fmt.Sprintf("%.4f,%.4f,%.4f\n", plant.Vout, currentGridVin, i)
			dataGraph <- line
			// fmt.Fprintf(luaStdin, "%.4f,%.4f,%.4f\n", plant.Vout, currentGridVin, i)
			os.Stdout.Sync()
		}
	}
	cmd.Wait()
}
