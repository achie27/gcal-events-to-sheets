package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func getGoogleClient(ctx context.Context) *http.Client {
	googleConfig := []byte(fmt.Sprintf(`
		{
			"installed": {
				"client_id": "%s",
				"project_id": "%s",
				"auth_uri": "https://accounts.google.com/o/oauth2/auth",
				"token_uri": "https://oauth2.googleapis.com/token",
				"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
				"client_secret": "%s",
				"redirect_uris": [
					"%s"
				]
			}
		}
	`, GOOGLE_CONFIG_CLIENT_ID, GOOGLE_CONFIG_PROJECT_ID, GOOGLE_CONFIG_CLIENT_SECRET, GOOGLE_CONFIG_REDIRECT_URI))

	config, err := google.ConfigFromJSON(googleConfig, calendar.CalendarReadonlyScope, sheets.SpreadsheetsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	expiry, err := time.Parse(time.RFC3339Nano, GOOGLE_OAUTH_EXPIRY)
	if err != nil {
		log.Fatalf("Incorrect expiry: %v", err)
	}

	return config.Client(ctx, &oauth2.Token{
		AccessToken:  GOOGLE_OAUTH_ACCESS_TOKEN,
		TokenType:    GOOGLE_OAUTH_TOKEN_TYPE,
		RefreshToken: GOOGLE_OAUTH_REFRESH_TOKEN,
		Expiry:       expiry,
	})
}

func getCalendarService(ctx context.Context) *calendar.Service {
	calendarSrv, err := calendar.NewService(ctx, option.WithHTTPClient(gClient))
	if err != nil {
		log.Fatalf("Unable to create Calendar service %v", err)
	}

	return calendarSrv
}

func getSheetsService(ctx context.Context) *sheets.Service {
	sheetsSrv, err := sheets.NewService(ctx, option.WithHTTPClient(gClient))
	if err != nil {
		log.Fatalf("Unable to create Sheets service %v", err)
	}

	return sheetsSrv
}
