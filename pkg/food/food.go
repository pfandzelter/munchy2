package food

import "fmt"

// Food is one food item at a canteen.
type Food struct {
	Name       string `json:"name"`
	StudPrice  int    `json:"studprice"`
	ProfPrice  int    `json:"profprice"`
	Vegan      bool   `json:"vgn"`
	Vegetarian bool   `json:"vgt"`
	Fish       bool   `json:"fish"`
	Climate    bool   `json:"climate"`
}

func (f Food) String() string {
	return fmt.Sprintf("%s: %dEUR/%dEUR, vegan: %t, vegetarian: %t, fish: %t, climate: %t", f.Name, f.StudPrice, f.ProfPrice, f.Vegan, f.Vegetarian, f.Fish, f.Climate)
}
