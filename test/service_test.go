// test/service_test.go
package test

import (
	"context"
	"testing"
	"time"

	"github.com/poportss/go-challenge-flight-price/internal/domain"
	"github.com/poportss/go-challenge-flight-price/internal/flights"
	"github.com/poportss/go-challenge-flight-price/internal/providers"
)

type fakeProv struct {
	name string
	qs   []domain.Quote
	err  error
}

func (f fakeProv) Name() string { return f.name }
func (f fakeProv) Search(ctx context.Context, o, d string, dt, et time.Time) ([]domain.Quote, error) {
	return f.qs, f.err
}

func TestAggregationCheapestFastest(t *testing.T) {
	now := time.Now()
	p1 := fakeProv{
		name: "amadeus",
		qs: []domain.Quote{
			{
				Provider:    "amadeus",
				Price:       100,
				Duration:    3 * time.Hour,
				DepartureAt: now,
				ArrivalAt:   now.Add(3 * time.Hour),
			},
		},
		err: nil,
	}

	p2 := fakeProv{
		name: "amadeus",
		qs: []domain.Quote{
			{
				Provider:    "Ports airlines",
				Price:       120,
				Duration:    2 * time.Hour,
				DepartureAt: now,
				ArrivalAt:   now.Add(2 * time.Hour),
			},
		},
		err: nil,
	}

	svc := flights.NewService([]providers.Provider{p1, p2}, 10*time.Second, flights.NewInMemoryTTL())

	requestData := domain.SearchRequest{
		Origin:      "GRU",
		Destination: "JFK",
		StartDate:   now,
		EndDate:     now.Add(86 * time.Hour),
	}

	resp, err := svc.Search(context.Background(), requestData)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Cheapest == nil {
		t.Fatalf("expected cheapest not nil")
	}

	if resp.Fastest == nil {
		t.Fatalf("expected fastest not nil")
	}

	if len(resp.Offers) < 1 {
		t.Fatalf("offers expected 1 or more, got %d", len(resp.Offers))
	}
}
