package product

type ServiceInterface interface {
	CreateCategory(tenantID, name string) (*Category, error)
	ListCategories(tenantID string) ([]Category, error)
	UpdateCategory(tenantID, categoryID, name string) (*Category, error)
	DeleteCategory(tenantID, categoryID string) error
	CreateProduct(tenantID string, in CreateProductInput) (*Product, error)
	ListProducts(tenantID string) ([]Product, error)
	FindProductByBarcode(tenantID, barcode string) (*Product, error)
	GetProduct(tenantID, productID string) (*Product, error)
	GetProductBarcodes(tenantID, productID string) ([]string, error)
	UpdateProduct(tenantID, productID string, in UpdateProductInput) (*Product, error)
	DeleteProduct(tenantID, productID string) error
	AddBarcode(tenantID, productID, barcode string) (*ProductBarcode, error)
	DeleteBarcode(tenantID, productID, barcode string) error
}

var _ ServiceInterface = (*Service)(nil)
