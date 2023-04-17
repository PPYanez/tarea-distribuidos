package utilities

import (
	"context"
	"distribuidos/tarea-1/db"
	"distribuidos/tarea-1/models"
	"log"
	"math/rand"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GenPNR() models.PNR {
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
			if pnrInDB.Pnr == string(generatedPNR) {
				pnrIsNew = false
				break
			}
		}

		if pnrIsNew {
			break
		}
	}

	return models.PNR{Pnr: string(generatedPNR)}
}

func getExistingPNR() []models.PNR {
	client := db.GetClient()

	var pnrCollection = client.Database("distribuidos").Collection("PNR")
	findOptions := options.Find()
	var results []models.PNR

	cur, err := pnrCollection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem models.PNR
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
