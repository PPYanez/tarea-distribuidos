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

func updateReservationMenu(pnr string, apellido string) {
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

		// Obtener reserva con los datos ingresados
		originalReservation := discreteGetReservation(pnr, apellido)

		if len(originalReservation.Pasajeros) == 0 && len(originalReservation.Vuelos) == 0 {
			log.Fatal("Reserva no encontrada")
			break
		}

		if option == 1 {
			// Mostrar vuelos que se pueden cambiar
			fmt.Println("Vuelos:")
			fmt.Printf("1. ida: %s %s - %s\n", originalReservation.Vuelos[0].NumeroVuelo, originalReservation.Vuelos[0].HoraSalida, originalReservation.Vuelos[0].HoraLlegada)
			fmt.Printf("2. vuelta: %s %s - %s\n", originalReservation.Vuelos[1].NumeroVuelo, originalReservation.Vuelos[1].HoraSalida, originalReservation.Vuelos[1].HoraLlegada)

			// Elegir vuelo a cambiar
			var vueloToReplace int
			fmt.Print("Ingrese una opción: ")
			fmt.Scanln(&vueloToReplace)
			vueloToReplace -= 1

			updatedReservation, err := changeFlightDate(originalReservation, vueloToReplace, pnr, apellido)
			if err != nil {
				log.Fatal(err)
				continue
			}

			updateReservation(pnr, apellido, updatedReservation)
			updateStocks(updatedReservation, updatedReservation.Vuelos[vueloToReplace], originalReservation.Vuelos[vueloToReplace])
		}
	}
}

func changeFlightDate(reserva models.Reserva, vueloToReplace int, pnr string, apellido string) (models.Reserva, error) {
	vuelo := reserva.Vuelos[vueloToReplace]

	// Elegir nueva fecha
	var newDate string
	fmt.Print("Ingrese nueva fecha: ")
	fmt.Scanln(&newDate)

	vuelos := getVuelos(vuelo, newDate)
	if len(vuelos) == 0 {
		return models.Reserva{}, fmt.Errorf("No hay vuelos disponibles para la nueva fecha")
	}

	// Mostrar vuelos disponibles para la nueva fecha
	fmt.Println("Vuelos disponibles:")
	for i, v := range vuelos {
		fmt.Printf("%d. %s %s - %s\n", i+1, v.NumeroVuelo, v.HoraSalida, v.HoraLlegada)
	}
	var selectedNewVuelo string
	fmt.Print("Ingrese una opción: ")
	fmt.Scanln(&selectedNewVuelo)

	// Actualizar objeto con reserva
	var newReserva models.Reserva
	newReserva.PNR = reserva.PNR
	newReserva.Pasajeros = append(newReserva.Pasajeros, reserva.Pasajeros...)
	newReserva.Vuelos = append(newReserva.Vuelos, reserva.Vuelos...)

	newReserva.Vuelos[vueloToReplace] = vuelos[vueloToReplace]

	return newReserva, nil
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

func updateReservation(pnr string, apellido string, updatedReservation models.Reserva) {
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

func updateStocks(reserva models.Reserva, newVuelo models.Vuelo, oldVuelo models.Vuelo) {
	updateStock(reserva, newVuelo, "decrease")
	updateStock(reserva, oldVuelo, "add")
}

func updateStock(reserva models.Reserva, vuelo models.Vuelo, operator string) {
	queries := map[string]string{
		"numero_vuelo": vuelo.NumeroVuelo,
		"origen":       vuelo.Origen,
		"destino":      vuelo.Destino,
		"fecha":        vuelo.Fecha,
	}
	url := utilities.CreateUrl("vuelo", queries)

	var stock models.Stock

	if operator == "add" {
		stock = models.Stock{
			Stock: vuelo.Avion.StockDePasajeros + len(reserva.Pasajeros),
		}
	} else {
		stock = models.Stock{
			Stock: vuelo.Avion.StockDePasajeros - len(reserva.Pasajeros),
		}
	}

	stockJson, err := json.Marshal(&stock)
	if err != nil {
		log.Fatal("Error al actualizar el stock")
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(stockJson))
	if err != nil {
		log.Fatal("Error al actualizar el stock")
	}

	req.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)

	defer req.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Error del servidor al actualizar el stock")
	}
}
