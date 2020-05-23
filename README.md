# Installation

## List S3 bucket

```
aws s3api list-buckets \
    #--query "Buckets[].Name"
```

## Create S3 bucket

```
aws s3api create-bucket --bucket \
    --region eu-west-1 --create-bucket-configuration LocationConstraint=eu-west-1 \
    --acl private
```

## List lambda

```
aws lambda list-functions --max-items 10 \
    #--query 'Functions[].FunctionArn'
```

## Create lambda

```
aws lambda create-function --function-name your-function \
    --zip-file fileb://your-function.zip --handler go-file-name --runtime go1.x \
    --environment "Variables={bucketPathFrom=your-bucket-path-for-upload-image,bucketPathTo=your-bucket-path-for-show-image,env=your-env,table=your-dynamodb-tablename}" \
    --role your-lambda-role \
    #--query 'FunctionArn'
```

## Search lambda

```
aws lambda get-function \
    --function-name s3-rekognition \
    #--query 'Configuration.FunctionArn'
```

## Invoke lambda

```
aws lambda invoke --function-name your-lambda-function-name out --log-type Tail \
    #--query 'LogResult' --output text |  base64 -d
```

## Update lambda code

```
aws lambda update-function-code \
    --function-name  your-lambda-function-name \
    --zip-file fileb://path_to_your_archive.zip
```

## Listen S3 event to Trgigger lambda

```
aws s3api put-bucket-notification-configuration \
    --bucket your-bucket \
    --notification-configuration file://lambda/s3-bucket-notification-configuration.json
```

Extract LambdaFunctionArn from : aws lambda get-function --function-name your-lambda-function-name --query 'Configuration.FunctionArn'

## DynamoDB table list

```
aws dynamodb list-tables
```

## DynamoDB show table informations

```
aws dynamodb describe-table \
    --table-name your-table-name-to-describe
```
