package utilities

import (
	"distribuidos/tarea-1/models"
	"fmt"
	"time"
)

func CalculateFlightPrice(vuelo models.Vuelo) int {
	// DD/MM/YY -> YYYY-MM-DD
	date := DateFormat(vuelo.Fecha)

	salida, _ := time.Parse(
		"2006-01-02 15:04",
		fmt.Sprintf("%s %s", date, vuelo.HoraSalida),
	)
	llegada, _ := time.Parse(
		"2006-01-02 15:04",
		fmt.Sprintf("%s %s", date, vuelo.HoraLlegada),
	)

	diff := llegada.Sub(salida)
	minutes := int(diff.Minutes())

	precioPasaje := minutes * 590

	return precioPasaje
}
