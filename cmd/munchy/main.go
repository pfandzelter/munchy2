package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/pfandzelter/munchy2/pkg/dynamo"
	"github.com/pfandzelter/munchy2/pkg/munchy"
)

var webhookURL = os.Getenv("WEBHOOK_URL")
var awsRegion = os.Getenv("DYNAMODB_REGION")
var awsTable = os.Getenv("DYNAMODB_TABLE")
var deepLTargetLang = os.Getenv("DEEPL_TARGET_LANG")
var deepLURL = os.Getenv("DEEPL_URL")
var deepLKey = os.Getenv("DEEPL_KEY")
var deepLSourceLang = "DE"

var longMsg = "Today is " + time.Now().Weekday().String() + ", the *" + time.Now().Format("01/02/2006") + "*, here is today's lunch menu.\n*Enjoy!* :drooling_face:"
var shortMsg = "Here is today's lunch menu!"

// HandleRequest handles one request to the Lambda function.
func HandleRequest(ctx context.Context, event events.CloudWatchEvent) {

	timezone := os.Getenv("MENSA_TIMEZONE")

	tz, err := time.LoadLocation(timezone)

	if err != nil {
		log.Fatal(err)
	}

	// see if this event was triggered by the DST eventbridge rule
	if strings.Contains(event.Resources[0], "dst") != time.Now().In(tz).IsDST() {
		return
	}

	f, err := dynamo.GetFood(awsRegion, awsTable)

	if err != nil {
		log.Fatalf("Error getting food from DynamoDB: %v", err)
	}

	f, err = munchy.TranslateFood(f, deepLSourceLang, deepLTargetLang, deepLURL, deepLKey)

	if err != nil {
		log.Fatalf("Error translating food: %v", err)
	}

	msg := ""

	// every day is English wednesday
	msg = munchy.GetMessage(f, longMsg, shortMsg)

	jsonStr := []byte(msg)
	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonStr))

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	log.Printf("sending %s to %s, got %d: %s", msg, webhookURL, resp.StatusCode, string(data))
}

func main() {
	lambda.Start(HandleRequest)
}
