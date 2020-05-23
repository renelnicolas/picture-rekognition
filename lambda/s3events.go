package main

// ./buid.sh s3events s3-rekognition dev-img-rekognition

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/aws/aws-sdk-go/service/s3"

	"ohmytech.io/picture-rekognition/awsbuckets3"
	"ohmytech.io/picture-rekognition/awsdynamodb"
)

// label :
type label struct {
	Confidence float64 `json:"confidence"`
	Label      string  `json:"label"`
}

// labels :
type labels []label

// appError :
type appError struct {
	Code string
	Err  string
}

// tableName :
var (
	tableName = os.Getenv("table")
)

// hash :
func hash(toHash []byte) string {
	h := sha1.New()

	h.Write(toHash)

	return hex.EncodeToString(h.Sum(nil))
}

// extractLabels :
func extractLabels(result *rekognition.DetectLabelsOutput) labels {
	lbls := labels{}

	for _, lbl := range result.Labels {
		lbls = append(lbls, label{
			Confidence: *lbl.Confidence,
			Label:      *lbl.Name,
		})
	}

	return lbls
}

// handler :
func handler(ctx context.Context, s3Event events.S3Event) error {
	for _, record := range s3Event.Records {
		rs3 := record.S3

		svc := rekognition.New(session.New())
		input := &rekognition.DetectLabelsInput{
			Image: &rekognition.Image{
				S3Object: &rekognition.S3Object{
					Bucket: aws.String(rs3.Bucket.Name),
					Name:   aws.String(rs3.Object.Key),
				},
			},
			MaxLabels:     aws.Int64(4),
			MinConfidence: aws.Float64(70.000000),
		}

		result, err := svc.DetectLabels(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				return fmt.Errorf("DetectLabels > Code : %s >> %s", aerr.Code(), aerr.Error())
			}

			return fmt.Errorf("DetectLabels > Code : %s >> %s", "UKN", err.Error())
		}

		lbls := extractLabels(result)
		if nil != lbls && 0 < len(lbls) {
			svc := awsdynamodb.CreateDynamoDBConn(nil)

			if _, err := awsdynamodb.TableExist(svc, tableName); nil != err {
				attributeDefinitions := []*dynamodb.AttributeDefinition{
					{
						AttributeName: aws.String("pk"),
						AttributeType: aws.String("S"),
					},
					// {
					// 	AttributeName: aws.String("sk"),
					// 	AttributeType: aws.String("S"),
					// },
				}

				keySchema := []*dynamodb.KeySchemaElement{
					{
						AttributeName: aws.String("pk"),
						KeyType:       aws.String("HASH"),
					},
					// {
					// 	AttributeName: aws.String("sk"),
					// 	KeyType:       aws.String("RANGE"),
					// },
				}

				provisionedThroughput := dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(5),
					WriteCapacityUnits: aws.Int64(5),
				}

				if err = awsdynamodb.CreateTable(svc, tableName, attributeDefinitions, keySchema, provisionedThroughput); nil != err {
					return fmt.Errorf("createTable > Code : %s >> %s", err.Code, err.Err)
				}

				counter := 0

				tableStatus, err := awsdynamodb.TableExist(svc, tableName)
				for "ACTIVE" != *tableStatus && nil == err {
					time.Sleep(3 * time.Second)

					tableStatus, err = awsdynamodb.TableExist(svc, tableName)

					counter++

					if 6 < counter {
						return errors.New("Database creation too long")
					}
				}
			}

			s3Svc := s3.New(session.New())

			toDir := strings.Replace(rs3.Object.Key, os.Getenv("bucketPathFrom"), os.Getenv("bucketPathTo"), 1)

			awsErr := awsbuckets3.CopyObject(s3Svc, rs3.Bucket.Name, rs3.Object.Key, toDir)
			if nil != awsErr {
				return fmt.Errorf("CopyObject > Code : %s >> %s", awsErr.Code, awsErr.Err)
			}

			URL, awsErr := awsbuckets3.GetObjectURL(s3Svc, rs3.Bucket.Name, toDir, 604799)
			if nil != awsErr {
				return fmt.Errorf("GetObjectURL > Code : %s >> %s", awsErr.Code, awsErr.Err)
			}

			awsErr = awsbuckets3.DeleteObject(s3Svc, rs3.Bucket.Name, rs3.Object.Key)
			if nil != awsErr {
				return fmt.Errorf("DeleteObject > Code : %s >> %s", awsErr.Code, awsErr.Err)
			}

			image := strings.Replace(rs3.Object.Key, os.Getenv("bucketPathFrom")+"/", "", 1)

			item := awsdynamodb.RekognitionItem{
				SortKey:    time.Now().Format("2006#01#02#15#04#05"),
				Keywords:   lbls,
				URL:        URL,
				PrimaryKey: hash([]byte(image)),
			}

			if aperr := awsdynamodb.InsertInto(svc, tableName, item); nil != aperr {
				return fmt.Errorf("insertIntoTable > Code : %s >> %s", aperr.Code, aperr.Err)
			}
		}
	}

	return nil
}

// main :
func main() {
	lambda.Start(handler)
}
