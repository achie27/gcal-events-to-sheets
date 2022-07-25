package main

import (
	"context"
	"errors"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
)

var calendarIdRegex = regexp.MustCompile(`calendar/v3/calendars/\s*(.*?)\s*/events`)

func handler(ctx context.Context, request *events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	if request.Headers["X-Goog-Channel-ID"] != WHITELISTED_CHANNEL_ID {
		return events.LambdaFunctionURLResponse{StatusCode: 400}, errors.New("Not it")
	}

	match := calendarIdRegex.FindStringSubmatch(request.Headers["X-Goog-Resource-URI"])
	event, err := gCalSrv.Events.Get(match[1], request.Headers["X-Goog-Resource-ID"]).Do()
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: 400}, errors.New("Not it")
	}

	extractedEmails := []string{}
	extractedEntities := []string{}

	for _, attn := range event.Attendees {
		extractedEmails = append(extractedEmails, attn.Email)

		if attn.DisplayName != "" {
			extractedEntities = append(extractedEntities, attn.DisplayName)
		}
	}

	// eventText := strings.Join(event.Summary, " ", event.Description)

	return events.LambdaFunctionURLResponse{StatusCode: 200}, nil
}
