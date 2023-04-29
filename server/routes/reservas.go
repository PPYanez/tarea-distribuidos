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
	nuevaReserva := new(models.Reserva)

	defer cancel()

	if err := c.BindJSON(&nuevaReserva); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	// Guardar reserva en la base de datos
	reservaCollection := client.Database("distribuidos").Collection("Reservas")

	pnrReserva := utilities.GenPNR()

	payload := models.Reserva{
		Vuelos:    nuevaReserva.Vuelos,
		Pasajeros: nuevaReserva.Pasajeros,
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

	// Guardar ruta en la base de datos (para facilitar estadísticas)
	comprasCollection := client.Database("distribuidos").Collection("Compras")
	totalPasajesIda := 0
	totalPasajesVuelta := 0

	for _, pasajero := range nuevaReserva.Pasajeros {
		totalPasajesIda += pasajero.Balances.VueloIda
		totalPasajesVuelta += pasajero.Balances.VueloVuelta
	}

	compra := models.Compra{
		Origen:            nuevaReserva.Vuelos[0].Origen,
		Destino:           nuevaReserva.Vuelos[0].Destino,
		TotalPasajes:      totalPasajesIda,
		CantidadPasajeros: len(nuevaReserva.Pasajeros),
		FechaVuelo:        nuevaReserva.Vuelos[0].Fecha,
	}

	_, err = comprasCollection.InsertOne(ctx, compra)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	if len(nuevaReserva.Vuelos) > 1 {
		compra = models.Compra{
			Origen:            nuevaReserva.Vuelos[1].Origen,
			Destino:           nuevaReserva.Vuelos[1].Destino,
			TotalPasajes:      totalPasajesVuelta,
			CantidadPasajeros: len(nuevaReserva.Pasajeros),
			FechaVuelo:        nuevaReserva.Vuelos[1].Fecha,
		}
		_, err = comprasCollection.InsertOne(ctx, compra)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err})
			return
		}
	}

	// Disminuir stock de ancillaries y de pasajeros
	vueloCollection := client.Database("distribuidos").Collection("Vuelos")

	for _, pasajero := range nuevaReserva.Pasajeros {
		ancillariesReservados := pasajero.Ancillaries

		vueloCollection.UpdateOne(
			ctx,
			bson.M{"numero_vuelo": nuevaReserva.Vuelos[0].NumeroVuelo},
			bson.M{"$inc": bson.M{"avion.stock_de_pasajeros": -1}},
		)

		if len(nuevaReserva.Vuelos) > 1 {
			vueloCollection.UpdateOne(
				ctx,
				bson.M{"numero_vuelo": nuevaReserva.Vuelos[1].NumeroVuelo},
				bson.M{"$inc": bson.M{"avion.stock_de_pasajeros": -1}},
			)
		}

		for _, ancillaries := range ancillariesReservados {
			if ancillaries.Ida != nil {
				for _, ancillary_ida := range ancillaries.Ida {
					vueloCollection.UpdateOne(
						ctx,
						bson.M{"numero_vuelo": nuevaReserva.Vuelos[0].NumeroVuelo, "ancillaries.ssr": ancillary_ida.Ssr},
						bson.M{"$inc": bson.M{"ancillaries.$.stock": -1 * ancillary_ida.Cantidad}},
					)
				}
			}

			if ancillaries.Vuelta != nil && len(nuevaReserva.Vuelos) > 1 {
				for _, ancillary_vuelta := range ancillaries.Vuelta {
					vueloCollection.UpdateOne(
						ctx,
						bson.M{"numero_vuelo": nuevaReserva.Vuelos[1].NumeroVuelo, "ancillaries.ssr": ancillary_vuelta.Ssr},
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

func UpdateReserva(c *gin.Context) {
	client := db.GetClient()

	PNR := c.Query("pnr")
	apellido := c.Query("apellido")

	if PNR == "" || apellido == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No se ha ingresado un PNR o apellido"})
		return
	}

	collection := client.Database("distribuidos").Collection("Reservas")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Bindear JSON a estructura Reserva
	reemplazo := new(models.Reserva)

	defer cancel()

	if err := c.BindJSON(&reemplazo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	reemplazo.PNR = PNR

	// Actualizar reserva
	filter := bson.M{"pnr": PNR, "pasajeros.apellido": apellido}
	update := bson.M{"$set": reemplazo}

	resp, err := collection.UpdateOne(ctx, filter, update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "No se pudo actualizar la reserva"})
		return
	}

	if resp.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "La reserva requerida no existe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"PNR": PNR})
}
