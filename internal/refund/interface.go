package refund

type ServiceInterface interface {
	CreateRefund(tenantID, userID string, in RefundInput) (*RefundResult, error)
	ListRefunds(tenantID string, limit, offset int) ([]Refund, error)
}

var _ ServiceInterface = (*Service)(nil)
