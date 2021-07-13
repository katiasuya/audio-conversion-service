package model

import "time"

// RequestInfo represents a history response.
type Request struct {
	ID           string
	UserID       string
	AudioName    string
	SourceID     string
	SourceFormat string
	TargetID     string
	TargetFormat string
	Created      time.Time
	Updated      time.Time
	Status       string
}
