package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Estructuras
type reserva struct {
	Vuelos    []vuelo    `json:"vuelos,omitempty" bson:"vuelos,omitempty"`
	Pasajeros []pasajero `json:"pasajeros,omitempty" bson:"pasajeros,omitempty"`
	PNR       string     `json:"pnr,omitempty" bson:"pnr,omitempty"`
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

type PNR struct {
	Value string `json:"value,omitempty" bson:"value,omitempty"`
}

var client *mongo.Client

func main() {
	// Conectarse a MongoDB
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

	// Levantar servidor Gin
	router := gin.Default()

	router.POST("/api/reserva", postReserva)
	router.GET("/api/reserva", getReserva)

	router.Run("localhost:6666")
}

func postReserva(c *gin.Context) {
	var reservaCollection = client.Database("distribuidos").Collection("Reservas")
	var pnrCollection = client.Database("distribuidos").Collection("PNR")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Bindear JSON a estructura Reserva
	nueva_reserva := new(reserva)

	defer cancel()

	if err := c.BindJSON(&nueva_reserva); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		log.Fatal(err)
		return
	}

	// Guardar reserva en la base de datos
	pnrReserva := GenPNR()

	payload := reserva{
		Vuelos:    nueva_reserva.Vuelos,
		Pasajeros: nueva_reserva.Pasajeros,
		PNR:       pnrReserva.Value,
	}

	_, err := pnrCollection.InsertOne(ctx, pnrReserva)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	_, err = reservaCollection.InsertOne(ctx, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"PNR": pnrReserva.Value})

}

func getReserva(c *gin.Context) {
	PNR := c.Query("pnr")
	apellido := c.Query("apellido")

	if PNR == "" || apellido == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No se ha ingresado un PNR o apellido"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var collection = client.Database("distribuidos").Collection("Reservas")

	var result reserva

	defer cancel()

	err := collection.FindOne(ctx, bson.M{"pnr": PNR, "pasajeros.apellido": apellido}).Decode(&result)
	res := map[string]interface{}{"vuelos": result.Vuelos, "pasajeros": result.Pasajeros}

	if err != nil {
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"message": "No existe una reserva con los datos ingresados, por favor verifique que están bien escritos.",
			},
		)
		return
	}

	c.JSON(http.StatusOK, gin.H{"vuelos": res["vuelos"], "pasajeros": res["pasajeros"]})
}

func GenPNR() PNR {
	var letters = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	var numbers = []rune("0123456789")
	var allCharacters = append(letters, numbers...)
	existingPNR := getExistingPNR()

	var generatedPNR []rune

	for true {
		pnrIsNew := true

		// generate pnr
		generatedPNR = make([]rune, 5)

		for i := range generatedPNR {
			generatedPNR[i] = allCharacters[rand.Intn(len(allCharacters))]
		}

		if !strings.ContainsAny(string(generatedPNR), "0123456789") {
			generatedPNR = append(generatedPNR, numbers[rand.Intn(len(numbers))])
		} else if !strings.ContainsAny(string(generatedPNR), "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			generatedPNR = append(generatedPNR, letters[rand.Intn(len(letters))])
		} else {
			generatedPNR = append(generatedPNR, allCharacters[rand.Intn(len(allCharacters))])
		}

		// check if pnr is new
		for _, pnrInDB := range existingPNR {
			if pnrInDB.Value == string(generatedPNR) {
				pnrIsNew = false
				break
			}
		}

		if pnrIsNew {
			break
		}
	}

	return PNR{Value: string(generatedPNR)}
}

func getExistingPNR() []PNR {
	var pnrCollection = client.Database("distribuidos").Collection("PNR")
	findOptions := options.Find()
	var results []PNR

	cur, err := pnrCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem PNR
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	return results
}
