package awsdynamodb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// AwsError :
type AwsError struct {
	Code string
	Err  string
}

// CreateTable :
func CreateTable(svc *dynamodb.DynamoDB, tableName string, attributeDefinitions []*dynamodb.AttributeDefinition, keySchemaElement []*dynamodb.KeySchemaElement, provisionedThroughput dynamodb.ProvisionedThroughput) *AwsError {
	params := &dynamodb.CreateTableInput{
		TableName:             aws.String(tableName),
		AttributeDefinitions:  attributeDefinitions,
		KeySchema:             keySchemaElement,
		ProvisionedThroughput: &provisionedThroughput,
	}

	_, err := svc.CreateTable(params)
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

	result, err := svc.DescribeTable(input)
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

	_, err = svc.PutItem(input)
	if err != nil {
		return &AwsError{
			Code: "UKN",
			Err:  "Got error calling PutItem: " + err.Error(),
		}
	}

	return nil
}
