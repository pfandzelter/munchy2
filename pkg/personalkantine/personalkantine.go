package personalkantine

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/pfandzelter/munchy2/pkg/food"
)

var blacklist = [...]string{
	"Gemüseplatte",
}

type kantine struct{}

// New creates a new service to pull the menu from Personalkantine.
func New() *kantine {
	return &kantine{}
}

func checkBlacklist(name string) bool {
	for _, item := range blacklist {
		if strings.Contains(name, item) {
			return true
		}
	}

	return false
}

func (m *kantine) GetFood(t time.Time) ([]food.Food, error) {
	// get today's date
	date := t.Format("02.01.2006")

	// download the correct website
	resp, err := http.Get("https://personalkantine.personalabteilung.tu-berlin.de/")

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

	doc.Find("main > #wochenkarte > div > div > .Menu__accordion > li").Each(func(i int, t *goquery.Selection) {
		if strings.Contains(t.Find("h2").Text(), date) {

			t.Find("ul > li").Each(func(i int, s *goquery.Selection) {

				name := s.Find("h4").Text()

				// get the text of the main element, but without children
				description := s.Contents().NotSelection(s.Children()).Text()

				// remove all bracketed text from the description using regex
				bracketed := regexp.MustCompile(`\([^)]*\)`)
				description = bracketed.ReplaceAllString(description, "")

				// trim whitespaces from the description
				description = strings.TrimSpace(description)

				name = fmt.Sprintf("%s (%s)", name, description)

				if checkBlacklist(name) {
					return
				}

				price := s.Find(".price").Text()
				price = strings.Replace(price, "\n", "", -1)
				price = strings.Replace(price, " ", "", -1)
				price = strings.Replace(price, "€", "", -1)
				price = strings.Replace(price, "&euro;", "", -1)
				price = strings.Replace(price, ",", "", -1)

				endprice, err := strconv.Atoi(price)

				if err != nil {
					return
				}

				vegetarian := strings.Contains(s.Text(), "(v)") || strings.Contains(s.Text(), "(V)")
				fish := strings.Contains(s.Text(), "(F)")

				foodstuff[name] = food.Food{
					Name:       name,
					StudPrice:  endprice,
					ProfPrice:  endprice,
					Vegan:      false,
					Vegetarian: vegetarian,
					Fish:       fish,
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
