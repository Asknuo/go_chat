package models

import "time"

type Group struct {
	ID        uint
	Name      string
	Avatar    string
	CreatedAt time.Time
}
