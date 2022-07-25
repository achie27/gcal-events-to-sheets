package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

func handler(ctx context.Context, request *events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	if request.Headers["X-Goog-Channel-ID"] != WHITELISTED_CHANNEL_ID {
		return events.LambdaFunctionURLResponse{StatusCode: 400}, errors.New("Not it")
	}

	log.Println(request)

	return events.LambdaFunctionURLResponse{StatusCode: 200}, nil
}
