package routes

import (
	"distribuidos/tarea-1/db"
	"distribuidos/tarea-1/models"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetEstadisticas(c *gin.Context) {
	client := db.GetClient()
	comprasCollection := client.Database("distribuidos").Collection("Compras")
	reservasCollection := client.Database("distribuidos").Collection("Reservas")

	// Obtener estadísticas de ganancias por ruta
	lowest, highest := RouteEarningStatistics(c, comprasCollection)

	// Obtener ranking de ansillaries
	totalEarnedPerAncillarie := AncillariesRankingStatistic(c, reservasCollection)

	// Obtener promedio de pasajeros por vuelo
	promedioPasajeros := AveragePassengerStatistic(c, comprasCollection)

	statistics := models.Statistics{
		RutaMayorGanancia:  highest,
		RutaMenorGanancia:  lowest,
		RankingAncillaries: totalEarnedPerAncillarie,
		PromedioPasajeros:  promedioPasajeros,
	}

	c.JSON(200, statistics)
}

func RouteEarningStatistics(ctx *gin.Context, comprasCollection *mongo.Collection) (lowest, highest string) {
	groupStage := bson.D{
		{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "ida", Value: "$ida"},
				{Key: "vuelta", Value: "$vuelta"},
			}},
			{Key: "total_pasajes", Value: bson.D{
				{Key: "$sum", Value: "$total_pasajes"},
			}},
		}},
	}

	sortStage := bson.D{
		{Key: "$sort", Value: bson.D{
			{Key: "total_pasajes", Value: -1},
		}},
	}

	secondGroupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: nil},

		{Key: "ida_min", Value: bson.D{{Key: "$first", Value: "$_id.ida"}}},
		{Key: "vuelta_min", Value: bson.D{{Key: "$first", Value: "$_id.vuelta"}}},

		{Key: "ida_max", Value: bson.D{{Key: "$last", Value: "$_id.ida"}}},
		{Key: "vuelta_max", Value: bson.D{{Key: "$last", Value: "$_id.vuelta"}}},
	}}}

	pipeline := mongo.Pipeline{groupStage, sortStage, secondGroupStage}

	cursor, err := comprasCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	var rutaMax string
	var rutaMin string

	if cursor.Next(ctx) {
		var result bson.M
		err := cursor.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		idaMin := result["ida_min"].(string)
		vueltaMin := result["vuelta_min"].(string)
		rutaMin = fmt.Sprintf("%s - %s", idaMin, vueltaMin)

		idaMax := result["ida_max"].(string)
		vueltaMax := result["vuelta_max"].(string)
		rutaMax = fmt.Sprintf("%s - %s", idaMax, vueltaMax)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	return rutaMin, rutaMax
}

func AncillariesRankingStatistic(c *gin.Context, reservasCollection *mongo.Collection) []models.StatisticAncillarie {
	cursor, err := reservasCollection.Find(c, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	totalEarnedPerAncillarie := make([]models.StatisticAncillarie, len(db.GetAncillaries()))
	ancillaries := db.GetAncillaries()

	for i, ancillarie := range ancillaries {
		intIndex, _ := strconv.Atoi(i)

		totalEarnedPerAncillarie[intIndex-1].Nombre = ancillarie.Nombre
		totalEarnedPerAncillarie[intIndex-1].Ssr = ancillarie.Ssr
		totalEarnedPerAncillarie[intIndex-1].Ganancia = 0
	}

	for cursor.Next(c) {
		var result models.Reserva
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}

		for _, pasajero := range result.Pasajeros {
			// Add ancillaries for every passenger in the flight
			ancillariesIda := pasajero.Ancillaries[0].Ida

			for _, ancillarie := range ancillariesIda {
				ancillarieIndex := indexAncillarie(ancillarie, ancillaries)
				totalEarnedPerAncillarie[ancillarieIndex-1].Ganancia += ancillarie.Cantidad * getAncillariePrice(ancillarie)
			}

			// If the reservation has departure and return, add the return ancillaries
			if len(result.Vuelos) > 1 {
				ancillariesVuelta := pasajero.Ancillaries[1].Vuelta

				for _, ancillarie := range ancillariesVuelta {
					ancillarieIndex := indexAncillarie(ancillarie, ancillaries)
					totalEarnedPerAncillarie[ancillarieIndex-1].Ganancia += ancillarie.Cantidad * getAncillariePrice(ancillarie)
				}
			}
		}
	}

	return totalEarnedPerAncillarie
}

func indexAncillarie(ancillarieObject models.AncillarieDetail, ancillaries map[string]models.AncillarieData) int {
	for i, ancillarie := range ancillaries {
		if ancillarie.Ssr == ancillarieObject.Ssr {
			idx, _ := strconv.Atoi(i)
			return idx
		}
	}

	return -1
}

func getAncillariePrice(ancillarieObject models.AncillarieDetail) int {
	ancillaries := db.GetAncillaries()
	for _, ancillarie := range ancillaries {
		if ancillarie.Ssr == ancillarieObject.Ssr {
			return ancillarie.Precio
		}
	}

	return -1
}

func AveragePassengerStatistic(c *gin.Context, comprasCollection *mongo.Collection) models.PromedioPasajeros {
	cursor, err := comprasCollection.Find(c, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	var results []bson.M
	if err = cursor.All(c, &results); err != nil {
		log.Fatal(err)
	}

	pasajerosPerMonth := make([]int, 12)

	// Iterate over the results
	for _, result := range results {
		monthNumber, _ := strconv.Atoi(result["fecha_vuelo"].(string)[3:5])
		pasajerosPerMonth[monthNumber-1] += int(result["cantidad_pasajeros"].(int32))
	}

	PromedioPasajeros := averagePasajeros(pasajerosPerMonth)

	return PromedioPasajeros
}

func averagePasajeros(pasajerosPerMonth []int) models.PromedioPasajeros {
	var PromedioPasajeros models.PromedioPasajeros

	for i := 0; i < 12; i++ {
		pasajerosThisMonth := pasajerosPerMonth[i]

		setMonth(&PromedioPasajeros, pasajerosThisMonth, i)
	}

	return PromedioPasajeros
}

func setMonth(PromedioPasajeros *models.PromedioPasajeros, pasajerosThisMonth, monthNumber int) {
	switch monthNumber {
	case 0:
		PromedioPasajeros.Enero += pasajerosThisMonth / 31
	case 1:
		PromedioPasajeros.Febrero += pasajerosThisMonth / 28
	case 2:
		PromedioPasajeros.Marzo += pasajerosThisMonth / 31
	case 3:
		PromedioPasajeros.Abril += pasajerosThisMonth / 30
	case 4:
		PromedioPasajeros.Mayo += pasajerosThisMonth / 31
	case 5:
		PromedioPasajeros.Junio += pasajerosThisMonth / 30
	case 6:
		PromedioPasajeros.Julio += pasajerosThisMonth / 31
	case 7:
		PromedioPasajeros.Agosto += pasajerosThisMonth / 31
	case 8:
		PromedioPasajeros.Septiembre += pasajerosThisMonth / 30
	case 9:
		PromedioPasajeros.Octubre += pasajerosThisMonth / 31
	case 10:
		PromedioPasajeros.Noviembre += pasajerosThisMonth / 30
	case 11:
		PromedioPasajeros.Diciembre += pasajerosThisMonth / 31
	}
}
