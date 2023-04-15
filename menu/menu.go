package main

import (
	"context"
	"distribuidos/tarea-1/db"
	"distribuidos/tarea-1/models"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	start()
}

func start() {
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
			manageReservation()
		}

		if option == 3 {
			keepRunning = false
		}
	}
}

func manageReservation() {
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

func getReservation(pnr string, apellido string) {
	// Obtener reserva
	client := db.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	reservaCollection := client.Database("distribuidos").Collection("Reservas")
	var reserva models.Reserva

	filter := bson.M{"pnr": pnr, "pasajeros.apellido": apellido}
	err := reservaCollection.FindOne(ctx, filter).Decode(&reserva)

	if err != nil {
		fmt.Println("No se encontró la reserva")
		return
	}

	// Mostrar detalles de vuelos de ida y vuelta
	for i, vuelo := range reserva.Vuelos {
		var seccion string
		if i == 0 {
			seccion = "Ida"
		} else {
			seccion = "Vuelta"
		}

		fmt.Print(
			fmt.Sprintf("%s:\n %s %s %s\n", seccion, vuelo.NumeroVuelo, vuelo.HoraSalida, vuelo.HoraLlegada),
		)
	}

	// Mostrar detalle de pasajeros
	fmt.Print(
		"Pasajeros:\n",
	)

	for _, pasajero := range reserva.Pasajeros {
		// Datos personales
		fmt.Print(
			fmt.Sprintf("%s %d\n", pasajero.Nombre, pasajero.Edad),
		)

		// Ancillaries
		for _, ancillaries := range pasajero.Ancillaries {
			if ancillaries.Ida != nil {
				fmt.Print("Ancillaries ida: ")
				for _, ancillary := range ancillaries.Ida {
					fmt.Print(
						fmt.Sprintf("%s", ancillary.Ssr),
					)
				}
				fmt.Println()
			}

			if ancillaries.Vuelta != nil {
				fmt.Print("Ancillaries vuelta: ")
				for _, ancillary := range ancillaries.Vuelta {
					fmt.Print(
						fmt.Sprintf("%s", ancillary.Ssr),
					)
				}
				fmt.Println()
			}
		}
	}
}
