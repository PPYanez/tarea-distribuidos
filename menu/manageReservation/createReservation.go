package menu

import (
	"bytes"
	"distribuidos/tarea-1/db"
	"distribuidos/tarea-1/models"
	"distribuidos/tarea-1/utilities"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
)

func createReservation(fechaIda string, fechaVuelta string, origen string, destino string, cantidadPasajeros int) {
	// Seleccionar vuelos
	fmt.Println("Vuelos disponibles:")
	vuelos, precioPasajeIda, precioPasajeVuelta := chooseVuelos(fechaIda, fechaVuelta, origen, destino)

	// Ingresar pasajeros
	pasajeros := setPassengersInfo(cantidadPasajeros)

	// Actualizar balances con precio de pasajes y obtener costo total
	var costoTotal int
	for i := range pasajeros {
		pasajeros[i].Balances.VueloIda = precioPasajeIda
		pasajeros[i].Balances.VueloVuelta = precioPasajeVuelta

		costoTotal += pasajeros[i].Balances.VueloIda +
			pasajeros[i].Balances.VueloVuelta +
			pasajeros[i].Balances.AncillariesIda +
			pasajeros[i].Balances.AncillariesVuelta
	}

	// Crear reserva
	reserva := models.Reserva{
		Pasajeros: pasajeros,
		Vuelos:    vuelos,
	}

	reservaJson, err := json.Marshal(reserva)
	if err != nil {
		log.Fatal("Error al crear reserva")
	}

	url := utilities.CreateUrl("reserva", nil)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reservaJson))
	if err != nil {
		log.Fatal("Error al crear reserva")
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error al crear reserva")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error al recibir respuesta del servidor")
	}

	var reservaResponse models.PNR

	if err := json.Unmarshal(body, &reservaResponse); err != nil {
		log.Fatal("Respuesta no válida")
	}

	fmt.Println("La reserva fue generada, el PNR es:", reservaResponse.Pnr)
	fmt.Println(fmt.Sprintf("El costo total de la reserva fué de $%d", costoTotal))
}

func chooseVuelos(fechaIda string, fechaVuelta string, origen string, destino string) ([]models.Vuelo, int, int) {
	// Vuelos de ida
	fmt.Println("Ida:")
	vueloIda, precioPasajeIda, err := chooseVuelo(origen, destino, fechaIda)
	if err != nil {
		fmt.Println("No se encontraron vuelos de ida, intente con otra fecha")
		return nil, 0, 0
	}

	// Vuelos de vuelta
	fmt.Println("Vuelta:")
	vueloVuelta, precioPasajeVuelta, err := chooseVuelo(destino, origen, fechaVuelta)
	if err != nil {
		fmt.Println("No se encontraron vuelos de vuelta, se reservará solo el vuelo de ida")
	}

	var vuelosReserva []models.Vuelo
	vuelosReserva = append(vuelosReserva, *vueloIda)
	vuelosReserva = append(vuelosReserva, *vueloVuelta)

	return vuelosReserva, precioPasajeIda, precioPasajeVuelta
}

func chooseVuelo(origen string, destino string, fecha string) (*models.Vuelo, int, error) {
	// Obtener vuelos
	queries := map[string]string{
		"origen":  origen,
		"destino": destino,
		"fecha":   fecha,
	}

	url := utilities.CreateUrl("vuelo", queries)
	vuelos := requestVuelos(url)
	var vuelo models.Vuelo

	if len(vuelos) > 0 {
		// Mostrar opciones
		showVuelos(vuelos)

		// Reservar ida
		fmt.Println("Ingrese una opción: ")

		var ida int
		fmt.Scanln(&ida)
		vuelo = vuelos[ida-1]
	} else {
		return nil, 0, fmt.Errorf("No se encontraron vuelos")
	}

	precioPasaje := utilities.CalculateFlightPrice(vuelo)

	return &vuelo, precioPasaje, nil
}

func requestVuelos(url string) []models.Vuelo {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("No se encontraron vuelos")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("No se encontraron vuelos")
	}

	var vuelos []models.Vuelo

	if err := json.Unmarshal(body, &vuelos); err != nil {
		log.Fatal("No se encontraron vuelos")
	}

	return vuelos
}

func showVuelos(vuelos []models.Vuelo) {
	for _, vuelo := range vuelos {
		precioPasaje := utilities.CalculateFlightPrice(vuelo)

		fmt.Println(
			fmt.Sprintf("%s %s %s $%d", vuelo.NumeroVuelo, vuelo.HoraSalida, vuelo.HoraLlegada, precioPasaje),
		)
	}
}

func setPassengersInfo(cantidadPasajeros int) []models.Pasajero {
	pasajeros := make([]models.Pasajero, cantidadPasajeros)

	for i := 0; i < cantidadPasajeros; i++ {
		fmt.Println(fmt.Sprintf("%s %d%s", "Pasajero", i+1, ":"))

		pasajero := setPassengerInfo()
		pasajeros[i] = pasajero
	}

	return pasajeros
}

func setPassengerInfo() models.Pasajero {
	nombre, apellido, edad := getPassengerData()
	ancillariesData := db.GetAncillaries()

	fmt.Println("Ancillaries de ida:")
	showAncillaries(ancillariesData)

	ancillariesIda, totalAncillariesIda := chooseAncillaries(ancillariesData)
	fmt.Println("Total ancillaries: ", totalAncillariesIda)

	fmt.Println("Ancillaries de vuelta:")
	showAncillaries(ancillariesData)

	ancillariesVuelta, totalAncillariesVuelta := chooseAncillaries(ancillariesData)
	fmt.Println("Total ancillaries: ", totalAncillariesVuelta)

	selectedAncillaries := []models.AncillariePasajero{
		{
			Ida:    ancillariesIda,
			Vuelta: ancillariesVuelta,
		},
	}

	balances := models.Balance{
		AncillariesIda:    totalAncillariesIda,
		AncillariesVuelta: totalAncillariesVuelta,
	}

	return models.Pasajero{
		Nombre:      nombre,
		Apellido:    apellido,
		Edad:        edad,
		Ancillaries: selectedAncillaries,
		Balances:    balances,
	}
}

func getPassengerData() (string, string, int) {
	var nombre string
	var apellido string
	var edad int

	fmt.Print("Ingrese el nombre: ")
	fmt.Scanln(&nombre)
	fmt.Print("Ingrese el apellido: ")
	fmt.Scanln(&apellido)
	fmt.Print("Ingrese la edad: ")
	fmt.Scanln(&edad)

	return nombre, apellido, edad
}

func showAncillaries(ancillariesData map[string]models.AncillarieData) {
	// Display sorted by key
	keys := make([]string, 0, len(ancillariesData))
	for k := range ancillariesData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := ancillariesData[k]
		fmt.Printf("%s: %s - %d\n", k, v.Nombre, v.Precio)
	}
}

func chooseAncillaries(ancillariesData map[string]models.AncillarieData) ([]models.AncillarieDetail, int) {
	var selection string
	fmt.Print("Ingrese los ancillaries (separados por coma): ")
	fmt.Scanln(&selection)

	selectedAncillariesSplitted := strings.Split(selection, ",")

	selectedAncillaries := []models.AncillarieDetail{}

	var selectedAncillariesTotalPrice int

	for _, ancillary := range selectedAncillariesSplitted {
		ancillaryObject := ancillariesData[ancillary]

		selectedAncillariesTotalPrice += ancillaryObject.Precio

		ssr := ancillaryObject.Ssr

		found := false
		for i, selectedAncillary := range selectedAncillaries {
			if selectedAncillary.Ssr == ssr {
				selectedAncillaries[i].Cantidad++
				found = true
				break
			}
		}

		if !found {
			selectedAncillaries = append(selectedAncillaries, models.AncillarieDetail{
				Ssr:      ancillaryObject.Ssr,
				Cantidad: 1,
			})
		}
	}

	return selectedAncillaries, selectedAncillariesTotalPrice
}
