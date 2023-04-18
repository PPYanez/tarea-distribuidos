package menu

import (
	"bytes"
	"distribuidos/tarea-1/models"
	"distribuidos/tarea-1/utilities"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func updateReservation(pnr string, apellido string) {
	keepRunning := true

	for keepRunning {
		fmt.Println(
			"Opciones:\n",
			"1. Cambiar fecha de vuelo\n",
			"2. Adicionar ancillaries\n",
			"3. Salir",
		)

		var option int
		fmt.Print("Ingrese una opción: ")
		fmt.Scanln(&option)

		if option == 3 {
			keepRunning = false
			break
		}

		var updatedReservation *models.Reserva

		if option == 1 {
			var err error
			updatedReservation, err = changeFlightDate(pnr, apellido)

			if err != nil {
				log.Fatal("No existen vuelos para la fecha ingresada")
				continue
			}
		}

		performUpdate(pnr, apellido, updatedReservation)
	}
}

func changeFlightDate(pnr string, apellido string) (*models.Reserva, error) {
	reserva := discreteGetReservation(pnr, apellido)

	if reserva == nil {
		return nil, fmt.Errorf("No existe una reserva con los datos ingresados")
	}

	fmt.Println("Vuelos:")
	fmt.Printf("1. ida: %s %s - %s\n", reserva.Vuelos[0].NumeroVuelo, reserva.Vuelos[0].HoraSalida, reserva.Vuelos[0].HoraLlegada)
	fmt.Printf("2. vuelta: %s %s - %s\n", reserva.Vuelos[1].NumeroVuelo, reserva.Vuelos[1].HoraSalida, reserva.Vuelos[1].HoraLlegada)

	var vueloToReplace int
	fmt.Print("Ingrese una opción: ")
	fmt.Scanln(&vueloToReplace)
	vueloToReplace -= 1

	vuelo := reserva.Vuelos[vueloToReplace]

	var newDate string
	fmt.Print("Ingrese nueva fecha: ")
	fmt.Scanln(&newDate)

	vuelos := getVuelos(vuelo, newDate)

	fmt.Println("Vuelos disponibles:")
	for i, vuelo := range vuelos {
		fmt.Printf("%d. %s %s - %s\n", i+1, vuelo.NumeroVuelo, vuelo.HoraSalida, vuelo.HoraLlegada)
	}
	var selectedNewVuelo string
	fmt.Print("Ingrese una opción: ")
	fmt.Scanln(&selectedNewVuelo)

	// Actualizar reserva
	reserva.Vuelos[vueloToReplace] = vuelos[vueloToReplace]

	return reserva, nil
}

func getVuelos(vuelo models.Vuelo, newDate string) []models.Vuelo {
	queries := map[string]string{
		"origen":  vuelo.Origen,
		"destino": vuelo.Destino,
		"fecha":   newDate,
	}
	url := utilities.CreateUrl("vuelo", queries)

	resp, err := http.Get(url)
	if err != nil {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var vuelos []models.Vuelo
	if err := json.Unmarshal(body, &vuelos); err != nil {
		return nil
	}

	return vuelos
}

func performUpdate(pnr string, apellido string, updatedReservation *models.Reserva) {
	// Actualizar la reserva en la base de datos
	queries := map[string]string{
		"pnr":      pnr,
		"apellido": apellido,
	}
	url := utilities.CreateUrl("reserva", queries)

	updatedReservationJson, err := json.Marshal(&updatedReservation)
	if err != nil {
		log.Fatal("Error al actualizar la reserva")
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(updatedReservationJson))
	if err != nil {
		log.Fatal("Error al actualizar la reserva")
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)

	defer req.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("La reserva fué modificada exitosamente!")
	} else {
		log.Fatal("Error del servidor al actualizar la reserva")
	}
}
