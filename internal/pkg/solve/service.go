package solve

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"strconv"

	"pow-f/internal/pkg/back"
)

type Service struct {
	token   string
	ip      string
	target  *big.Int
	backend backend
}

type backend interface {
	Auth(ctx context.Context) (*back.AuthResponse, error)
	Solve(ctx context.Context, r *back.Request) (*back.SolveResponse, bool, error)
}

type result struct {
	hash  string
	nonce int
}

var maxNonce = math.MaxInt64

func New(b backend, t int) *Service {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-t))
	return &Service{
		backend: b,
		target:  target,
	}
}

func (s *Service) auth(ctx context.Context) error {
	res, err := s.backend.Auth(ctx)
	if err != nil {
		return fmt.Errorf("couldn't autorize in backend: %w", err)
	}
	s.token = res.Token
	s.ip = res.Ip

	return nil
}

func (s *Service) Solve(ctx context.Context) (string, error) {
	if s.token == "" {
		err := s.auth(ctx)
		if err != nil {
			return "", err
		}
	}

	result := s.calculate(ctx)
	req := &back.Request{
		Token: s.token,
		Ip:    s.ip,
		Hash:  result.hash,
		Nonce: result.nonce,
	}

	res, isUnauthorized, err := s.backend.Solve(ctx, req)
	if isUnauthorized {
		s.backend.Auth(ctx)
		res, isUnauthorized, err = s.backend.Solve(ctx, req)
	}
	if err != nil {
		return "", fmt.Errorf("couldn't solve: %w", err)
	}

	return res.Phrase, nil
}

func (s *Service) calculate(_ context.Context) *result {
	var nonce int
	var hashInt big.Int
	var hash [32]byte
	res := &result{}

	for nonce < maxNonce {
		if nonce%100 == 0 {
			log.Println("done %d nonce", nonce)
		}
		data := s.prepare(nonce)

		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(s.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	res.hash = hex.EncodeToString(hash[:])
	res.nonce = nonce

	return res
}

func (s *Service) prepare(nonce int) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(s.token),
			[]byte(s.ip),
			[]byte(strconv.FormatInt(int64(nonce), 16)),
		},
		[]byte{},
	)
}
