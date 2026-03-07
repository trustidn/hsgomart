package product

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
// ListProducts
// ---------------------------------------------------------------------------

func TestListProducts_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products", nil)

	h := newTestHandler()
	h.ListProducts(c)

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
// CreateProduct
// ---------------------------------------------------------------------------

func TestCreateProduct_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products", nil)

	h := newTestHandler()
	h.CreateProduct(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCreateProduct_NoBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.CreateProduct(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateProduct_MissingName(t *testing.T) {
	payload := `{"sell_price": 100}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.CreateProduct(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "Name") {
		t.Fatalf("expected error about Name field, got %q", body["error"])
	}
}

func TestCreateProduct_MissingSellPrice(t *testing.T) {
	payload := `{"name": "Widget"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.CreateProduct(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "SellPrice") {
		t.Fatalf("expected error about SellPrice field, got %q", body["error"])
	}
}

// ---------------------------------------------------------------------------
// GetProduct
// ---------------------------------------------------------------------------

func TestGetProduct_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products/abc", nil)

	h := newTestHandler()
	h.GetProduct(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetProduct_MissingID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products/", nil)
	c.Set("tenant_id", "t-1")
	// Gin params are empty by default; Param("id") returns ""

	h := newTestHandler()
	h.GetProduct(c)

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
// UpdateProduct
// ---------------------------------------------------------------------------

func TestUpdateProduct_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/products/abc", nil)

	h := newTestHandler()
	h.UpdateProduct(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestUpdateProduct_MissingID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/products/", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.UpdateProduct(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateProduct_NoBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/products/abc", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h := newTestHandler()
	h.UpdateProduct(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// ---------------------------------------------------------------------------
// DeleteProduct
// ---------------------------------------------------------------------------

func TestDeleteProduct_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/products/abc", nil)

	h := newTestHandler()
	h.DeleteProduct(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestDeleteProduct_MissingID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/products/", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.DeleteProduct(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// ---------------------------------------------------------------------------
// AddBarcode
// ---------------------------------------------------------------------------

func TestAddBarcode_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products/abc/barcodes", nil)

	h := newTestHandler()
	h.AddBarcode(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestAddBarcode_MissingProductID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products//barcodes", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.AddBarcode(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "product id required" {
		t.Fatalf("expected 'product id required', got %q", body["error"])
	}
}

func TestAddBarcode_NoBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products/abc/barcodes", nil)
	c.Set("tenant_id", "t-1")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h := newTestHandler()
	h.AddBarcode(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestAddBarcode_MissingBarcodeField(t *testing.T) {
	payload := `{}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/products/abc/barcodes", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")
	c.Params = gin.Params{{Key: "id", Value: "abc"}}

	h := newTestHandler()
	h.AddBarcode(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "Barcode") {
		t.Fatalf("expected error about Barcode field, got %q", body["error"])
	}
}

// ---------------------------------------------------------------------------
// ListCategories
// ---------------------------------------------------------------------------

func TestListCategories_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/categories", nil)

	h := newTestHandler()
	h.ListCategories(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

// ---------------------------------------------------------------------------
// CreateCategory
// ---------------------------------------------------------------------------

func TestCreateCategory_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/categories", nil)

	h := newTestHandler()
	h.CreateCategory(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestCreateCategory_NoBody(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/categories", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.CreateCategory(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateCategory_MissingName(t *testing.T) {
	payload := `{}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(payload))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.CreateCategory(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if !strings.Contains(body["error"], "Name") {
		t.Fatalf("expected error about Name field, got %q", body["error"])
	}
}

// ---------------------------------------------------------------------------
// GetProductByBarcode
// ---------------------------------------------------------------------------

func TestGetProductByBarcode_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products/barcode/12345", nil)

	h := newTestHandler()
	h.GetProductByBarcode(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestGetProductByBarcode_MissingBarcode(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products/barcode/", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.GetProductByBarcode(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["error"] != "barcode required" {
		t.Fatalf("expected 'barcode required', got %q", body["error"])
	}
}

// ---------------------------------------------------------------------------
// DeleteBarcode
// ---------------------------------------------------------------------------

func TestDeleteBarcode_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/products/abc/barcodes/123", nil)

	h := newTestHandler()
	h.DeleteBarcode(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestDeleteBarcode_MissingParams(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/products//barcodes/", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.DeleteBarcode(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

// ---------------------------------------------------------------------------
// ListBarcodes
// ---------------------------------------------------------------------------

func TestListBarcodes_NoTenantContext(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products/abc/barcodes", nil)

	h := newTestHandler()
	h.ListBarcodes(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestListBarcodes_MissingProductID(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/products//barcodes", nil)
	c.Set("tenant_id", "t-1")

	h := newTestHandler()
	h.ListBarcodes(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}
