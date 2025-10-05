package domain

type GoogleFlightsResponse struct {
	BestFlights []struct {
		Flights []struct {
			DepartureAirport struct {
				Name string `json:"name"`
				ID   string `json:"id"`
				Time string `json:"time"`
			} `json:"departure_airport"`
			ArrivalAirport struct {
				Name string `json:"name"`
				ID   string `json:"id"`
				Time string `json:"time"`
			} `json:"arrival_airport"`
			Duration int    `json:"duration"`
			Airline  string `json:"airline"`
			Price    float64
		} `json:"flights"`
		Price float64 `json:"price"`
	} `json:"best_flights"`
	OtherFlights []struct {
		Flights []struct {
			DepartureAirport struct {
				ID   string `json:"id"`
				Time string `json:"time"`
			} `json:"departure_airport"`
			ArrivalAirport struct {
				ID   string `json:"id"`
				Time string `json:"time"`
			} `json:"arrival_airport"`
			Duration int    `json:"duration"`
			Airline  string `json:"airline"`
		} `json:"flights"`
		Price float64 `json:"price"`
	} `json:"other_flights"`
}
