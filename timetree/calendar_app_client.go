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

type Attributes struct {
	Category      string `json:"category,omitempty"`
	Title         string `json:"title,omitempty"`
	AllDay        *bool  `json:"all_day,omitempty"`
	StartAt       string `json:"start_at,omitempty"`
	StartTimezone string `json:"start_timezone,omitempty"`
	EndAt         string `json:"end_at,omitempty"`
	EndTimezone   string `json:"end_timezone,omitempty"`
	Description   string `json:"description,omitempty"`
	Location      string `json:"location,omitempty"`
	URL           string `json:"url,omitempty"`
}

type CalendarEventRequestData struct {
	Attributes *Attributes `json:"attributes,omitempty"`
}

type CalendarEventRequest struct {
	Data *CalendarEventRequestData `json:"data,omitempty"`
}

type CalendarEventResponseData struct {
	ID         string      `json:"id,omitempty"`
	Type       string      `json:"type,omitempty"`
	Attributes *Attributes `json:"attributes,omitempty"`
}

type CalendarEventResponse struct {
	api.ErrorResponse
	Data *CalendarEventResponseData `json:"data,omitempty"`
}

// CreateCalendarEventRequest type
type CreateCalendarEventRequest CalendarEventRequest

// CreateCalendarEventResponse type
type CreateCalendarEventResponse CalendarEventResponse

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

// UpdateCalendarEventRequest type
type UpdateCalendarEventRequest CalendarEventRequest

// UpdateCalendarEventResponse type
type UpdateCalendarEventResponse CalendarEventResponse

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
