package limiter

import (
	"context"
	"sync"
	"time"
)

type TokenBucket struct {
	rate     int64
	capacity int64
	tokens   int64
	mu       sync.Mutex
	ticker   *time.Ticker
	notifyCh chan struct{}
	stopCh   chan struct{}
}

func New(rate, capacity int) *TokenBucket {
	if rate <= 0 || capacity <= 0 {
		return nil
	}

	tb := &TokenBucket{
		rate:     int64(rate),
		capacity: int64(capacity),
		tokens:   int64(capacity),
		notifyCh: make(chan struct{}),
		stopCh:   make(chan struct{}),
	}

	interval := time.Second / time.Duration(rate)
	tb.ticker = time.NewTicker(interval)

	go tb.refill()

	return tb
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func (tb *TokenBucket) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	for {
		tb.mu.Lock()
		if tb.tokens > 0 {
			tb.tokens--
			tb.mu.Unlock()
			return nil
		}
		tb.mu.Unlock()

		select {
		case <-tb.notifyCh:
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (tb *TokenBucket) Stop() {
	tb.ticker.Stop()
	close(tb.stopCh)
}

func (tb *TokenBucket) refill() {
	for {
		select {
		case <-tb.ticker.C:
			tb.mu.Lock()
			if tb.tokens < tb.capacity {
				tb.tokens++
				select {
				case tb.notifyCh <- struct{}{}:
					//
				default:
					//
				}
			}
			tb.mu.Unlock()
		case <-tb.stopCh:
			return
		}
	}
}
