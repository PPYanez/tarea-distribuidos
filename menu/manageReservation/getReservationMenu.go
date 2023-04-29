package menu

import (
	"distribuidos/tarea-1/models"
	"distribuidos/tarea-1/utilities"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Obtener reserva mostrando los datos necesarios para el menú
func getReservationMenu(pnr string, apellido string) {
	queries := map[string]string{
		"pnr":      pnr,
		"apellido": apellido,
	}

	url := utilities.CreateUrl("reserva", queries)

	resp, err := http.Get(url)
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

		fmt.Printf("%s:\n %s %s %s\n", seccion, vuelo.NumeroVuelo, vuelo.HoraSalida, vuelo.HoraLlegada)
	}

	// Mostrar detalle de pasajeros
	fmt.Print(
		"Pasajeros:\n",
	)

	for _, pasajero := range reserva.Pasajeros {
		// Datos personales
		fmt.Printf("%s %d\n", pasajero.Nombre, pasajero.Edad)

		// Ancillaries
		for _, ancillaries := range pasajero.Ancillaries {
			if ancillaries.Ida != nil {
				fmt.Print("Ancillaries ida: ")
				for _, ancillary := range ancillaries.Ida {
					fmt.Printf("%s ", ancillary.Ssr)
				}
				fmt.Println()
			}

			if ancillaries.Vuelta != nil {
				fmt.Print("Ancillaries vuelta: ")
				for _, ancillary := range ancillaries.Vuelta {
					fmt.Printf("%s ", ancillary.Ssr)
				}
				fmt.Println()
			}
		}
	}
}

// Obtener reserva, pero sin prints
func discreteGetReservation(pnr string, apellido string) models.Reserva {
	queries := map[string]string{
		"pnr":      pnr,
		"apellido": apellido,
	}

	url := utilities.CreateUrl("reserva", queries)

	resp, err := http.Get(url)
	if err != nil {
		return models.Reserva{}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Reserva{}
	}

	var reserva models.Reserva
	if err := json.Unmarshal(body, &reserva); err != nil {
		return models.Reserva{}
	}

	return reserva
}
