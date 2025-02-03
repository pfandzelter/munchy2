package kaiserstueck

import (
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pfandzelter/munchy2/pkg/food"
)

const price = 890

type kaiserstk struct{}

// New creates a new service to pull the menu from Kaiserstück.
func New() *kaiserstk {
	return &kaiserstk{}
}

func (m *kaiserstk) GetFood(t time.Time) ([]food.Food, error) {
	// download the correct website
	resp, err := http.Get("https://kaiserstueck.de/")

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

	doc.Find(".entry-content").Each(func(i int, s *goquery.Selection) {
		s.Find("p").Each(func(i int, t *goquery.Selection) {
			name := t.Text()
			name = strings.Replace(name, "\n", " ", -1)

			if strings.Contains(name, "Corona") {
				return
			}

			veg := strings.Contains(name, "veg.") || strings.Contains(name, "Veggie")

			foodstuff[name] = food.Food{
				Name:       name,
				StudPrice:  price,
				ProfPrice:  price,
				Vegan:      false,
				Vegetarian: veg,
				Fish:       false,
			}
		})
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
