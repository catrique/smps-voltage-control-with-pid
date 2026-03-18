package main

import (
	"fmt"
	"smps-voltage-control-with-pid/engine"
)

const TimeStep float64 = 0.0001

func main() {
	var plant engine.Plant
	fmt.Println("Informe a tensão de entrada:")
	fmt.Scan(&plant.V_in)
	vinOriginal := plant.V_in
	fmt.Println("Informe a tensão de saída desejada:")
	fmt.Scan(&plant.V_target)
	fmt.Println("Informe a indutância (L) em henrys")
	fmt.Scan(&plant.L)
	fmt.Println("Informe a capacitância (C) em farads")
	fmt.Scan(&plant.C)
	fmt.Println("Informa a resitência (R) em ohms")
	fmt.Scan(&plant.R)

	if plant.V_in == 0 || plant.V_target == 0 || plant.R == 0 || plant.L == 0 || plant.C == 0 {
		fmt.Println("ERRO: L, C e R devem ser maiores que zero para o simulaodr funcionar!")
		return
	}

	z, status := engine.RouthHurwitz(plant)

	fmt.Printf("Zeta: %f, %s \n\n", z, status)

	if z == 0 {
		fmt.Printf("%s", status)
	}

	var meuPID engine.PID
	alfaIdeal := meuPID.AutoSintonizar(plant)
	meuPID.Sintonizar(plant, alfaIdeal)

	fmt.Printf("Sintonia Automática Concluída! Alfa escolhido: %.2f\n", alfaIdeal)
	fmt.Printf("Ganhos Calculados -> Kp: %.2f | Ki: %.2f | Kd: %.2f\n", meuPID.Kp, meuPID.Ki, meuPID.Kd)

	fmt.Println("\nIniciando Simulação...")
	for i := 0.0; i < 50.0; i += TimeStep {
		controle := meuPID.CalcularControlador(plant.V_target, plant.Vout, TimeStep)

		plant.V_in = controle * vinOriginal
		plant.CalcularPasso(TimeStep)

		if int(i/TimeStep)%5000 == 0 {
			fmt.Printf("Tempo: %.3fs | Vout: %.2fV | Duty Cycle: %.1f%%\n", i, plant.Vout, controle*100)
		}

	}

}
