package opname

type ServiceInterface interface {
	StartOpname(tenantID, userID string) (*StockOpname, error)
	SubmitItems(tenantID, opnameID string, items []SubmitItemInput) (*StockOpname, error)
	ApproveOpname(tenantID, opnameID string) (*StockOpname, error)
	GetOpname(tenantID, opnameID string) (*StockOpname, []StockOpnameItem, error)
	ListOpnames(tenantID string, limit, offset int) ([]StockOpname, error)
}

var _ ServiceInterface = (*Service)(nil)
