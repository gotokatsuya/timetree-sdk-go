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

func NewCalendarAppClient(accessToken string, httpClient *http.Client) (*CalendarAppClient, error) {
	c := &CalendarAppClient{}
	cli, err := api.NewClient(accessToken, httpClient)
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

type CreateCalendarEventRequest CalendarEventRequest

type CreateCalendarEventResponse CalendarEventResponse

// CreateCalendarEvent 予定の作成
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

type UpdateCalendarEventRequest CalendarEventRequest

type UpdateCalendarEventResponse CalendarEventResponse

// UpdateCalendarEvent 予定の更新
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

// DeleteCalendarEvent 予定の削除
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
