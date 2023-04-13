package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type reserva struct {
	Vuelos    []vuelo    `json:"vuelos,omitempty" bson:"vuelos,omitempty"`
	Pasajeros []pasajero `json:"pasajeros,omitempty" bson:"pasajeros,omitempty"`
}

type pasajero struct {
	Nombre      string                `json:"nombre,omitempty" bson:"nombre,omitempty"`
	Apellido    string                `json:"apellido,omitempty" bson:"apellido,omitempty"`
	Edad        int                   `json:"edad,omitempty" bson:"edad,omitempty"`
	Ancillaries []ancillarie_pasajero `json:"ancillaries,omitempty" bson:"ancillaries,omitempty"`
	Balances    balance               `json:"balances,omitempty" bson:"balances,omitempty"`
}

type ancillarie_pasajero struct {
	Ida []struct {
		Ssr      string
		Cantidad int
	} `json:"ida,omitempty" bson:"ida,omitempty"`
	Vuelta []struct {
		Ssr      string
		Cantidad int
	} `json:"vuelta,omitempty" bson:"vuelta,omitempty"`
}

type balance struct {
	AncillariesIda    int `json:"ancillariesida,omitempty" bson:"ancillariesida,omitempty"`
	VueloIda          int `json:"vueloida,omitempty" bson:"vueloida,omitempty"`
	AncillariesVuelta int `json:"ancillariesvuelta,omitempty" bson:"ancillariesvuelta,omitempty"`
	VueloVuelta       int `json:"vuelovuelta,omitempty" bson:"vuelovuelta,omitempty"`
}

type vuelo struct {
	NumeroVuelo string       `json:"numero_vuelo,omitempty" bson:"numero_vuelo,omitempty"`
	Origen      string       `json:"origen,omitempty" bson:"origen,omitempty"`
	Destino     string       `json:"destino,omitempty" bson:"destino,omitempty"`
	HoraSalida  string       `json:"hora_salida,omitempty" bson:"hora_salida,omitempty"`
	HoraLlegada string       `json:"hora_llegada,omitempty" bson:"hora_llegada,omitempty"`
	Fecha       string       `json:"fecha,omitempty" bson:"fecha,omitempty"`
	Avion       avion        `json:"avion,omitempty" bson:"avion,omitempty"`
	Ancillaries []ancillarie `json:"ancillaries,omitempty" bson:"ancillaries,omitempty"`
}

type ancillarie struct {
	Nombre string `json:"nombre,omitempty" bson:"nombre,omitempty"`
	Stock  int    `json:"stock,omitempty" bson:"stock,omitempty"`
	Ssr    string `json:"ssr,omitempty" bson:"ssr,omitempty"`
}

type avion struct {
	Modelo           string `json:"modelo,omitempty" bson:"modelo,omitempty"`
	NumeroDeSerie    string `json:"numero_de_serie,omitempty" bson:"numero_de_serie,omitempty"`
	StockDePasajeros int    `json:"stock_de_pasajeros,omitempty" bson:"stock_de_pasajeros,omitempty"`
}

func main() {
	/* Connect to mongodb */
	Mongo_URL := "mongodb://localhost:27017"
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(Mongo_URL).SetServerAPIOptions(serverAPI)

	var err error
	client, err = mongo.Connect(context.TODO(), opts)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	var result bson.M
	if err := client.Database("distribuidos").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to MongoDB!")

	/* Set up gin server */
	router := gin.Default()

	router.POST("/api/reserva", postReserva)

	router.Run("localhost:6666")
}

func postReserva(c *gin.Context) {
	var collection = client.Database("distribuidos").Collection("Reservas")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	nueva_reserva := new(reserva)

	defer cancel()

	if err := c.BindJSON(&nueva_reserva); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	payload := reserva{
		Vuelos:    nueva_reserva.Vuelos,
		Pasajeros: nueva_reserva.Pasajeros,
	}

	result, err := collection.InsertOne(ctx, payload)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Posted successfully", "Data": map[string]interface{}{"data": result}})
}

func getVuelos(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) *mongo.Cursor {
	fmt.Println("Entered getvuelos")
	collection := client.Database("distribuidos").Collection("vuelos")
	fmt.Println("got colletion " + collection.Name() + " " + collection.Database().Name())

	defer cancel()
	fmt.Println("after cancel")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	fmt.Println("got cursor")

	if err != nil {
		panic(err)
	}

	fmt.Println("return")
	return cursor
}
