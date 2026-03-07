package shift

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestHandler() *Handler {
	return NewHandler(nil)
}

// ---------------------------------------------------------------------------
// OpenShift
// ---------------------------------------------------------------------------

func TestOpenShift_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/shifts/open", nil)

	h := newTestHandler()
	h.OpenShift(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "tenant context required" {
		t.Fatalf("expected 'tenant context required', got %q", body["error"])
	}
}

func TestOpenShift_NoUserContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/shifts/open", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.OpenShift(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "user context required" {
		t.Fatalf("expected 'user context required', got %q", body["error"])
	}
}

func TestOpenShift_NoBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/shifts/open", nil)
	c.Set("tenant_id", "t-1")
	c.Set("user_id", "u-1")

	h := newTestHandler()
	h.OpenShift(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestOpenShift_MissingOpeningCash(t *testing.T) {
	payload := `{}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/shifts/open", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")
	c.Set("user_id", "u-1")

	h := newTestHandler()
	h.OpenShift(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "OpeningCash") {
		t.Fatalf("expected error about OpeningCash, got %q", body["error"])
	}
}

// ---------------------------------------------------------------------------
// CloseShift
// ---------------------------------------------------------------------------

func TestCloseShift_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/shifts/close", nil)

	h := newTestHandler()
	h.CloseShift(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "tenant context required" {
		t.Fatalf("expected 'tenant context required', got %q", body["error"])
	}
}

func TestCloseShift_NoUserContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/shifts/close", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.CloseShift(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "user context required" {
		t.Fatalf("expected 'user context required', got %q", body["error"])
	}
}

func TestCloseShift_NoBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/shifts/close", nil)
	c.Set("tenant_id", "t-1")
	c.Set("user_id", "u-1")

	h := newTestHandler()
	h.CloseShift(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// ---------------------------------------------------------------------------
// GetCurrentShift
// ---------------------------------------------------------------------------

func TestGetCurrentShift_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/shifts/current", nil)

	h := newTestHandler()
	h.GetCurrentShift(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "tenant context required" {
		t.Fatalf("expected 'tenant context required', got %q", body["error"])
	}
}

func TestGetCurrentShift_NoUserContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/shifts/current", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.GetCurrentShift(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "user context required" {
		t.Fatalf("expected 'user context required', got %q", body["error"])
	}
}

// ---------------------------------------------------------------------------
// ListShifts
// ---------------------------------------------------------------------------

func TestListShifts_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/shifts", nil)

	h := newTestHandler()
	h.ListShifts(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "tenant context required" {
		t.Fatalf("expected 'tenant context required', got %q", body["error"])
	}
}
