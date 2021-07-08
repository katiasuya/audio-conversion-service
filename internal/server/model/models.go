package model

import "time"

// AudioInfo represents downloaded audio information.
type AudioInfo struct {
	ID       string `json:"ID"`
	Name     string `json:"name"`
	Format   string `json:"format"`
	Location string `json:"location"`
}

// RequestInfo represents a history response.
type RequestInfo struct {
	ID           string    `json:"ID"`
	AudioName    string    `json:"audioName"`
	SourceFormat string    `json:"sourceFormat"`
	TargetFormat string    `json:"targetFormat"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	Status       string    `json:"status"`
}
