package user

import (
	"errors"

	"github.com/trustidn/hsmart-saas/internal/subscription"
	"github.com/trustidn/hsmart-saas/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailExists      = errors.New("email already registered in this tenant")
	ErrInvalidRole      = errors.New("role must be owner or cashier")
	ErrCannotDeleteSelf = errors.New("cannot delete your own account")
)

const bcryptCost = 10

// PlanLimitChecker is used to enforce plan max_users (e.g. subscription.Service).
type PlanLimitChecker interface {
	CheckSubscription(tenantID string) (*subscription.SubscriptionWithPlan, error)
}

type Service struct {
	db               *gorm.DB
	planLimitChecker PlanLimitChecker
}

func NewService(db *gorm.DB, planLimitChecker PlanLimitChecker) *Service {
	return &Service{db: db, planLimitChecker: planLimitChecker}
}

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
}

type UpdateUserInput struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
	Role     *string `json:"role"`
	Status   *string `json:"status"`
}

func validRole(role string) bool {
	return role == "owner" || role == "cashier"
}

func (s *Service) CreateUser(tenantID string, in CreateUserInput) (*User, error) {
	if !validRole(in.Role) {
		return nil, ErrInvalidRole
	}
	if err := utils.ValidatePasswordStrength(in.Password); err != nil {
		return nil, err
	}
	if s.planLimitChecker != nil {
		subWithPlan, err := s.planLimitChecker.CheckSubscription(tenantID)
		if err != nil {
			return nil, err
		}
		var count int64
		s.db.Model(&User{}).Where("tenant_id = ?", tenantID).Count(&count)
		if int(count) >= subWithPlan.Plan.MaxUsers {
			return nil, subscription.ErrPlanLimitReached
		}
	}

	var existing User
	if err := s.db.Where("tenant_id = ? AND email = ?", tenantID, in.Email).First(&existing).Error; err == nil {
		return nil, ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcryptCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		TenantID:     tenantID,
		Name:         in.Name,
		Email:        in.Email,
		PasswordHash: string(hash),
		Role:         in.Role,
		Status:       "active",
	}
	if err := CreateUser(s.db, u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) ListUsers(tenantID string) ([]User, error) {
	return FindUsersByTenant(s.db, tenantID)
}

func (s *Service) GetUser(tenantID, userID string) (*User, error) {
	u, err := FindUserByID(s.db, tenantID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return u, nil
}

func (s *Service) UpdateUser(tenantID, userID string, in UpdateUserInput) (*User, error) {
	updates := make(map[string]interface{})
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Email != nil {
		updates["email"] = *in.Email
	}
	if in.Role != nil {
		if !validRole(*in.Role) {
			return nil, ErrInvalidRole
		}
		updates["role"] = *in.Role
	}
	if in.Status != nil {
		updates["status"] = *in.Status
	}
	if in.Password != nil && *in.Password != "" {
		if err := utils.ValidatePasswordStrength(*in.Password); err != nil {
			return nil, err
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(*in.Password), bcryptCost)
		if err != nil {
			return nil, err
		}
		updates["password_hash"] = string(hash)
	}

	if len(updates) == 0 {
		return s.GetUser(tenantID, userID)
	}

	if err := UpdateUser(s.db, tenantID, userID, updates); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return s.GetUser(tenantID, userID)
}

func (s *Service) DeleteUser(tenantID, currentUserID, targetUserID string) error {
	if currentUserID == targetUserID {
		return ErrCannotDeleteSelf
	}
	if err := DeleteUser(s.db, tenantID, targetUserID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
