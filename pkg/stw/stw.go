package stw

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/pfandzelter/munchy2/pkg/food"
)

// Mensen:
// 321 - TU Hardenbergstr
// 631 - TU Pastaria Veggie 2.0
// 538 - TU Marchstr
// 540 - TU Pastaria Architektur
type mensa struct {
	id int
}

var blocklist = [...]string{
	"kuchen",
	"creme",
	"torte",
	"Brownie",
	"Apfeltasche",
	"Gouda",
	"Veganer Schmelz gerieben",
	"Hartkäse gerieben",
}

// New creates a new service to pull the menu for an STW Mensa based on an id.
func New(id int) *mensa {
	return &mensa{
		id: id,
	}
}

func checkblocklist(name string) bool {
	for _, item := range blocklist {
		if strings.Contains(strings.ToUpper(name), strings.ToUpper(item)) {
			return true
		}
	}

	return false
}

func (m *mensa) GetFood(t time.Time) ([]food.Food, error) {
	// get today's date
	date := t.Format("2006-01-02")

	// download the correct website
	// should be something like:
	// $ curl 'https://www.stw.berlin/xhr/speiseplan-wochentag.html' -v --data 'resources_id=321&date=2020-02-21' --compressed
	data := []byte(fmt.Sprintf("resources_id=%d&date=%s", m.id, date))

	resp, err := http.Post("https://www.stw.berlin/xhr/speiseplan-wochentag.html",
		"application/x-www-form-urlencoded", bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// parse the results
	foodstuff := make(map[string]food.Food)

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, err
	}

	doc.Find(".splGroupWrapper").Each(func(i int, t *goquery.Selection) {
		if t.Find(".row > .splGroup").Text() == "Aktionen" || t.Find(".row > .splGroup").Text() == "Essen" {
			t.Find(".splMeal").Each(func(i int, s *goquery.Selection) {
				name := s.Find("div > .bold").Text()

				if checkblocklist(name) {
					return
				}

				price := s.Find(".col-xs-12.col-md-3.text-right").Text()
				price = strings.Replace(price, "\n", "", -1)
				price = strings.Replace(price, " ", "", -1)
				price = strings.Replace(price, "€", "", -1)
				price = strings.Replace(price, "&euro;", "", -1)
				price = strings.Replace(price, ",", "", -1)

				prices := strings.Split(price, "/")

				studPrice, err := strconv.Atoi(prices[0])

				if err != nil {
					return
				}

				profPrice, err := strconv.Atoi(prices[1])

				if err != nil {
					return
				}

				vegetarian := false
				vegan := false
				fish := false
				climate := false

				s.Find("div > .splIcon").Each(func(i int, x *goquery.Selection) {
					src, ok := x.Attr("src")

					if !ok {
						return
					}

					if src == "/vendor/infomax/mensen/icons/15.png" {
						vegan = true
						return
					}

					if src == "/vendor/infomax/mensen/icons/1.png" {
						vegetarian = true
						return
					}

					if src == "/vendor/infomax/mensen/icons/38.png" {
						fish = true
						return
					}

					if src == "/vendor/infomax/mensen/icons/43.png" {
						climate = true
						return
					}
				})

				// we only check for the fishing symbol, for unfair fishing that doesn't appear
				// we can alternatively check for the "24" allergen, meaning fish or "22" meaning sea food

				if !fish {
					allergens := s.Find("div > .kennz").Text()
					fish = strings.Contains(allergens, "24") || strings.Contains(allergens, "22")
				}

				foodstuff[name] = food.Food{
					Name:       name,
					StudPrice:  studPrice,
					ProfPrice:  profPrice,
					Vegan:      vegan,
					Vegetarian: vegetarian,
					Fish:       fish,
					Climate:    climate,
				}
			})
		}
	})

	// return stuff
	foodlist := make([]food.Food, len(foodstuff))
	i := 0

	for _, f := range foodstuff {
		foodlist[i] = f
		i++
	}

	return foodlist, nil

}
