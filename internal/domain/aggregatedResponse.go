package domain

type AggregatedResponse struct {
	Cheapest *Quote  `json:"cheapest"`
	Fastest  *Quote  `json:"fastest"`
	Offers   []Quote `json:"offers"` // order by price, after by duration
}
