package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/poportss/go-challenge-flight-price/internal/domain"
	"github.com/poportss/go-challenge-flight-price/internal/util"
)

type Amadeus struct {
	client       *http.Client
	baseURL      string
	clientID     string
	clientSecret string
	token        string
	tokenExp     time.Time
}

func NewAmadeus(client *http.Client, token string) *Amadeus {
	return &Amadeus{
		client:  client,
		baseURL: util.EnvOr("AMADEUS_BASE_URL", "https://test.api.amadeus.com"),
		token:   token,
	}
}

func (a *Amadeus) Name() string { return "Amadeus" }

func (a *Amadeus) Search(ctx context.Context, origin, destination string, startDate, endDate time.Time) ([]domain.Quote, error) {
	url := fmt.Sprintf("%s/v2/shopping/flight-offers?originLocationCode=%s&destinationLocationCode=%s&departureDate=%s&returnDate=%s&adults=1&max=3",
		a.baseURL, origin, destination, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("amadeus: build request failed: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+a.token)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("amadeus: http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("amadeus: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var out domain.AmadeusResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("amadeus: json decode failed: %w | raw: %s", err, string(body))
	}

	quotes := make([]domain.Quote, 0, len(out.Data))
	for _, d := range out.Data {
		if len(d.Itineraries) == 0 || len(d.Itineraries[0].Segments) == 0 {
			continue
		}

		segments := d.Itineraries[0].Segments
		first := segments[0]
		last := segments[len(segments)-1]

		layout := "2006-01-02T15:04:05"
		dep, err1 := time.Parse(layout, first.Departure.At)
		arr, err2 := time.Parse(layout, last.Arrival.At)
		if err1 != nil || err2 != nil {
			continue
		}

		price, _ := strconv.ParseFloat(d.Price.GrandTotal, 64)

		quotes = append(quotes, domain.Quote{
			Provider:    a.Name(),
			Airline:     first.CarrierCode,
			Price:       price,
			Currency:    d.Price.Currency,
			Duration:    arr.Sub(dep),
			DepartureAt: dep,
			ArrivalAt:   arr,
			Origin:      origin,
			Destination: destination,
		})
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("amadeus: no valid quotes found for %s→%s", origin, destination)
	}

	return quotes, nil
}

func GetAmadeusAccessToken(ctx context.Context) (string, error) {
	client := util.NewHTTPClient(10 * time.Second)

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", os.Getenv("AMADEUS_CLIENT_ID"))
	form.Set("client_secret", os.Getenv("AMADEUS_CLIENT_SECRET"))

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://test.api.amadeus.com/v1/security/oauth2/token",
		strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("amadeus auth failed: %s", string(body))
	}

	var data struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		TokenType   string `json:"token_type"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	return data.AccessToken, nil
}
