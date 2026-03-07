package purchase

type ServiceInterface interface {
	CreatePurchase(tenantID string, in CreatePurchaseInput) (*Purchase, error)
	ListPurchases(tenantID string) ([]PurchaseListRow, error)
	GetPurchase(tenantID, purchaseID string) (*Purchase, []PurchaseItem, error)
	GetPurchaseWithItemRows(tenantID, purchaseID string) (*Purchase, []PurchaseItemRow, error)
}

var _ ServiceInterface = (*Service)(nil)
