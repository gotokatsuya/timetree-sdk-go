package timetree

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gotokatsuya/timetree-sdk-go/timetree/api"
)

type CalendarAppClient struct {
	client *api.Client
}

func NewCalendarAppClient(accessToken string) (*CalendarAppClient, error) {
	c := &CalendarAppClient{}
	cli, err := api.NewClient(accessToken, http.DefaultClient)
	if err != nil {
		return nil, err
	}
	c.client = cli
	return c, nil
}

// EventAttributes type
type EventAttributes struct {
	Category      string `json:"category"`
	Title         string `json:"title"`
	AllDay        bool   `json:"all_day"`
	StartAt       string `json:"start_at"`
	StartTimezone string `json:"start_timezone,omitempty"`
	EndAt         string `json:"end_at"`
	EndTimezone   string `json:"end_timezone,omitempty"`
	Description   string `json:"description,omitempty"`
	Location      string `json:"location,omitempty"`
	URL           string `json:"url,omitempty"`
}

type CreateCalendarEventRequestData struct {
	Attributes EventAttributes `json:"attributes"`
}

// CreateCalendarEventRequest type
type CreateCalendarEventRequest struct {
	Data CreateCalendarEventRequestData `json:"data"`
}

// CreateCalendarEventResponse type
type CreateCalendarEventResponse struct {
	api.ErrorResponse
	Data struct {
		ID         string          `json:"id"`
		Type       string          `json:"type"`
		Attributes EventAttributes `json:"attributes"`
	} `json:"data"`
}

func (c *CalendarAppClient) CreateCalendarEvent(ctx context.Context, req *CreateCalendarEventRequest) (*CreateCalendarEventResponse, *http.Response, error) {
	path := "/calendar/events"
	httpReq, err := c.client.NewRequest(http.MethodPost, path, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(CreateCalendarEventResponse)
	httpResp, err := c.client.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

type UpdateCalendarEventRequestData struct {
	Attributes EventAttributes `json:"attributes"`
}

// UpdateCalendarEventRequest type
type UpdateCalendarEventRequest struct {
	Data UpdateCalendarEventRequestData `json:"data"`
}

// UpdateCalendarEventResponse type
type UpdateCalendarEventResponse struct {
	api.ErrorResponse
	Data struct {
		ID         string          `json:"id"`
		Type       string          `json:"type"`
		Attributes EventAttributes `json:"attributes"`
	} `json:"data"`
}

func (c *CalendarAppClient) UpdateCalendarEvent(ctx context.Context, eventID string, req *UpdateCalendarEventRequest) (*UpdateCalendarEventResponse, *http.Response, error) {
	path := fmt.Sprintf("/calendar/events/%s", eventID)
	httpReq, err := c.client.NewRequest(http.MethodPut, path, req)
	if err != nil {
		return nil, nil, err
	}
	resp := new(UpdateCalendarEventResponse)
	httpResp, err := c.client.Do(ctx, httpReq, resp)
	if err != nil {
		return nil, httpResp, err
	}
	return resp, httpResp, nil
}

func (c *CalendarAppClient) DeleteCalendarEvent(ctx context.Context, eventID string) (*http.Response, error) {
	path := fmt.Sprintf("/calendar/events/%s", eventID)
	httpReq, err := c.client.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	httpResp, err := c.client.Do(ctx, httpReq, nil)
	if err != nil {
		return nil, err
	}
	return httpResp, nil
}
