package providers

import (
	"context"
	"time"

	"github.com/poportss/go-challenge-flight-price/internal/domain"
)

// Provider interface padronizada para todos os providers de voos
type Provider interface {
	Name() string
	Search(ctx context.Context, origin, destination string, startDate, endDate time.Time) ([]domain.Quote, error)
}
