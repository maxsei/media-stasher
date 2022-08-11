package limiter

import (
  "sync"
  "time"

  "github.com/gin-gonic/gin"
)

func NewLimiter(limit int, freq time.Duration) *Limiter {
  var m sync.Mutex
  ret := Limiter{
    Limit:     limit,
    Freq:      freq,
    AddrStore: make(map[string]Time),
    m:         &m,
  }
  return &ret
}

func NewLimiter() *Limiter {
	return &Limiter{
		Ch: make(chan time.Time),
		Count: 0,
	}
}

type CompletedRequest struct {
	ID string
	Time time.Time
}

type Limiter struct {
	Ch chan CompletedRequest
	Count int
}

type MiddlewareCancelFirst struct {
  Limit     int
  Freq      time.Duration
  AddrStore map[string]*Limiter // mayble clean these up every so often
  m         sync.Mutex
}

func (m *MiddlewareCancelFirst) Handle(c *gin.Context) {
	// Get client ip address.
  ip := c.ClientIP()

	// Create a rate limiter linked to the IP address if one doesn't exists and
	// increment.
  m.m.Lock()
  rateLimiter, ok := m.AddrStore[ip]
  if !ok {
    m.AddrStore[ip] = NewLimiter()
  }
	m.AddrStore[ip].Count++
  m.m.Unlock()

	// On exit of this function we must decrement out rate limiter.
	defer func() {
		m.m.Lock()
		rateLimiter.Count--
		m.m.Unlock()
	}()

	// Create a rate limit request and perform the route.
	// TODO: may have to recover from panics here.
	req := CompletedRequest{time.Now(), uuid.NewString()}
	go func() {
		c.Next()
		limiter.Ch <- req
	}

  for {
		// Receive all finished rate limit requests.
		reqRecv := <-rateLimiter.Ch
		// If the request we received has taken longer than the given frequency and
		// we are passed the count then abort.
		m.m.Lock()
		count := limiter.Count 
		m.m.Unlock()
		if reqRecv.Time.Sub(req.Time) > m.Freq && count >= m.Limit {
			c.Abort()
		}
		// else check if we this is the same request and exit gracefully
		if req.ID == reqRecv.ID {
			return
		}
  }
}
