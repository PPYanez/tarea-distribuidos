package menu

import (
	"distribuidos/tarea-1/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

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
