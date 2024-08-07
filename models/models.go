package models

import "time"

type TimedItem struct {
	Borntime time.Time `json:"born_time"`
	Lifetime int       `json:"life_time"`
}
