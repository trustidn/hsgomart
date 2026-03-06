package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/trustidn/hsmart-saas/internal/subscription"
	"github.com/trustidn/hsmart-saas/internal/tenant"
	"github.com/trustidn/hsmart-saas/internal/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already registered")
)

const bcryptCost = 10

type JWTClaims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Service struct {
	db   *gorm.DB
	secret []byte
}

func NewService(db *gorm.DB, jwtSecret string) *Service {
	return &Service{db: db, secret: []byte(jwtSecret)}
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (s *Service) Register(in RegisterInput) (string, error) {
	var existing user.User
	if err := s.db.Where("email = ?", in.Email).First(&existing).Error; err == nil {
		return "", ErrEmailExists
	}

	t := tenant.Tenant{
		Name:   in.Name,
		Email:  in.Email,
		Status: "active",
	}
	if err := s.db.Create(&t).Error; err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcryptCost)
	if err != nil {
		return "", err
	}

	u := user.User{
		TenantID:     t.ID,
		Name:         in.Name,
		Email:        in.Email,
		PasswordHash: string(hash),
		Role:         "owner",
		Status:       "active",
	}
	if err := s.db.Create(&u).Error; err != nil {
		return "", err
	}

	plan, err := s.getOrCreateTrialPlan()
	if err != nil {
		return "", err
	}

	now := time.Now()
	end := now.AddDate(0, 1, 0) // 30 days trial
	sub := subscription.Subscription{
		TenantID:  t.ID,
		PlanID:    plan.ID,
		Status:    "trial",
		StartDate: &now,
		EndDate:   &end,
	}
	if err := s.db.Create(&sub).Error; err != nil {
		return "", err
	}

	return s.generateToken(u.ID, t.ID, u.Role)
}

func (s *Service) getOrCreateTrialPlan() (subscription.Plan, error) {
	var plan subscription.Plan
	err := s.db.Where("name = ?", "Trial").First(&plan).Error
	if err == nil {
		return plan, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return subscription.Plan{}, err
	}
	plan = subscription.Plan{
		Name:        "Trial",
		Price:       0,
		MaxUsers:    5,
		MaxProducts: 100,
	}
	if err := s.db.Create(&plan).Error; err != nil {
		return subscription.Plan{}, err
	}
	return plan, nil
}

func (s *Service) Login(in LoginInput) (string, error) {
	var u user.User
	if err := s.db.Where("email = ?", in.Email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.generateToken(u.ID, u.TenantID, u.Role)
}

func (s *Service) generateToken(userID, tenantID, role string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secret)
}

func (s *Service) ValidateToken(tokenString string) (*JWTClaims, error) {
	t, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(*JWTClaims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (s *Service) UserEmail(userID string) (string, error) {
	var u user.User
	if err := s.db.Where("id = ?", userID).First(&u).Error; err != nil {
		return "", err
	}
	return u.Email, nil
}
