package revenuecat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://api.revenuecat.com/v1"

type Service struct {
	apiKey     string
	httpClient *http.Client
}

func New(apiKey string) *Service {
	return &Service{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

type subscriberResponse struct {
	Subscriber struct {
		Entitlements map[string]struct {
			ExpiresDate *string `json:"expires_date"`
		} `json:"entitlements"`
	} `json:"subscriber"`
}

// getSubscriber fetches the subscriber from RevenueCat. Returns nil, nil if the subscriber does not exist.
func (s *Service) getSubscriber(ctx context.Context, userID string) (*subscriberResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/subscribers/%s", baseURL, userID), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("revenuecat returned status %d", resp.StatusCode)
	}

	var body subscriberResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	return &body, nil
}

// CustomerExists returns true if the subscriber exists in RevenueCat.
func (s *Service) CustomerExists(ctx context.Context, userID string) (bool, error) {
	sub, err := s.getSubscriber(ctx, userID)
	return sub != nil, err
}

// HasActiveSubscription returns true if the user has at least one active entitlement.
func (s *Service) HasActiveSubscription(ctx context.Context, userID string) (bool, error) {
	sub, err := s.getSubscriber(ctx, userID)
	if err != nil || sub == nil {
		return false, err
	}
	return len(sub.Subscriber.Entitlements) > 0, nil
}
