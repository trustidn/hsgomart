package subscription

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func newNilHandler() *Handler {
	return NewHandler(NewService(nil))
}

func newDryRunHandler(t *testing.T) *Handler {
	t.Helper()
	db, err := gorm.Open(postgres.Open("host=localhost dbname=testdb sslmode=disable"), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
	})
	if err != nil {
		t.Fatalf("gorm.Open dry-run: %v", err)
	}
	return NewHandler(NewService(db))
}

func TestGetSubscription_NoTenant(t *testing.T) {
	h := newNilHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	h.GetSubscription(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if msg, _ := body["error"].(string); msg != "tenant context required" {
		t.Errorf("error = %q, want %q", msg, "tenant context required")
	}
}

func TestListPlans_OK(t *testing.T) {
	h := newDryRunHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	h.ListPlans(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusOK, w.Body.String())
	}
}

func TestChangePlan_NoTenant(t *testing.T) {
	h := newNilHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)

	h.ChangePlan(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
	var body map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &body)
	if msg, _ := body["error"].(string); msg != "tenant context required" {
		t.Errorf("error = %q, want %q", msg, "tenant context required")
	}
}

func TestChangePlan_NoBody(t *testing.T) {
	h := newNilHandler()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
	c.Set("tenant_id", "test-tenant")

	h.ChangePlan(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusBadRequest, w.Body.String())
	}
}
