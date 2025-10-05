package providers

import (
	"context"
	"math/rand"
	"time"

	"github.com/poportss/go-challenge-flight-price/internal/domain"
)

type MockProvider struct {
	name string
}

func NewMockProvider(name string) *MockProvider {
	return &MockProvider{name: name}
}

func (m *MockProvider) Name() string { return m.name }

func (m *MockProvider) Search(ctx context.Context, origin, destination string, startDate, endDate time.Time) ([]domain.Quote, error) {
	// simula tempo de resposta variável (200–600ms)
	time.Sleep(time.Duration(200+rand.Intn(400)) * time.Millisecond)

	// simula cotações falsas
	price := 500 + float64(rand.Intn(400))
	duration := time.Duration(6+rand.Intn(4)) * time.Hour // 6–10h

	departure := startDate.Add(10 * time.Hour)
	arrival := departure.Add(duration)

	return []domain.Quote{
		{
			Provider:    m.Name(),
			Airline:     "MockAir",
			Price:       price,
			Currency:    "USD",
			DepartureAt: departure,
			ArrivalAt:   arrival,
			Duration:    duration,
			Origin:      origin,
			Destination: destination,
		},
	}, nil
}
