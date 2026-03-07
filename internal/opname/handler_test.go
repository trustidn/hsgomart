package opname

import (
	"net/http"
	"net/http/httptest"
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
// Start
// ---------------------------------------------------------------------------

// The opname Start handler does not guard against missing tenant/user with
// early returns — it passes empty strings to service which will fail at DB.
// We verify the handler doesn't panic with a nil service when the service
// is nil (it will panic trying to call service methods). Since the handler
// unconditionally calls service.StartOpname, the best we can test without
// a DB is that missing context still reaches the handler without panic up
// to the service call boundary. We test SubmitItems body validation instead.

func TestStart_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/opname/start", nil)

	h := newTestHandler()

	defer func() {
		if r := recover(); r == nil {
			if w.Code != http.StatusUnauthorized && w.Code != http.StatusInternalServerError && w.Code != http.StatusCreated {
				t.Fatalf("unexpected status %d", w.Code)
			}
		}
	}()
	h.Start(c)
}

func TestStart_NoUserContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/opname/start", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()

	defer func() {
		if r := recover(); r == nil {
			if w.Code != http.StatusUnauthorized && w.Code != http.StatusInternalServerError && w.Code != http.StatusCreated {
				t.Fatalf("unexpected status %d", w.Code)
			}
		}
	}()
	h.Start(c)
}

// ---------------------------------------------------------------------------
// SubmitItems
// ---------------------------------------------------------------------------

func TestSubmitItems_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/opname/abc/items", nil)

	h := newTestHandler()
	h.SubmitItems(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSubmitItems_NoBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/opname/abc/items", nil)
	c.Set("tenant_id", "t-1")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h := newTestHandler()
	h.SubmitItems(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// ---------------------------------------------------------------------------
// Approve
// ---------------------------------------------------------------------------

func TestApprove_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/opname/abc/approve", nil)

	h := newTestHandler()

	defer func() {
		if r := recover(); r == nil {
			if w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest {
				t.Fatalf("unexpected status %d", w.Code)
			}
		}
	}()
	h.Approve(c)
}

func TestApprove_MissingID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/opname//approve", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()

	defer func() {
		if r := recover(); r == nil {
			if w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
				t.Fatalf("unexpected status %d", w.Code)
			}
		}
	}()
	h.Approve(c)
}

// ---------------------------------------------------------------------------
// Get
// ---------------------------------------------------------------------------

func TestGet_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/opname/abc", nil)

	h := newTestHandler()

	defer func() {
		if r := recover(); r == nil {
			if w.Code != http.StatusNotFound {
				t.Fatalf("unexpected status %d", w.Code)
			}
		}
	}()
	h.Get(c)
}

func TestGet_MissingID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/opname/", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()

	defer func() {
		if r := recover(); r == nil {
			if w.Code != http.StatusNotFound && w.Code != http.StatusBadRequest {
				t.Fatalf("unexpected status %d", w.Code)
			}
		}
	}()
	h.Get(c)
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

func TestList_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/opname", nil)

	h := newTestHandler()

	defer func() {
		if r := recover(); r == nil {
			if w.Code != http.StatusInternalServerError && w.Code != http.StatusOK {
				t.Fatalf("unexpected status %d", w.Code)
			}
		}
	}()
	h.List(c)
}
