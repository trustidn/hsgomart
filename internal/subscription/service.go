package subscription

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrSubscriptionRequired = errors.New("subscription required")
	ErrSubscriptionExpired  = errors.New("subscription expired")
	ErrPlanLimitReached     = errors.New("plan limit reached")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CheckSubscription returns the active/trial subscription and plan for the tenant.
// Returns ErrSubscriptionRequired if no active/trial subscription, ErrSubscriptionExpired if end_date has passed.
func (s *Service) CheckSubscription(tenantID string) (*SubscriptionWithPlan, error) {
	subWithPlan, err := GetActiveSubscription(s.db, tenantID)
	if err != nil {
		return nil, err
	}
	if subWithPlan == nil {
		return nil, ErrSubscriptionRequired
	}
	// Trial and active both valid until end_date
	if subWithPlan.Subscription.EndDate != nil && subWithPlan.Subscription.EndDate.Before(time.Now()) {
		return nil, ErrSubscriptionExpired
	}
	return subWithPlan, nil
}
