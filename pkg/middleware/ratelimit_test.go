package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRateLimiter_Allow_UnderLimit(t *testing.T) {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		max:      5,
		window:   time.Minute,
	}

	for i := 0; i < 5; i++ {
		if !rl.Allow("192.168.1.1") {
			t.Fatalf("request %d should be allowed (limit is 5)", i+1)
		}
	}
}

func TestRateLimiter_Allow_ExceedsLimit(t *testing.T) {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		max:      3,
		window:   time.Minute,
	}

	for i := 0; i < 3; i++ {
		if !rl.Allow("10.0.0.1") {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}

	if rl.Allow("10.0.0.1") {
		t.Error("4th request should be rejected after limit of 3")
	}
}

func TestRateLimiter_Allow_DifferentIPs(t *testing.T) {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		max:      1,
		window:   time.Minute,
	}

	if !rl.Allow("ip-a") {
		t.Error("first request from ip-a should be allowed")
	}
	if !rl.Allow("ip-b") {
		t.Error("first request from ip-b should be allowed")
	}
	if rl.Allow("ip-a") {
		t.Error("second request from ip-a should be rejected")
	}
}

func TestRateLimiter_Allow_WindowReset(t *testing.T) {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		max:      1,
		window:   50 * time.Millisecond,
	}

	if !rl.Allow("x") {
		t.Fatal("first request should be allowed")
	}
	if rl.Allow("x") {
		t.Fatal("second request should be rejected")
	}

	time.Sleep(60 * time.Millisecond)

	if !rl.Allow("x") {
		t.Error("request after window reset should be allowed")
	}
}

func TestRateLimitMiddleware_AllowsRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		max:      10,
		window:   time.Minute,
	}

	router := gin.New()
	router.Use(RateLimit(rl))
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestRateLimitMiddleware_Blocks(t *testing.T) {
	gin.SetMode(gin.TestMode)

	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		max:      2,
		window:   time.Minute,
	}

	router := gin.New()
	router.Use(RateLimit(rl))
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("request %d: expected 200, got %d", i+1, w.Code)
		}
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected status 429, got %d", w.Code)
	}
}
