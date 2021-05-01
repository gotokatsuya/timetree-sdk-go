package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gotokatsuya/timetree-sdk-go/timetree"
)

var (
	installationID = os.Getenv("INSTALLATION_ID")
	calendarAppID  = os.Getenv("CALENDAR_APP_ID")
	privateKeyPath = os.Getenv("PRIVATE_KEY_PATH")
)

func main() {
	ctx := context.Background()

	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		panic(err)
	}

	authenticator, err := timetree.NewCalendarAppAuthenticator(calendarAppID, timetree.DefaultAccessTokenLifetime, privateKey, http.DefaultClient)
	if err != nil {
		panic(err)
	}
	accessTokenRes, httpRes, err := authenticator.AccessToken(ctx, installationID)
	if err != nil {
		panic(err)
	}
	switch httpRes.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		log.Printf("EventID: %s\n", accessTokenRes.AccessToken)
	default:
		log.Printf("Error: %s\n", accessTokenRes.ErrorResponse.Title)
		return
	}
	accessToken := accessTokenRes.AccessToken

	client, err := timetree.NewCalendarAppClient(http.DefaultClient)
	if err != nil {
		panic(err)
	}
	req := &timetree.CreateCalendarEventRequest{
		Data: &timetree.CalendarEventRequestData{
			Attributes: &timetree.Attributes{
				Category:      "schedule",
				Title:         "予定日1",
				AllDay:        timetree.Bool(true),
				StartAt:       "2020-11-18T00:00:00.000Z",
				StartTimezone: "Asia/Tokyo",
				EndAt:         "2020-11-23T00:00:00.000Z",
				EndTimezone:   "Asia/Tokyo",
				Description:   "予定日1の詳細",
			},
		},
	}
	calendarEventRes, httpRes, err := client.CreateCalendarEvent(ctx, accessToken, req)
	if err != nil {
		panic(err)
	}
	rateLimit := timetree.ParseRateLimit(httpRes)
	log.Printf("RateLimit: %v\n", rateLimit)

	switch httpRes.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent:
		log.Printf("EventID: %s\n", calendarEventRes.Data.ID)
	default:
		log.Printf("Error: %s\n", calendarEventRes.ErrorResponse.Title)
	}
}
