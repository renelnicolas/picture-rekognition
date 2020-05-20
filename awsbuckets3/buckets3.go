package awsbuckets3

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// AwsError :
type AwsError struct {
	Code string
	Err  string
}

// CopyObject :
func CopyObject(svc *s3.S3, bucketName, objectKey, toDir string) *AwsError {
	source := bucketName + "/" + objectKey

	input := &s3.CopyObjectInput{
		Bucket:     aws.String(bucketName),
		CopySource: aws.String(source),
		Key:        aws.String(toDir),
	}

	_, err := svc.CopyObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return &AwsError{
				Code: aerr.Code(),
				Err:  aerr.Error(),
			}
		}

		return &AwsError{
			Code: "UKN",
			Err:  err.Error(),
		}
	}

	return nil
}

// GetObjectURL :
func GetObjectURL(svc *s3.S3, bucketName, objectKey string, secondExpiration time.Duration) (string, *AwsError) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	req, _ := svc.GetObjectRequest(params)

	url, err := req.Presign(secondExpiration * time.Second) // Set link expiration time
	if err != nil {
		return "", &AwsError{
			Code: "UKN",
			Err:  "[AWS GET LINK]: " + err.Error(),
		}
	}

	return url, nil
}

// DeleteObject :
func DeleteObject(svc *s3.S3, bucketName, objectKey string) *AwsError {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	_, err := svc.DeleteObject(params)
	if err != nil {
		return &AwsError{
			Code: "UKN",
			Err:  "Unable to delete object from bucket :" + err.Error(),
		}
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		return &AwsError{
			Code: "UKN",
			Err:  "Unable to delete object from bucket :" + err.Error(),
		}
	}

	return nil
}

// UploadObject :
func UploadObject(sess *session.Session, fileLocalDir, bucketName, objectKey string) *AwsError {
	file, err := os.Open(fileLocalDir)
	if err != nil {
		return &AwsError{
			Code: "UKN",
			Err:  "UploadObject - Open :" + err.Error(),
		}
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	size := fileInfo.Size()

	fmt.Println("UploadObject > fileInfo", fileInfo)
	fmt.Println("UploadObject > size", size)

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})

	if nil != err {
		return &AwsError{
			Code: "UKN",
			Err:  "UploadObject - PutObject :" + err.Error(),
		}
	}

	return nil
}
