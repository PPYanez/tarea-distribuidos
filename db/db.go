package db

import (
	"context"
	"distribuidos/tarea-1/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	ctx    context.Context
)

func init() {
	// initialize the MongoDB client
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		panic(err)
	}

	// establish a connection to the MongoDB server
	err = client.Connect(context.Background())
	if err != nil {
		panic(err)
	}

	// check the connection
	connected := pingDb(client)
	if !connected {
		panic("Failed to connect to database")
	}

}

func pingDb(client *mongo.Client) bool {
	var result bson.M
	if err := client.Database("distribuidos").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return false
	}
	return true
}

func GetClient() *mongo.Client {
	return client
}

// Static data
func GetAncillaries() map[string]models.AncillarieData {
	ancillariesData := map[string]models.AncillarieData{
		"1": {Nombre: "Equipaje de mano", Precio: 10000, Ssr: "BGH"},
		"2": {Nombre: "Equipaje de bodega", Precio: 30000, Ssr: "BGR"},
		"3": {Nombre: "Asiento", Precio: 5000, Ssr: "STDF"},
		"4": {Nombre: "Embarque y Check In prioritario", Precio: 2000, Ssr: "PAXS"},
		"5": {Nombre: "Mascota en cabina", Precio: 40000, Ssr: "PTCR"},
		"6": {Nombre: "Mascota en bodega", Precio: 40000, Ssr: "AVIH"},
		"7": {Nombre: "Equipaje especial", Precio: 35000, Ssr: "SPML"},
		"8": {Nombre: "Acceso a Salón VIP", Precio: 15000, Ssr: "LNGE"},
		"9": {Nombre: "Wi-Fi a bordo", Precio: 20000, Ssr: "WIFI"},
	}

	return ancillariesData
}
