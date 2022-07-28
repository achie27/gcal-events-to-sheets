package main

import (
	"context"
	"errors"
	"log"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/sheets/v4"
)

type EventInfo struct {
	Id          string
	EmailIds    []string
	Entities    []string
	Summary     string
	Description string
	Organizer   string
	Start       string
	Created     string
	CalendarId  string
}

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

	eventInfo := extractInfoFromEvent(event, match[1])
	if len(eventInfo.EmailIds) == 0 {
		log.Printf("Event not relevant: %+v\n", eventInfo)
		return events.LambdaFunctionURLResponse{StatusCode: 200}, nil
	}

	sheetRow := [][]interface{}{{eventInfo.Id, eventInfo.Summary, eventInfo.Description, eventInfo.EmailIds, eventInfo.Start, eventInfo.Organizer, eventInfo.CalendarId, eventInfo.Created, eventInfo.Entities}}

	_, err = gSheetsSrv.Spreadsheets.Values.Append(SPREADSHEET_ID, "Sheet1", &sheets.ValueRange{Values: sheetRow}).InsertDataOption("INSERT_ROWS").Do()
	if err != nil {
		log.Printf("Sheets erred: %v", err)
		return events.LambdaFunctionURLResponse{StatusCode: 500}, err
	}

	return events.LambdaFunctionURLResponse{StatusCode: 200}, nil
}

func extractInfoFromEvent(event *calendar.Event, calendarId string) *EventInfo {
	organizer, start := "", ""

	if event.Organizer != nil {
		organizer = event.Organizer.Email
	} else if event.Creator != nil {
		organizer = event.Creator.Email
	}

	if event.Start != nil {
		start = event.Start.DateTime
	}

	extractedEmails := []string{}
	extractedEntities := []string{}

	for _, attn := range event.Attendees {
		extractedEmails = append(extractedEmails, attn.Email)

		if attn.DisplayName != "" {
			extractedEntities = append(extractedEntities, attn.DisplayName)
		}
	}
	extractedEmails = filterEmails(extractedEmails)

	eventText := event.Summary + "\n" + event.Description
	extractedEmails = append(extractedEmails, extractAndFilterEmailsFromText(eventText)...)
	extractedEntities = append(extractedEntities, extractEntitiesFromText(eventText)...)

	return &EventInfo{
		Id:          event.Id,
		EmailIds:    extractedEmails,
		Entities:    extractedEntities,
		Summary:     event.Summary,
		Description: event.Description,
		Created:     event.Created,
		Organizer:   organizer,
		Start:       start,
		CalendarId:  calendarId,
	}
}
