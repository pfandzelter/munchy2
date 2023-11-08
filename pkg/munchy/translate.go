package munchy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/pfandzelter/munchy2/pkg/dynamo"
)

type translator struct {
	deepLSourceLang string
	deepLTargetLang string
	deepLURL        string
	deepLKey        string
}

func (t *translator) translate(name string) (string, error) {
	reqBody := "text=" + url.QueryEscape(name)
	reqBody += "&source_lang=" + t.deepLSourceLang
	reqBody += "&target_lang=" + t.deepLTargetLang

	req, err := http.NewRequest("POST", t.deepLURL, bytes.NewBuffer([]byte(reqBody)))

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("DeepL-Auth-Key %s", t.deepLKey))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	s := struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&s)

	if err != nil {
		return "", err
	}

	return s.Translations[0].Text, nil
}

func TranslateFood(f []dynamo.DBEntry, deepLSourceLang string, deepLTargetLang string, deepLURL string, deepLKey string) ([]dynamo.DBEntry, error) {

	if deepLSourceLang == deepLTargetLang {
		return f, nil
	}

	t := &translator{
		deepLSourceLang: deepLSourceLang,
		deepLTargetLang: deepLTargetLang,
		deepLURL:        deepLURL,
		deepLKey:        deepLKey,
	}

	for i, entry := range f {
		for j, item := range entry.Items {
			name, err := t.translate(item.Name)

			if err != nil {
				return nil, err
			}

			f[i].Items[j].Name = name

			log.Printf("translated %s to %s", item.Name, name)
		}
	}

	return f, nil
}
