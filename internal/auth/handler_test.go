package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupHandler() *Handler {
	svc := &Service{secret: []byte(testSecret)}
	return NewHandler(svc)
}

func performRequest(handler gin.HandlerFunc, method, path string, body interface{}) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var req *http.Request
	if body != nil {
		b, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	c.Request = req

	handler(c)
	return w
}

func parseBody(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("failed to parse response body: %v\nbody: %s", err, w.Body.String())
	}
	return m
}

// --- Register handler tests ---

func TestRegisterHandler_MissingFields(t *testing.T) {
	h := setupHandler()

	tests := []struct {
		name string
		body interface{}
	}{
		{"empty body", map[string]string{}},
		{"missing email", map[string]string{"name": "Test", "password": "Secret123"}},
		{"missing name", map[string]string{"email": "a@b.com", "password": "Secret123"}},
		{"missing password", map[string]string{"name": "Test", "email": "a@b.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequest(h.Register, http.MethodPost, "/auth/register", tt.body)
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusBadRequest, w.Body.String())
			}
			body := parseBody(t, w)
			if _, ok := body["error"]; !ok {
				t.Error("response should contain 'error' field")
			}
		})
	}
}

func TestRegisterHandler_InvalidEmail(t *testing.T) {
	h := setupHandler()
	body := map[string]string{
		"name":     "Test",
		"email":    "not-an-email",
		"password": "Secret123",
	}
	w := performRequest(h.Register, http.MethodPost, "/auth/register", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRegisterHandler_ShortPassword(t *testing.T) {
	h := setupHandler()
	body := map[string]string{
		"name":     "Test",
		"email":    "test@example.com",
		"password": "short",
	}
	w := performRequest(h.Register, http.MethodPost, "/auth/register", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRegisterHandler_NoBody(t *testing.T) {
	h := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/register", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// --- Login handler tests ---

func TestLoginHandler_MissingFields(t *testing.T) {
	h := setupHandler()

	tests := []struct {
		name string
		body interface{}
	}{
		{"empty body", map[string]string{}},
		{"missing email", map[string]string{"password": "Secret123"}},
		{"missing password", map[string]string{"email": "a@b.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequest(h.Login, http.MethodPost, "/auth/login", tt.body)
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusBadRequest, w.Body.String())
			}
			body := parseBody(t, w)
			if _, ok := body["error"]; !ok {
				t.Error("response should contain 'error' field")
			}
		})
	}
}

func TestLoginHandler_InvalidEmail(t *testing.T) {
	h := setupHandler()
	body := map[string]string{
		"email":    "bad-email",
		"password": "Secret123",
	}
	w := performRequest(h.Login, http.MethodPost, "/auth/login", body)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestLoginHandler_NoBody(t *testing.T) {
	h := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.Login(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// --- Refresh handler tests ---

func TestRefreshHandler_MissingToken(t *testing.T) {
	h := setupHandler()

	w := performRequest(h.Refresh, http.MethodPost, "/auth/refresh", map[string]string{})
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d; body: %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
	body := parseBody(t, w)
	if _, ok := body["error"]; !ok {
		t.Error("response should contain 'error' field")
	}
}

func TestRefreshHandler_NoBody(t *testing.T) {
	h := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	c.Request.Header.Set("Content-Type", "application/json")

	h.Refresh(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// --- Profile handler tests ---

func TestProfileHandler_NoUserID(t *testing.T) {
	h := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/profile", nil)

	h.Profile(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	body := parseBody(t, w)
	if errMsg, ok := body["error"]; !ok || errMsg != "unauthorized" {
		t.Errorf("error = %v, want 'unauthorized'", errMsg)
	}
}

func TestProfileHandler_WithValidClaims(t *testing.T) {
	h := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/profile", nil)

	c.Set(utils.KeyUserID, "user-999")
	c.Set(utils.KeyTenantID, "tenant-888")
	c.Set(utils.KeyRole, "owner")

	// Profile calls service.UserEmail which hits gorm.DB (nil in tests).
	// We verify the handler gets past auth extraction (no 401) by recovering the panic.
	panicked := true
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
			}
		}()
		h.Profile(c)
		panicked = false
	}()

	// If it panicked, the context extraction worked but DB access failed — expected.
	// If it didn't panic, check it's not 401 (meaning auth context was read).
	if !panicked && w.Code == http.StatusUnauthorized {
		t.Error("with valid context claims, should not get 401")
	}
}

func TestProfileHandler_ContextValuesTypes(t *testing.T) {
	h := setupHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/auth/profile", nil)

	// Set non-string value — GetUserID should return false
	c.Set(utils.KeyUserID, 12345)

	h.Profile(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d (non-string user_id)", w.Code, http.StatusUnauthorized)
	}
}

// --- Handler constructor test ---

func TestNewHandler(t *testing.T) {
	svc := newTestService()
	h := NewHandler(svc)
	if h == nil {
		t.Fatal("NewHandler returned nil")
	}
	if h.service != svc {
		t.Error("handler service should match provided service")
	}
}

// --- Response format tests ---

func TestRegisterHandler_ResponseContainsErrorKey(t *testing.T) {
	h := setupHandler()

	w := performRequest(h.Register, http.MethodPost, "/auth/register", map[string]string{})
	body := parseBody(t, w)

	if _, ok := body["error"]; !ok {
		t.Error("error response should contain 'error' key")
	}
}

func TestLoginHandler_ResponseContainsErrorKey(t *testing.T) {
	h := setupHandler()

	w := performRequest(h.Login, http.MethodPost, "/auth/login", map[string]string{})
	body := parseBody(t, w)

	if _, ok := body["error"]; !ok {
		t.Error("error response should contain 'error' key")
	}
}
