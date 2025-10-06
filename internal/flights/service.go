package flights

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/poportss/go-challenge-flight-price/internal/domain"
	"github.com/poportss/go-challenge-flight-price/internal/providers"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	mu        sync.RWMutex
	providers []providers.Provider
	timeout   time.Duration
	cache     Cache
}

func NewService(p []providers.Provider, timeout time.Duration, cache Cache) *Service {
	return &Service{providers: p, timeout: timeout, cache: cache}
}

// AddProvider dynamically adds a new provider to the service
func (s *Service) AddProvider(p providers.Provider) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.providers = append(s.providers, p)
	log.Printf("✓ Provider %s added dynamically", p.Name())
}

// RemoveProvider removes a provider by its name
func (s *Service) RemoveProvider(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, p := range s.providers {
		if p.Name() == name {
			s.providers = append(s.providers[:i], s.providers[i+1:]...)
			log.Printf("✓ Provider %s removed", name)
			return
		}
	}
}

// Search queries all active providers concurrently and aggregates the results
func (s *Service) Search(ctx context.Context, req domain.SearchRequest) (domain.AggregatedResponse, error) {
	cacheKey := fmt.Sprintf("%s|%s|%s|%s",
		req.Origin,
		req.Destination,
		req.StartDate.Format("2006-01-02"),
		req.EndDate.Format("2006-01-02"))

	// Try fetching from cache first
	if v, ok := s.cache.Get(cacheKey); ok {
		log.Printf("✓ Cache HIT for %s", cacheKey)
		return v.(domain.AggregatedResponse), nil
	}

	log.Printf("✗ Cache MISS for %s", cacheKey)

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	// Lock for safe concurrent access to the provider list
	s.mu.RLock()
	providersCopy := make([]providers.Provider, len(s.providers))
	copy(providersCopy, s.providers)
	s.mu.RUnlock()

	type result struct {
		provider string
		quotes   []domain.Quote
		err      error
	}
	resCh := make(chan result, len(providersCopy))

	eg, ctx := errgroup.WithContext(ctx)
	for _, p := range providersCopy {
		prov := p
		eg.Go(func() error {
			log.Printf("→ Fetching from %s...", prov.Name())
			qs, err := prov.Search(ctx, req.Origin, req.Destination, req.StartDate, req.EndDate)
			if err != nil {
				log.Printf("✗ Error from provider %s: %v", prov.Name(), err)
			} else {
				log.Printf("✓ Provider %s returned %d quotes", prov.Name(), len(qs))
			}
			resCh <- result{provider: prov.Name(), quotes: qs, err: err}
			return nil
		})
	}

	go func() {
		_ = eg.Wait()
		close(resCh)
	}()

	all := make([]domain.Quote, 0, 16)
	var atLeastOne bool
	for r := range resCh {
		if r.err == nil && len(r.quotes) > 0 {
			all = append(all, r.quotes...)
			atLeastOne = true
		}
	}

	if !atLeastOne {
		return domain.AggregatedResponse{}, errors.New("no providers returned valid quotes")
	}

	// Sort by price, then by duration
	sort.Slice(all, func(i, j int) bool {
		if all[i].Price == all[j].Price {
			return all[i].Duration < all[j].Duration
		}
		return all[i].Price < all[j].Price
	})

	// The cheapest offer is the first
	cheapest := all[0]

	// Find the fastest offer
	fastest := all[0]
	for _, q := range all {
		if q.Duration < fastest.Duration {
			fastest = q
		}
	}

	resp := domain.AggregatedResponse{
		Cheapest: &cheapest,
		Fastest:  &fastest,
		Offers:   all,
	}

	// Store the response in cache for 30 seconds
	s.cache.Set(cacheKey, resp, 30*time.Second)
	log.Printf("✓ Response cached: %s", cacheKey)

	return resp, nil
}
