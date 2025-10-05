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

type SearchRequest struct {
	Origin      string    `form:"origin" binding:"required,len=3"`
	Destination string    `form:"destination" binding:"required,len=3"`
	StartDate   time.Time `form:"starDate" time_format:"2006-01-02" binding:"required"`
	EndDate     time.Time `form:"endDate" time_format:"2006-01-02" binding:"required"`
}

type AggregatedResponse struct {
	Cheapest *Quote  `json:"cheapest"`
	Fastest  *Quote  `json:"fastest"`
	Offers   []Quote `json:"offers"` // order by price, after by duration
}
