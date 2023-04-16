package main

import (
	"distribuidos/tarea-1/models"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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
	url, err := url.Parse("http://localhost:5000/api/reserva")
	if err != nil {
		log.Fatal("URL no válida")
	}

	values := url.Query()
	values.Add("pnr", pnr)
	values.Add("apellido", apellido)

	url.RawQuery = values.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		log.Fatal("Reserva no encontrada")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Reserva no encontrada")
	}

	var reserva models.Reserva

	if err := json.Unmarshal(body, &reserva); err != nil {
		log.Fatal("Respuesta no válida")
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
						fmt.Sprintf("%s ", ancillary.Ssr),
					)
				}
				fmt.Println()
			}

			if ancillaries.Vuelta != nil {
				fmt.Print("Ancillaries vuelta: ")
				for _, ancillary := range ancillaries.Vuelta {
					fmt.Print(
						fmt.Sprintf("%s ", ancillary.Ssr),
					)
				}
				fmt.Println()
			}
		}
	}
}
