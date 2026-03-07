package subscription

type ServiceInterface interface {
	CheckSubscription(tenantID string) (*SubscriptionWithPlan, error)
}

var _ ServiceInterface = (*Service)(nil)
