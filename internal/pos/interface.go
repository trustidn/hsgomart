package pos

type ServiceInterface interface {
	Checkout(tenantID, userID string, in CheckoutInput) (*CheckoutResult, error)
}

var _ ServiceInterface = (*Service)(nil)
