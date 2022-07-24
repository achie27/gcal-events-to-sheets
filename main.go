package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var WHITELISTED_CHANNEL_ID string

func init() {
	WHITELISTED_CHANNEL_ID = os.Getenv("WHITELISTED_CHANNEL_ID")
}

func main() {
	lambda.Start(handler)
}
