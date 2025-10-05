package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/poportss/go-challenge-flight-price/internal/domain"
	"github.com/poportss/go-challenge-flight-price/internal/util"
)

// GoogleFlights provider
type GoogleFlights struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

func NewGoogleFlights(client *http.Client, apiKey string) *GoogleFlights {
	base := util.EnvOr("GOOGLE_FLIGHTS_BASE_URL", "https://serpapi.com/search.json")
	if !strings.HasPrefix(base, "http") {
		base = "https://" + base
	}
	return &GoogleFlights{
		client:  client,
		baseURL: base,
		apiKey:  apiKey,
	}
}

func (g *GoogleFlights) Name() string { return "GoogleFlights" }

// Search implements the Provider interface and fetches flight data from SerpAPI (Google Flights)
func (g *GoogleFlights) Search(ctx context.Context, origin, destination string, startDate, endDate time.Time) ([]domain.Quote, error) {
	url := fmt.Sprintf("%s?engine=google_flights&departure_id=%s&arrival_id=%s&outbound_date=%s&return_date=%s&currency=USD&hl=en&api_key=%s",
		g.baseURL, origin, destination, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"), g.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("googleflights: build request failed: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("googleflights: http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("googleflights: bad status %d: %s", resp.StatusCode, string(body))
	}

	var out domain.GoogleFlightsResponse

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("googleflights: decode failed: %w", err)
	}

	layout := "2006-01-02 15:04"
	quotes := make([]domain.Quote, 0, len(out.BestFlights)+len(out.OtherFlights))

	// process best flights
	for _, f := range out.BestFlights {
		if len(f.Flights) == 0 {
			continue
		}
		first := f.Flights[0]
		last := f.Flights[len(f.Flights)-1]

		dep, err1 := time.Parse(layout, first.DepartureAirport.Time)
		arr, err2 := time.Parse(layout, last.ArrivalAirport.Time)
		if err1 != nil || err2 != nil {
			continue
		}

		quotes = append(quotes, domain.Quote{
			Provider:    g.Name(),
			Airline:     first.Airline,
			Price:       f.Price,
			Currency:    "USD",
			DepartureAt: dep,
			ArrivalAt:   arr,
			Duration:    arr.Sub(dep),
			Origin:      first.DepartureAirport.ID,
			Destination: last.ArrivalAirport.ID,
		})
	}

	// process other flights (optional)
	for _, f := range out.OtherFlights {
		if len(f.Flights) == 0 {
			continue
		}
		first := f.Flights[0]
		last := f.Flights[len(f.Flights)-1]

		dep, err1 := time.Parse(layout, first.DepartureAirport.Time)
		arr, err2 := time.Parse(layout, last.ArrivalAirport.Time)
		if err1 != nil || err2 != nil {
			continue
		}

		quotes = append(quotes, domain.Quote{
			Provider:    g.Name(),
			Airline:     first.Airline,
			Price:       f.Price,
			Currency:    "USD",
			DepartureAt: dep,
			ArrivalAt:   arr,
			Duration:    arr.Sub(dep),
			Origin:      first.DepartureAirport.ID,
			Destination: last.ArrivalAirport.ID,
		})
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("googleflights: no valid quotes parsed")
	}

	return quotes, nil
}
