package user

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
// Create
// ---------------------------------------------------------------------------

func TestCreate_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", nil)

	h := newTestHandler()
	h.Create(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "tenant context required" {
		t.Fatalf("expected 'tenant context required', got %q", body["error"])
	}
}

func TestCreate_NoBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreate_MissingEmail(t *testing.T) {
	payload := `{"name":"Bob","password":"secret123","role":"cashier"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "Email") {
		t.Fatalf("expected error about Email field, got %q", body["error"])
	}
}

func TestCreate_MissingPassword(t *testing.T) {
	payload := `{"name":"Bob","email":"bob@example.com","role":"cashier"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "Password") {
		t.Fatalf("expected error about Password field, got %q", body["error"])
	}
}

func TestCreate_MissingName(t *testing.T) {
	payload := `{"email":"bob@example.com","password":"secret123","role":"cashier"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "Name") {
		t.Fatalf("expected error about Name field, got %q", body["error"])
	}
}

func TestCreate_MissingRole(t *testing.T) {
	payload := `{"name":"Bob","email":"bob@example.com","password":"secret123"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "Role") {
		t.Fatalf("expected error about Role field, got %q", body["error"])
	}
}

func TestCreate_InvalidEmail(t *testing.T) {
	payload := `{"name":"Bob","email":"not-an-email","password":"secret123","role":"cashier"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "Email") {
		t.Fatalf("expected error about Email field, got %q", body["error"])
	}
}

func TestCreate_PasswordTooShort(t *testing.T) {
	payload := `{"name":"Bob","email":"bob@example.com","password":"ab","role":"cashier"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.Create(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "Password") {
		t.Fatalf("expected error about Password field, got %q", body["error"])
	}
}

// ---------------------------------------------------------------------------
// Update
// ---------------------------------------------------------------------------

func TestUpdate_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/users/abc", nil)

	h := newTestHandler()
	h.Update(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUpdate_MissingID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/users/", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.Update(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "user id required" {
		t.Fatalf("expected 'user id required', got %q", body["error"])
	}
}

func TestUpdate_NoBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/users/abc", nil)
	c.Set("tenant_id", "t-1")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h := newTestHandler()
	h.Update(c)

	// UpdateUserInput has all optional fields, so ShouldBindJSON with nil body => 400 (EOF)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// ---------------------------------------------------------------------------
// Delete
// ---------------------------------------------------------------------------

func TestDelete_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/users/abc", nil)

	h := newTestHandler()
	h.Delete(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestDelete_MissingID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/users/", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.Delete(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "user id required" {
		t.Fatalf("expected 'user id required', got %q", body["error"])
	}
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

func TestList_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/users", nil)

	h := newTestHandler()
	h.List(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
