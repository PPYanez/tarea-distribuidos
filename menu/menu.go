package main

import (
	statisticsMenu "distribuidos/tarea-1/menu/getStatistics"
	manageMenu "distribuidos/tarea-1/menu/manageReservation"
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
			"Ingrese una opción: ",
		)

		var option int
		fmt.Scanln(&option)

		if option == 1 {
			manageMenu.ManageReservationMenu()
		}

		if option == 2 {
			statisticsMenu.GetStatisticsMenu()
		}

		if option == 3 {
			keepRunning = false
		}
	}
}
