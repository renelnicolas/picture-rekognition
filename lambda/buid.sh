#!/bin/bash

rm -rf bin/$1
rm -rf bin/$1.zip

env GOOS=linux GOARCH=amd64 go build -ldflags="-d -s -w" -o bin/$1 $1.go

chmod +x bin/$1

zip -j bin/$1.zip bin/$1

aws lambda update-function-code \
    --function-name  $2 \
    --zip-file fileb://bin/$1.zip

sleep 8

echo "aws s3 cp loup.jpeg s3://$3/waiting/loup.jpeg"
aws s3 cp loup.jpeg s3://$3/waiting/
