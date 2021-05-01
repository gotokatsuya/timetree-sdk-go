package timetree

type ErrorResponse struct {
	Type   string `json:"type,omitempty"`
	Status int    `json:"status,omitempty"`
	Title  string `json:"title,omitempty"`
	// Errors string `json:"errors"`
}
