package main

import (
	"context"
	"log"
	"os"

	"github.com/gotokatsuya/timetree-sdk-go/timetree"
)

var (
	InstallationID = os.Getenv("INSTALLATION_ID")
	CalendarAppID  = os.Getenv("CALENDAR_APP_ID")
	PrivateKeyPath = os.Getenv("PRIVATE_KEY_PATH")
)

func main() {
	ctx := context.Background()

	authenticator, err := timetree.NewCalendarAppAuthenticator(CalendarAppID, PrivateKeyPath)
	if err != nil {
		panic(err)
	}
	accessTokenRes, httpRes, err := authenticator.AccessToken(ctx, InstallationID)
	if err != nil {
		panic(err)
	}
	switch httpRes.StatusCode {
	case 200, 201, 204:
		log.Printf("EventID: %s\n", accessTokenRes.AccessToken)
	default:
		log.Printf("Error: %s\n", accessTokenRes.ErrorResponse.Title)
		return
	}

	client, err := timetree.NewCalendarAppClient(accessTokenRes.AccessToken)
	if err != nil {
		panic(err)
	}
	req := &timetree.CreateCalendarEventRequest{
		Data: timetree.CreateCalendarEventRequestData{
			Attributes: timetree.EventAttributes{
				Category:      "schedule",
				Title:         "予定日",
				AllDay:        true,
				StartAt:       "2020-11-18T00:00:00.000Z",
				StartTimezone: "Asia/Tokyo",
				EndAt:         "2020-11-23T00:00:00.000Z",
				EndTimezone:   "Asia/Tokyo",
			},
		},
	}
	calendarEventRes, httpRes, err := client.CreateCalendarEvent(ctx, req)
	if err != nil {
		panic(err)
	}
	switch httpRes.StatusCode {
	case 200, 201, 204:
		log.Printf("EventID: %s\n", calendarEventRes.Data.ID)
	default:
		log.Printf("Error: %s\n", calendarEventRes.ErrorResponse.Title)
	}
}
