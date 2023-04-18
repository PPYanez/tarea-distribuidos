package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

	defer cancel()

	opts := options.Find().SetProjection(bson.D{primitive.E{Key: "ancillaries", Value: 0}})

	cursor, err := collection.Find(
		ctx,
		bson.M{"origen": origen, "destino": destino, "fecha": fecha},
		opts,
	)

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			make([]string, 0),
		)
		return
	}

	var results []models.Vuelo

	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, results)
}

func UpdateStock(c *gin.Context) {
	client := db.GetClient()

	// Obtener query strings
	numero_vuelo := c.Query("numero_vuelo")
	origen := c.Query("origen")
	destino := c.Query("destino")
	fecha := c.Query("fecha")

	if numero_vuelo == "" || origen == "" || destino == "" || fecha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No se han entregado todos los datos necesarios"})
		return
	}

	// Bindear JSON a estructura newStock
	newStock := new(models.Stock)

	if err := c.BindJSON(&newStock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "JSON incorrecto"})
		log.Fatal(err)
		return
	}

	fmt.Printf("---%+v---\n", newStock)

	// Actualizar stock en la base de datos
	filter := bson.M{"numero_vuelo": numero_vuelo, "origen": origen, "destino": destino, "fecha": fecha}
	update := bson.M{"$set": bson.M{"avion.stock_de_pasajeros": newStock.Stock}}
	after := options.After
	opts := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	collection := client.Database("distribuidos").Collection("Vuelos")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	resp := collection.FindOneAndUpdate(ctx, filter, update, &opts)
	if resp.Err() == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontró el vuelo a actualizar"})
		return
	}

	result := new(models.Vuelo)
	decodeErr := resp.Decode(&result)
	if decodeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error al obtener el vuelo actualizado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"numero_vuelo": result.NumeroVuelo,
		"origen":       result.Origen,
		"destino":      result.Destino,
		"hora_salida":  result.HoraSalida,
		"hora_llegada": result.HoraLlegada,
	})
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
