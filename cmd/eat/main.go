package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pfandzelter/munchy2/pkg/dynamo"
	"github.com/pfandzelter/munchy2/pkg/food"
	"github.com/pfandzelter/munchy2/pkg/kaiserstueck"
	"github.com/pfandzelter/munchy2/pkg/personalkantine"
	"github.com/pfandzelter/munchy2/pkg/singh"
	"github.com/pfandzelter/munchy2/pkg/stw"
)

type mensa interface {
	GetFood(date time.Time) ([]food.Food, error)
}

type Canteen struct {
	Name     string
	SpecDiet bool
}

// HandleRequest handles one request to the lambda function.
func HandleRequest(event events.CloudWatchEvent) {

	timezone := os.Getenv("MENSA_TIMEZONE")

	tz, err := time.LoadLocation(timezone)

	if err != nil {
		log.Fatal(err)
	}

	// see if this event was triggered by the DST eventbridge rule
	if strings.Contains(event.Resources[0], "dst") != time.Now().In(tz).IsDST() {
		return
	}

	tablename := os.Getenv("DYNAMODB_TABLE")
	region := os.Getenv("DYNAMODB_REGION")

	db, err := dynamo.New(region, tablename)

	if err != nil {
		log.Fatal(err)
	}

	type Canteen struct {
		Name     string
		SpecDiet bool
	}

	canteens := make(map[Canteen]mensa)

	canteens[Canteen{
		Name:     "Hauptmensa",
		SpecDiet: true,
	}] = stw.New(321)
	canteens[Canteen{
		Name:     "Pasteria Veggie 2.0",
		SpecDiet: true,
	}] = stw.New(631)
	canteens[Canteen{
		Name:     "Marchstr",
		SpecDiet: true,
	}] = stw.New(538)
	// canteens[Canteen{
	// Name:     "Pastaria Architektur",
	// SpecDiet: true,
	// }] = stw.New(540)
	canteens[Canteen{
		Name:     "Kaiserstück",
		SpecDiet: false,
	}] = kaiserstueck.New()
	canteens[Canteen{
		Name:     "Personalkantine",
		SpecDiet: true,
	}] = personalkantine.New()
	canteens[Canteen{
		Name:     "Mathe Café",
		SpecDiet: true,
	}] = singh.New()

	t := time.Now()

	foodlists := make(map[Canteen][]food.Food)

	for c, m := range canteens {
		fl, err := m.GetFood(t)
		if err != nil {
			log.Printf("Error getting food for %s: %s", c.Name, err)
			continue
		}
		log.Printf("Got %d items for %s", len(fl), c.Name)
		foodlists[c] = fl
	}

	for c, f := range foodlists {
		err := db.PutFood(c.Name, c.SpecDiet, f, t)
		if err != nil {
			log.Print(err)
		}
	}

}

func main() {
	lambda.Start(HandleRequest)
}
