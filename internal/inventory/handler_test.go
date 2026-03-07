package inventory

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
// LowStock
// ---------------------------------------------------------------------------

func TestLowStock_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/inventory/low-stock", nil)

	h := newTestHandler()
	h.LowStock(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "tenant context required" {
		t.Fatalf("expected 'tenant context required', got %q", body["error"])
	}
}

// ---------------------------------------------------------------------------
// List (inventory)
// ---------------------------------------------------------------------------

func TestList_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/inventory", nil)

	h := newTestHandler()
	h.List(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// ---------------------------------------------------------------------------
// ListMovements
// ---------------------------------------------------------------------------

func TestListMovements_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/inventory/movements", nil)

	h := newTestHandler()
	h.ListMovements(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// ---------------------------------------------------------------------------
// Expiring
// ---------------------------------------------------------------------------

func TestExpiring_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/inventory/expiring", nil)

	h := newTestHandler()
	h.Expiring(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// ---------------------------------------------------------------------------
// GetStock
// ---------------------------------------------------------------------------

func TestGetStock_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products/abc/stock", nil)

	h := newTestHandler()
	h.GetStock(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetStock_MissingProductID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products//stock", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.GetStock(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "product id required" {
		t.Fatalf("expected 'product id required', got %q", body["error"])
	}
}

// ---------------------------------------------------------------------------
// AdjustStock
// ---------------------------------------------------------------------------

func TestAdjustStock_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products/abc/adjust-stock", nil)

	h := newTestHandler()
	h.AdjustStock(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAdjustStock_MissingProductID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products//adjust-stock", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.AdjustStock(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAdjustStock_NoBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products/abc/adjust-stock", nil)
	c.Set("tenant_id", "t-1")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h := newTestHandler()
	h.AdjustStock(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAdjustStock_MissingQuantity(t *testing.T) {
	payload := `{"type": "adjustment"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products/abc/adjust-stock", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h := newTestHandler()
	h.AdjustStock(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "Quantity") {
		t.Fatalf("expected error about Quantity field, got %q", body["error"])
	}
}

func TestAdjustStock_MissingType(t *testing.T) {
	payload := `{"quantity": -5}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products/abc/adjust-stock", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h := newTestHandler()
	h.AdjustStock(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "Type") {
		t.Fatalf("expected error about Type field, got %q", body["error"])
	}
}

func TestAdjustStock_InvalidJSON(t *testing.T) {
	payload := `{invalid`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products/abc/adjust-stock", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h := newTestHandler()
	h.AdjustStock(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
