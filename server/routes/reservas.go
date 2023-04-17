package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"

	"distribuidos/tarea-1/db"
	"distribuidos/tarea-1/models"
	"distribuidos/tarea-1/utilities"
)

// Endpoints para reservas

func PostReserva(c *gin.Context) {
	client := db.GetClient()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Bindear JSON a estructura Reserva
	nueva_reserva := new(models.Reserva)

	defer cancel()

	if err := c.BindJSON(&nueva_reserva); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	// Guardar reserva en la base de datos
	reservaCollection := client.Database("distribuidos").Collection("Reservas")

	pnrReserva := utilities.GenPNR()

	payload := models.Reserva{
		Vuelos:    nueva_reserva.Vuelos,
		Pasajeros: nueva_reserva.Pasajeros,
		PNR:       pnrReserva.Pnr,
	}

	_, err := reservaCollection.InsertOne(ctx, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	// Guardar PNR en la base de datos (para mantener unicidad)
	pnrCollection := client.Database("distribuidos").Collection("PNR")
	_, err = pnrCollection.InsertOne(ctx, pnrReserva)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	// Disminuir stock de ancillaries y de pasajeros
	vueloCollection := client.Database("distribuidos").Collection("Vuelos")

	for _, pasajero := range nueva_reserva.Pasajeros {
		ancillaries_reservados := pasajero.Ancillaries

		vueloCollection.UpdateOne(
			ctx,
			bson.M{"numero_vuelo": nueva_reserva.Vuelos[0].NumeroVuelo},
			bson.M{"$inc": bson.M{"avion.stock_de_pasajeros": -1}},
		)

		vueloCollection.UpdateOne(
			ctx,
			bson.M{"numero_vuelo": nueva_reserva.Vuelos[1].NumeroVuelo},
			bson.M{"$inc": bson.M{"avion.stock_de_pasajeros": -1}},
		)

		for _, ancillaries := range ancillaries_reservados {
			if ancillaries.Ida != nil {
				for _, ancillary_ida := range ancillaries.Ida {
					vueloCollection.UpdateOne(
						ctx,
						bson.M{"numero_vuelo": nueva_reserva.Vuelos[0].NumeroVuelo, "ancillaries.ssr": ancillary_ida.Ssr},
						bson.M{"$inc": bson.M{"ancillaries.$.stock": -1 * ancillary_ida.Cantidad}},
					)
				}
			}

			if ancillaries.Vuelta != nil {
				for _, ancillary_vuelta := range ancillaries.Vuelta {
					vueloCollection.UpdateOne(
						ctx,
						bson.M{"numero_vuelo": nueva_reserva.Vuelos[1].NumeroVuelo, "ancillaries.ssr": ancillary_vuelta.Ssr},
						bson.M{"$inc": bson.M{"ancillaries.$.stock": -1 * ancillary_vuelta.Cantidad}},
					)
				}
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{"PNR": pnrReserva.Pnr})
}

func GetReserva(c *gin.Context) {
	client := db.GetClient()

	PNR := c.Query("pnr")
	apellido := c.Query("apellido")

	if PNR == "" || apellido == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No se ha ingresado un PNR o apellido"})
		return
	}

	collection := client.Database("distribuidos").Collection("Reservas")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var result models.Reserva

	defer cancel()

	err := collection.FindOne(ctx, bson.M{"pnr": PNR, "pasajeros.apellido": apellido}).Decode(&result)
	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"message": "No existe una reserva con los datos ingresados, por favor verifique que están bien escritos.",
			},
		)
		return
	}

	res := map[string]interface{}{"vuelos": result.Vuelos, "pasajeros": result.Pasajeros}

	c.JSON(http.StatusOK, gin.H{"vuelos": res["vuelos"], "pasajeros": res["pasajeros"]})
}
