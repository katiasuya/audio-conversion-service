package model

import "time"

type User struct {
	ID       string
	Username string
	Password string
	Created  time.Time
	Updated  time.Time
}
