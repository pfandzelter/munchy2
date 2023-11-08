package dynamo

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/pfandzelter/munchy2/pkg/food"
)

// DB is a DynamoDB service for a particular table.
type DB struct {
	dynamodb *dynamodb.DynamoDB
	table    string
}

// DBEntry is the entry in our DynamoDB table for a particular day.
type DBEntry struct {
	Canteen  string      `json:"canteen"`
	SpecDiet bool        `json:"spec_diet"`
	Date     string      `json:"date"`
	Items    []food.Food `json:"items"`
}

// New creates a new DynamoDB session.
func New(region string, table string) (*DB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return nil, err
	}

	return &DB{
		dynamodb: dynamodb.New(sess),
		table:    table,
	}, nil
}

// PutFood puts one food item into the DynamoDB table.
func (d *DB) PutFood(c string, specdiet bool, f []food.Food, t time.Time) error {
	item := struct {
		Canteen  string      `json:"canteen"`
		SpecDiet bool        `json:"spec_diet"`
		Date     string      `json:"date"`
		Items    []food.Food `json:"items"`
	}{
		Canteen:  c,
		SpecDiet: specdiet,
		Date:     t.Format("2006-01-02"),
		Items:    f,
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: &d.table,
	}

	_, err = d.dynamodb.PutItem(input)

	if err != nil {
		return err
	}

	return nil
}

func GetFood(region string, table string) ([]DBEntry, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	svc := dynamodb.New(sess)

	items := []DBEntry{}

	date := time.Now().Format("2006-01-02")

	filt := expression.Name("date").Equal(expression.Value(date))

	proj := expression.NamesList(expression.Name("canteen"), expression.Name("date"), expression.Name("spec_diet"), expression.Name("items"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		return nil, err
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(table),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)

	if err != nil {
		return nil, err
	}

	for _, i := range result.Items {
		item := DBEntry{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
