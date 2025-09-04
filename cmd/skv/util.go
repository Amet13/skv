package main

import (
	"context"
	crand "crypto/rand"
	"errors"
	"math/big"
	"time"

	"skv/internal/provider"
)

func fetchWithRetry(ctx context.Context, p provider.Provider, spec provider.SecretSpec, retries int, delay time.Duration) (string, error) {
	val, err := p.FetchSecret(ctx, spec)
	if err == nil || retries <= 0 {
		return val, err
	}
	attempts := retries
	backoff := delay
	for attempts > 0 {
		// Do not retry not found
		if errors.Is(err, provider.ErrNotFound) {
			return "", err
		}
		// jitter +/-20%
		sleep := backoff
		maxJitter := backoff / 5
		if maxJitter > 0 {
			if n, nErr := crand.Int(crand.Reader, big.NewInt(int64(maxJitter))); nErr == nil {
				j := time.Duration(n.Int64())
				sleep = backoff - j/2
			}
		}
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(sleep):
		}
		val, err = p.FetchSecret(ctx, spec)
		if err == nil {
			return val, nil
		}
		attempts--
		if backoff < 10*time.Second {
			backoff = backoff * 2
		}
	}
	return "", err
}

