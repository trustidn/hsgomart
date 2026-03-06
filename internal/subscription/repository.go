package subscription

import (
	"gorm.io/gorm"
)

// SubscriptionWithPlan holds subscription and its plan for enforcement checks.
type SubscriptionWithPlan struct {
	Subscription Subscription
	Plan         Plan
}

// GetActiveSubscription returns the active or trial subscription with plan for the tenant.
// status IN ('active', 'trial'). Returns nil if none found.
func GetActiveSubscription(db *gorm.DB, tenantID string) (*SubscriptionWithPlan, error) {
	var sub Subscription
	err := db.Where("tenant_id = ? AND status IN ?", tenantID, []string{"active", "trial"}).
		Order("end_date DESC NULLS LAST").
		First(&sub).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var plan Plan
	if err := db.Where("id = ?", sub.PlanID).First(&plan).Error; err != nil {
		return nil, err
	}

	return &SubscriptionWithPlan{Subscription: sub, Plan: plan}, nil
}
