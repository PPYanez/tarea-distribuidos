package main

import (
	menu "distribuidos/tarea-1/menu/manageReservation"
	"fmt"
)

func main() {
	keepRunning := true

	for keepRunning {
		fmt.Print(
			"Menu:\n",
			"1. Gestionar reserva\n",
			"2. Obtener estadísticas\n",
			"3. Salir\n",
		)

		var option int
		fmt.Scanln(&option)

		if option == 1 {
			menu.ManageReservation()
		}

		if option == 3 {
			keepRunning = false
		}
	}
}
