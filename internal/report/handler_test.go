package report

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

func parseJSON(t *testing.T, w *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("unmarshal response body: %v\nbody: %s", err, w.Body.String())
	}
	return m
}

func TestHandlers_NoTenantContext(t *testing.T) {
	h := newNilHandler()

	tests := []struct {
		name    string
		handler gin.HandlerFunc
	}{
		{"SalesSummary", h.SalesSummary},
		{"SalesDaily", h.SalesDaily},
		{"SalesTransactions", h.SalesTransactions},
		{"SalesHourly", h.SalesHourly},
		{"PaymentsReport", h.PaymentsReport},
		{"TopProducts", h.TopProducts},
		{"InventorySummary", h.InventorySummary},
		{"ProfitReport", h.ProfitReport},
		{"CashiersReport", h.CashiersReport},
		{"ShiftsReport", h.ShiftsReport},
		{"SalesCompare", h.SalesCompare},
		{"ProductMargin", h.ProductMargin},
		{"GetReceipt", h.GetReceipt},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			tt.handler(c)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
			}
			body := parseJSON(t, w)
			if msg, _ := body["error"].(string); msg != "tenant context required" {
				t.Errorf("error = %q, want %q", msg, "tenant context required")
			}
		})
	}
}

func TestSalesCompare_MissingQueryParams(t *testing.T) {
	h := newDryRunHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("tenant_id", "test-tenant")

	h.SalesCompare(c)

	// Missing query params should not produce 400; the handler defaults to current month.
	// Without a live DB the downstream query returns 500, which is expected in unit tests.
	if w.Code == http.StatusBadRequest {
		t.Fatal("missing query params should not produce 400; handler should use defaults")
	}
	if w.Code == http.StatusUnauthorized {
		t.Fatal("status should not be 401 when tenant is set")
	}
}

func TestGetReceipt_MissingID(t *testing.T) {
	h := newDryRunHandler(t)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/receipt/", nil)
	c.Set("tenant_id", "test-tenant")

	h.GetReceipt(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d; body = %s", w.Code, http.StatusNotFound, w.Body.String())
	}
	body := parseJSON(t, w)
	if msg, _ := body["error"].(string); msg != "transaction not found" {
		t.Errorf("error = %q, want %q", msg, "transaction not found")
	}
}
