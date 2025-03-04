package singh

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pfandzelter/munchy2/pkg/food"
)

type singh struct{}

// New creates a new service to pull the menu from Mathe Cafe.
func New() *singh {
	return &singh{}
}

func (m *singh) GetFood(t time.Time) ([]food.Food, error) {
	// get today's date
	date := t.Weekday().String()

	switch date {
	case "Monday":
		date = "MONTAG"
	case "Tuesday":
		date = "DIENSTAG"
	case "Wednesday":
		date = "MITTWOCH"
	case "Thursday":
		date = "DONNERSTAG"
	case "Friday":
		date = "FREITAG"
	}

	// download the correct website
	resp, err := http.Get("http://mathe-cafe-tu.de/cafe/")

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

	doc.Find(".content_wrapper.clearfix").Each(func(i int, s *goquery.Selection) {
		// println(s.Text())
		s.Find(".wrap.mcb-wrap.one-second.valign-top.clearfix").Each(func(i int, m *goquery.Selection) {

			// if the name of the day does not appear, we have the wrong menu
			if !strings.Contains(m.Text(), date) {
				// log.Printf("%s is not in %s", date, m.Text())
				return
			}
			// print("have something")

			var nextVeg bool
			var nextVgn bool

			items := m.Find(".column.mcb-column.one.column_column.column-margin-")

			for c := range items.Nodes {
				curr := items.Eq(c)
				name := curr.Text()
				println(name)
				name = strings.Replace(name, "\n", " ", -1)

				if strings.Contains(name, date) {
					continue
				}

				if strings.Contains(name, "VEGETARISCH") {
					nextVeg = true
					continue
				}

				if strings.Contains(name, "VEGAN") {
					nextVgn = true
					continue
				}

				// find the description
				desc := curr.Find("p").Text() + " (" + curr.Find("th").First().Text() + ")"
				log.Printf("%s", desc)

				// find the price
				endprice := 999
				for p := curr.Find("th").First(); len(p.Nodes) != 0; p = p.Next() {
					price := p.Text()

					if !strings.Contains(price, "€") {
						continue
					}

					price = strings.Replace(price, "\n", "", -1)
					price = strings.Replace(price, " ", "", -1)
					price = strings.Replace(price, "€", "", -1)
					price = strings.Replace(price, "&euro;", "", -1)
					price = strings.Replace(price, ",", "", -1)
					price = strings.Replace(price, ".", "", -1)

					endprice, err = strconv.Atoi(price)

					if err != nil {
						return
					}

					break
				}

				foodstuff[name] = food.Food{
					Name:       desc,
					StudPrice:  endprice,
					ProfPrice:  endprice,
					Vegan:      nextVgn,
					Vegetarian: nextVeg,
					Fish:       false,
				}
				nextVgn = false
				nextVeg = false
			}
		})
	})

	// remove vegetarian and vegan chicken... thanks, Mathe Cafe
	for k, v := range foodstuff {
		if strings.Contains(k, "Hähnchen") {
			v.Vegetarian = false
			v.Vegan = false
			foodstuff[k] = v
		}

		if strings.Contains(k, "SAMOSA") {
			v.Vegetarian = true
			v.Vegan = true
			foodstuff[k] = v
		}
	}

	// return stuff
	foodlist := make([]food.Food, len(foodstuff))
	i := 0

	for _, f := range foodstuff {
		foodlist[i] = f
		i++
	}

	return foodlist, nil

}
