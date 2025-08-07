package limiter

import (
	"sync"
	"sync/atomic"
	"context"
	"time"
	"testing"
)

func TestAllowBasic(t *testing.T) {
	limiter := New(10, 5)
	defer limiter.Stop()

	for i := 0; i < 5; i++ {
		if !limiter.Allow() {
			t.Errorf("Ожидалось true, получено false на итерации %d", i+1)
		}
	}
	if limiter.Allow() {
		t.Errorf("Ожидалось false, получено true")
	}
}

func TestRefill(t *testing.T) {
	limiter := New(10, 1)
	defer limiter.Stop()

	if !limiter.Allow() {
		t.Fatal("Не удалось получить единственный начальный токен")
	}

	if limiter.Allow() {
		t.Fatal("Лишний токен доступен сразу после исчерпания")
	}

	time.Sleep(101 * time.Millisecond)

	if !limiter.Allow() {
		t.Fatal("Новый токен не был добавлен")
	}
}

func TestWaitWithContextCancel(t *testing.T) {
	limiter := New(10, 1)
	defer limiter.Stop()
	limiter.Allow()

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	err := limiter.Wait(ctx)
	if err == nil {
		t.Errorf("Ожидалась ошибка отмены контекста, получено nil")
	}

	if err != context.Canceled {
		t.Errorf("Ожидалась ошибка отмены контекста, получена другая ошибка: %v", err)
	}
}

func TestConcurrency(t *testing.T) {
	limiter := New(100, 50)
	defer limiter.Stop()

	var wg sync.WaitGroup

	var successCount int64

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if limiter.Allow() {
				atomic.AddInt64(&successCount, 1)
			}
		}()
	}

	wg.Wait()

	if successCount != 50 {
		t.Errorf("Ожидалось 50 успешных запросов, но получено: %d", successCount)
	}
}