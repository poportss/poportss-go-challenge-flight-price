package domain

import "time"

type Quote struct {
	Provider    string        `json:"provider"`
	Airline     string        `json:"airline"`
	Price       float64       `json:"price"`
	Currency    string        `json:"currency"`
	Duration    time.Duration `json:"duration"`
	DepartureAt time.Time     `json:"departure_at"`
	ArrivalAt   time.Time     `json:"arrival_at"`
	Origin      string        `json:"origin"`
	Destination string        `json:"destination"`
}
