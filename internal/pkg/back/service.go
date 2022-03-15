package back

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	requestChallenge = "/request-challenge"
	solveChallenge   = "/solve-challenge"
)

type AuthResponse struct {
	Token string `json:"token"`
	Ip    string `json:"ip"`
}

type SolveResponse struct {
	Phrase string `json:"phrase"`
}

type Request struct {
	Token string `json:"token"`
	Ip    string `json:"ip"`
	Hash  string `json:"hash"`
	Nonce int    `json:"nonce"`
}

type Service struct {
	url    string
	client *http.Client
}

func New(url string, timeout time.Duration) *Service {
	return &Service{
		url: url,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *Service) Auth(ctx context.Context) (*AuthResponse, error) {
	res := &AuthResponse{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.url+requestChallenge, nil)
	if err != nil {
		return res, fmt.Errorf("couldn't create request challenge: %w", err)
	}
	r, err := s.client.Do(req)
	if err != nil {
		return res, fmt.Errorf("couldn't get response for request challenge: %w", err)
	}

	err = json.NewDecoder(r.Body).Decode(res)
	if err != nil {
		return res, fmt.Errorf("couldn't decode response for request challenge: %w", err)
	}

	return res, nil
}

func (s *Service) Solve(ctx context.Context, r *Request) (*SolveResponse, bool, error) {
	res := &SolveResponse{}
	body, err := json.Marshal(r)
	if err != nil {
		return res, false, fmt.Errorf("couldn't marshal request for solve challenge: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url+solveChallenge, bytes.NewBuffer(body))
	if err != nil {
		return res, false, fmt.Errorf("couldn't create request solve challenge: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	result, err := s.client.Do(req)
	if err != nil {
		if result != nil && result.StatusCode == http.StatusUnauthorized {
			return res, true, nil
		}
		return res, false, fmt.Errorf("couldn't get response for request challenge: %w", err)
	}

	err = json.NewDecoder(result.Body).Decode(res)
	if err != nil {
		return res, false, fmt.Errorf("couldn't decode response for solve challenge: %w", err)
	}

	return res, false, nil
}
