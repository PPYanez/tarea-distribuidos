package menu

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func createReservation(fechaIda string, fechaVuelta string, origen string, destino string, cantidadPasajeros int) {
	// Obtener vuelos
	url, err := url.Parse("http://localhost:5000/api/vuelo")
	if err != nil {
		log.Fatal("URL no válida")
	}

	values := url.Query()
	values.Add("origen", origen)
	values.Add("destino", destino)
	values.Add("fecha", fechaIda)

	url.RawQuery = values.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		log.Fatal("No se encontraron vuelos")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("No se encontraron vuelos")
	}

	fmt.Print(string(body))

	// var vuelos []models.Vuelo

	// if err := json.Unmarshal(body, &vuelos); err != nil {
	// 	log.Fatal("Respuesta no válida")
	// }
}
