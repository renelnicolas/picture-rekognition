package awsdynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// AwsError :
type AwsError struct {
	Code string
	Err  string
}

var (
	dynamoDBConn *dynamodb.DynamoDB
)

// CreateDynamoDBConn :
func CreateDynamoDBConn(svc *dynamodb.DynamoDB) *dynamodb.DynamoDB {
	if nil == svc {
		if nil == dynamoDBConn {
			dynamoDBConn = dynamodb.New(session.New())
			return dynamoDBConn
		}
	}

	return svc
}

// CreateTable :
func CreateTable(svc *dynamodb.DynamoDB, tableName string, attributeDefinitions []*dynamodb.AttributeDefinition, keySchemaElement []*dynamodb.KeySchemaElement, provisionedThroughput dynamodb.ProvisionedThroughput) *AwsError {
	params := &dynamodb.CreateTableInput{
		TableName:             aws.String(tableName),
		AttributeDefinitions:  attributeDefinitions,
		KeySchema:             keySchemaElement,
		ProvisionedThroughput: &provisionedThroughput,
	}

	_, err := CreateDynamoDBConn(svc).CreateTable(params)
	if err != nil {
		return &AwsError{
			Code: "UKN",
			Err:  err.Error(),
		}
	}

	return nil
}

// TableExist :
func TableExist(svc *dynamodb.DynamoDB, tableName string) (*string, *AwsError) {
	input := &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	}

	result, err := CreateDynamoDBConn(svc).DescribeTable(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return nil, &AwsError{
				Code: aerr.Code(),
				Err:  aerr.Error(),
			}
		}

		return nil, &AwsError{
			Code: "UKN",
			Err:  err.Error(),
		}
	}

	return result.Table.TableStatus, nil
}

// InsertInto :
func InsertInto(svc *dynamodb.DynamoDB, tableName string, item interface{}) *AwsError {
	damp, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return &AwsError{
			Code: "UKN",
			Err:  err.Error(),
		}
	}

	input := &dynamodb.PutItemInput{
		Item:      damp,
		TableName: aws.String(tableName),
	}

	_, err = CreateDynamoDBConn(svc).PutItem(input)
	if err != nil {
		return &AwsError{
			Code: "UKN",
			Err:  "Got error calling PutItem: " + err.Error(),
		}
	}

	return nil
}

// FindAll :
func FindAll(svc *dynamodb.DynamoDB, tableName string) ([]RekognitionItem, *AwsError) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc = dynamodb.New(sess)

	proj := expression.NamesList(
		expression.Name("pk"),
		expression.Name("url"),
		expression.Name("keywords"),
	)

	// expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()
	expr, err := expression.NewBuilder().WithProjection(proj).Build()
	if err != nil {
		return nil, &AwsError{
			Code: "UKN",
			Err:  "FindAll > Got error building expression: " + err.Error(),
		}
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(tableName),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	if err != nil {
		return nil, &AwsError{
			Code: "UKN",
			Err:  "FindAll > Query API call failed: " + err.Error(),
		}
	}

	items := []RekognitionItem{}

	for _, i := range result.Items {
		item := RekognitionItem{}

		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, &AwsError{
				Code: "UKN",
				Err:  "FindAll > Got error unmarshalling: " + err.Error(),
			}
		}

		items = append(items, item)
	}

	return items, nil
}
