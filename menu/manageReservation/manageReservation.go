package menu

import (
	"fmt"
)

func ManageReservation() {
	var keepRunning = true

	for keepRunning {
		fmt.Print(
			"Submenu:\n",
			"1. Crear reserva\n",
			"2. Obtener reserva\n",
			"3. Modificar reserva\n",
			"4. Volver al menú principal\n",
		)

		var option int
		fmt.Scanln(&option)

		if option == 1 {
			var fechaIda string
			var fechaVuelta string
			var origen string
			var destino string
			var cantidadPasajeros int

			fmt.Print("Ingrese la fecha de ida: ")
			fmt.Scanln(&fechaIda)
			fmt.Print("Ingrese la fecha de vuelta: ")
			fmt.Scanln(&fechaVuelta)
			fmt.Print("Ingrese el origen: ")
			fmt.Scanln(&origen)
			fmt.Print("Ingrese el destino: ")
			fmt.Scanln(&destino)
			fmt.Print("Ingrese la cantidad de pasajeros: ")
			fmt.Scanln(&cantidadPasajeros)

			createReservation(fechaIda, fechaVuelta, origen, destino, cantidadPasajeros)
		}

		if option == 2 {
			var pnr string
			var apellido string

			fmt.Print("Ingrese el PNR: ")
			fmt.Scanln(&pnr)

			fmt.Print("Ingrese el apellido: ")
			fmt.Scanln(&apellido)

			getReservation(pnr, apellido)
		}

		if option == 4 {
			keepRunning = false
		}
	}
}
