package model

import "time"

type User struct {
	ID       string
	Username string
	Password string
	Created  time.Time
	Updated  time.Time
}

type Audio struct {
	ID       string
	Name     string
	Format   string
	Location string
}

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
