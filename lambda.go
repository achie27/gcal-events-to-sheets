package main

import (
	"context"
	"errors"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

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
	log.Printf("Event: %+v", request)

	if request.Headers["x-goog-channel-id"] != WHITELISTED_CHANNEL_ID {
		return events.LambdaFunctionURLResponse{StatusCode: 400}, errors.New("Not a whitelisted channel")
	}

	if request.Headers["x-goog-resource-state"] == "sync" {
		return events.LambdaFunctionURLResponse{StatusCode: 200}, nil
	}

	resourceUrl, err := url.Parse(request.Headers["x-goog-resource-uri"])
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: 400}, err
	}

	match := calendarIdRegex.FindStringSubmatch(resourceUrl.Path)

	// TODO: add synctoken
	calEvents, err := gCalSrv.Events.List(match[1]).TimeMin(time.Now().Format(time.RFC3339)).Do()
	if err != nil {
		return events.LambdaFunctionURLResponse{StatusCode: 400}, err
	}

	for _, event := range calEvents.Items {
		eventInfo := extractInfoFromEvent(event, match[1])
		if len(eventInfo.EmailIds) == 0 {
			log.Printf("Event not relevant: %+v\n", eventInfo)
			continue
		}

		sheetRow := [][]interface{}{{eventInfo.Id, eventInfo.Summary, eventInfo.Description, strings.Join(eventInfo.EmailIds, ","), eventInfo.Start, eventInfo.Organizer, eventInfo.CalendarId, eventInfo.Created, strings.Join(eventInfo.Entities, ",")}}

		_, err = gSheetsSrv.Spreadsheets.Values.Append(SPREADSHEET_ID, "Sheet1", &sheets.ValueRange{Values: sheetRow}).InsertDataOption("INSERT_ROWS").ValueInputOption("RAW").Do()
		if err != nil {
			log.Printf("Sheets erred: %v", err)
			return events.LambdaFunctionURLResponse{StatusCode: 500}, err
		}
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
