package munchy

import (
	"fmt"
	"strings"

	"github.com/pfandzelter/munchy2/pkg/dynamo"
)

func fmtPrice(p int) string {
	cents := p % 100
	euros := int((p - cents) / 100)
	return fmt.Sprintf("%d,%02d", euros, cents)
}

// you should never do this
func escape(s string) string {
	s = strings.ReplaceAll(s, `"`, `\"`)

	return s
}

func GetMessage(f []dynamo.DBEntry, long string, short string) string {

	foods := ""

	for _, entry := range f {
		if len(entry.Items) == 0 {
			continue
		}

		foods += fmt.Sprintf(`,{"type": "section","text": {"type": "mrkdwn","text": "*%s*`, escape(entry.Canteen))

		for _, item := range entry.Items {
			foods += "\n"

			if !entry.SpecDiet {
				foods += "• "
			} else if item.Vegan {
				foods += ":seedling:"
			} else if item.Vegetarian {
				foods += ":cheese_wedge:"
			} else if item.Fish {
				foods += ":fish:"
			} else {
				foods += ":cut_of_meat:"
			}

			foods += escape(item.Name)

			if item.Climate {
				foods += " (:earth_africa:)"
			}

			foods += " "

			foods += fmtPrice(item.StudPrice)
			foods += "€"

			if item.StudPrice != item.ProfPrice {
				foods += "/" + fmtPrice(item.ProfPrice) + "€"
			}

		}
		foods += `"}}`
	}

	return fmt.Sprintf(`{"text":"%s", "blocks": [{"type": "section","text": {"type": "mrkdwn","text": "%s"}},{"type": "divider"}%s]}`, short, long, foods)
}
