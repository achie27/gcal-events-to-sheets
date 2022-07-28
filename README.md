## gcal-to-sheets
I created this to find out about all the events with external attendees happening in my org.

### Pre-requisites
- Get creds, with `calendar.CalendarReadonlyScope` and `sheets.SpreadsheetsScope`, scopes using [this guide](https://developers.google.com/calendar/api/quickstart/go)
- Post deploying this lambda, you need to call the [watch events](https://developers.google.com/calendar/api/v3/reference/events/watch) API, with this lambda's function URL as part of the request and all the public calendars you want to watch, for Google to actually send webhooks.

### Pseudo-issues
- *Event duplicacy*: The webhooks from Google unfortunately don't contain ANY info about addition/removal/updation of events - just that something has changed in the events resource of a calendar (id). The logic in this repo gets the list of ([hopefully](https://developers.google.com/calendar/api/v3/reference/events/list#orderBy)) the most recent events and just writes them all to the spreadsheet. This means there can be duplicacy. While duplicates can be removed using the tools this repo already has, duplicates can also be removed by a 2-click journey on the spreadsheet itself manually, so I'm not handling it and indulging my procrastinating self.