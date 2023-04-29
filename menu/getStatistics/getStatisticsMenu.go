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

func GetStatisticsMenu() {
	statistics := getStatistcs()

	fmt.Println("Ruta con mayor ganancia:", statistics.RutaMayorGanancia)
	fmt.Println("Ruta con menor ganancia:", statistics.RutaMenorGanancia)
	fmt.Println("Ranking de ancillaries:")

	for _, ancillary := range statistics.RankingAncillaries {
		fmt.Println(" -", ancillary.Ssr, ancillary.Nombre, ancillary.Ganancia)
	}

	showPromedioPasajeros(statistics.PromedioPasajeros)
}

func getStatistcs() models.Statistics {
	url := utilities.CreateUrl("estadisticas", nil)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Tuvimos problemas obteniendo las estadísticas, intente más tarde.")
	}

	body, _ := ioutil.ReadAll(resp.Body)

	var estadisticas models.Statistics

	if err := json.Unmarshal(body, &estadisticas); err != nil {
		log.Fatal("Tuvimos problemas obteniendo las estadísticas, intente más tarde.")
	}

	return estadisticas
}

func showPromedioPasajeros(promedioPasajeros models.PromedioPasajeros) {
	fmt.Println("Promedio de pasajeros por mes:")
	fmt.Println(" - Enero:", promedioPasajeros.Enero)
	fmt.Println(" - Febrero:", promedioPasajeros.Febrero)
	fmt.Println(" - Marzo:", promedioPasajeros.Marzo)
	fmt.Println(" - Abril:", promedioPasajeros.Abril)
	fmt.Println(" - Mayo:", promedioPasajeros.Mayo)
	fmt.Println(" - Junio:", promedioPasajeros.Junio)
	fmt.Println(" - Julio:", promedioPasajeros.Julio)
	fmt.Println(" - Agosto:", promedioPasajeros.Agosto)
	fmt.Println(" - Septiembre:", promedioPasajeros.Septiembre)
	fmt.Println(" - Octubre:", promedioPasajeros.Octubre)
	fmt.Println(" - Noviembre:", promedioPasajeros.Noviembre)
	fmt.Println(" - Diciembre:", promedioPasajeros.Diciembre)
}
