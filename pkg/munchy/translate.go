package munchy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error translating %s: %s", name, resp.Status)
	}

	defer resp.Body.Close()

	s := struct {
		Translations []struct {
			Text string `json:"text"`
		} `json:"translations"`
	}{}

	b, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	err = json.Unmarshal(b, &s)

	if err != nil {
		log.Printf("error unmarshalling %s: %v", string(b), err)
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
				log.Printf("error translating %s: %v", item.Name, err)
				return nil, err
			}

			f[i].Items[j].Name = name

			log.Printf("translated %s to %s", item.Name, name)

			time.Sleep(500 * time.Millisecond)
		}
	}

	return f, nil
}
