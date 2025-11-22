<div align="center">

# Token Bucket Rate Limiter

### Высокопроизводительный rate limiter на основе алгоритма Token Bucket для Go

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Tests](https://img.shields.io/badge/tests-passing-success?style=for-the-badge)](/)
[![Concurrency Safe](https://img.shields.io/badge/concurrency-safe-success?style=for-the-badge)](/)

</div>

---

## О проекте

**Token Bucket Rate Limiter** — это легковесная, потокобезопасная библиотека для ограничения частоты запросов (rate limiting) в Go приложениях.

### Зачем это нужно?

- 🛡️ **Защита API** от перегрузки и DDoS атак
- ⚡ **Контроль нагрузки** на внешние сервисы
- 💰 **Соблюдение лимитов** сторонних API (например, 100 запросов/сек)
- 🎮 **Fair use** ресурсов между пользователями

---

## ✨ Возможности

- ✅ **Простой API** — всего 3 метода: `New()`, `Allow()`, `Wait()`
- 🔒 **Concurrency-safe** — безопасен для использования из множества горутин
- ⚡ **Высокая производительность** — минимальные блокировки, эффективный refill
- 🕐 **Context-aware** — поддержка отмены через `context.Context`
- 🧪 **100% покрытие тестами** — включая concurrency тесты
- 📦 **Zero dependencies** — только стандартная библиотека Go
- 🪶 **Легковесный** — всего ~100 строк кода

---

## 📥 Установка

```bash
go get github.com/yourusername/token-bucket-limiter
```

**Требования**: Go 1.23.3+

---


## Тестирование

### Запуск тестов

```bash
# Все тесты
go test -v

# С race detector
go test -race

# С покрытием
go test -cover
```

### Тесты включают:

- ✅ **Базовая функциональность** (`TestAllowBasic`)
- ✅ **Refill механизм** (`TestRefill`)
- ✅ **Context отмена** (`TestWaitWithContextCancel`)
- ✅ **Concurrency** (`TestConcurrency`)

---

## Производительность

### Benchmarks

```bash
go test -bench=. -benchmem
```

### Оптимизации

- ✅ **Минимальные аллокации** (0 allocs/op для Allow)
- ✅ **Эффективные блокировки** (sync.Mutex только для критических секций)
- ✅ **Non-blocking refill** (через ticker в отдельной горутине)
- ✅ **Channel optimization** (non-blocking send в notify)

---

## 📖 Примеры использования

### Пример 1: HTTP API защита

```go
package main

import (
    "fmt"
    "net/http"
    limiter "github.com/yourusername/token-bucket-limiter"
)

func main() {
    rl := limiter.New(100, 10)
    defer rl.Stop()

    http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
        if !rl.Allow() {
            w.Header().Set("Retry-After", "1")
            http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
            return
        }
        fmt.Fprintf(w, "Hello, World!")
    })

    http.ListenAndServe(":8080", nil)
}
```

---

### Пример 2: Batch обработка

```go
package main

import (
    "context"
    "fmt"
    "time"
    limiter "github.com/yourusername/token-bucket-limiter"
)

func processBatch(items []string) {
    // Ограничение: 5 запросов/сек к внешнему API
    rl := limiter.New(5, 1)
    defer rl.Stop()

    for i, item := range items {
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        if err := rl.Wait(ctx); err != nil {
            fmt.Printf("Skipping item %d: %v\n", i, err)
            continue
        }

        // Вызов внешнего API
        callExternalAPI(item)
    }
}
```

---

### Пример 3: WebSocket rate limiting

```go
package main

import (
    "context"
    "time"
    "github.com/gorilla/websocket"
    limiter "github.com/yourusername/token-bucket-limiter"
)

type Client struct {
    conn    *websocket.Conn
    limiter *limiter.TokenBucket
}

func NewClient(conn *websocket.Conn) *Client {
    return &Client{
        conn:    conn,
        limiter: limiter.New(10, 5), // 10 msg/sec per client
    }
}

func (c *Client) ReadMessages() {
    defer c.limiter.Stop()

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }

        ctx, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        if err := c.limiter.Wait(ctx); err != nil {
            c.conn.WriteMessage(websocket.TextMessage, []byte("Rate limit exceeded"))
            continue
        }

        // Обработка сообщения
        processMessage(message)
    }
}
```

---

<div align="center">

**⭐ Если проект полезен, поставьте звезду!**

Made with ❤️ and Go

[⬆ Наверх](#-token-bucket-rate-limiter)

</div>

