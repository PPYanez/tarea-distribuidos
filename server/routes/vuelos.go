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
)

// Endpoints para vuelos
func GetVuelo(c *gin.Context) {
	client := db.GetClient()

	origen := c.Query("origen")
	destino := c.Query("destino")
	fecha := c.Query("fecha")

	if origen == "" || destino == "" || fecha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Faltan parámetros en la consulta"})
		return
	}

	collection := client.Database("distribuidos").Collection("Vuelos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var result models.Vuelo

	defer cancel()

	err := collection.FindOne(ctx, bson.M{"origen": origen, "destino": destino, "fecha": fecha}).Decode(&result)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			make([]string, 0),
		)
		return
	}

	res := []map[string]interface{}{
		{
			"numero_vuelo": result.NumeroVuelo,
			"origen":       result.Origen,
			"destino":      result.Destino,
			"hora_salida":  result.HoraSalida,
			"hora_llegada": result.HoraLlegada,
			"fecha":        result.Fecha,
			"avion":        result.Avion,
		},
	}

	c.JSON(http.StatusOK, res)
}

func PostVuelo(c *gin.Context) {
	client := db.GetClient()

	vueloCollection := client.Database("distribuidos").Collection("Vuelos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Bindear JSON a estructura Reserva
	nuevo_vuelo := new(models.Vuelo)

	defer cancel()

	if err := c.BindJSON(&nuevo_vuelo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	// Guardar reserva en la base de datos

	payload := models.Vuelo{
		NumeroVuelo: nuevo_vuelo.NumeroVuelo,
		Origen:      nuevo_vuelo.Origen,
		Destino:     nuevo_vuelo.Destino,
		HoraSalida:  nuevo_vuelo.HoraSalida,
		HoraLlegada: nuevo_vuelo.HoraLlegada,
		Fecha:       nuevo_vuelo.Fecha,
		Avion:       nuevo_vuelo.Avion,
		Ancillaries: nuevo_vuelo.Ancillaries,
	}

	_, err := vueloCollection.InsertOne(ctx, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"vuelo": payload})

}
