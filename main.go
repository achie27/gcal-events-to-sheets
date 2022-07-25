package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/sheets/v4"
)

var GOOGLE_CONFIG_CLIENT_ID = os.Getenv("GOOGLE_CONFIG_CLIENT_ID")
var GOOGLE_CONFIG_PROJECT_ID = os.Getenv("GOOGLE_CONFIG_PROJECT_ID")
var GOOGLE_CONFIG_CLIENT_SECRET = os.Getenv("GOOGLE_CONFIG_CLIENT_SECRET")
var GOOGLE_CONFIG_REDIRECT_URI = os.Getenv("GOOGLE_CONFIG_REDIRECT_URI")
var GOOGLE_OAUTH_ACCESS_TOKEN = os.Getenv("GOOGLE_OAUTH_ACCESS_TOKEN")
var GOOGLE_OAUTH_TOKEN_TYPE = os.Getenv("GOOGLE_OAUTH_TOKEN_TYPE")
var GOOGLE_OAUTH_REFRESH_TOKEN = os.Getenv("GOOGLE_OAUTH_REFRESH_TOKEN")
var GOOGLE_OAUTH_EXPIRY = os.Getenv("GOOGLE_OAUTH_EXPIRY")
var WHITELISTED_CHANNEL_ID = os.Getenv("WHITELISTED_CHANNEL_ID")

var gClient *http.Client
var gCalSrv *calendar.Service
var gSheetsSrv *sheets.Service

func init() {
	ctx := context.Background()
	gClient = getGoogleClient(ctx)
	gCalSrv = getCalendarService(ctx)
	gSheetsSrv = getSheetsService(ctx)
}

func main() {
	lambda.Start(handler)
}
