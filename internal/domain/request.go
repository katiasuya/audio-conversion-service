package domain

// Request represents audio conversion request.
type Request struct {
	ID             string `json:"ID"`
	OriginalID     string `json:"originalID"`
	OriginalFormat string `json:"originalFormat"`
	TargetID       string `json:"targetID"`
	TargetFormat   string `json:"targetFormat"`
	Created        string `json:"created"`
	Updated        string `json:"updated"`
	Status         string `json:"status"`
}
