package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/trustidn/hsmart-saas/internal/subscription"
	"github.com/trustidn/hsmart-saas/internal/tenant"
	"github.com/trustidn/hsmart-saas/internal/user"
	"github.com/trustidn/hsmart-saas/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrEmailExists         = errors.New("email already registered")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

const (
	bcryptCost           = 10
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

type JWTClaims struct {
	UserID   string `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type RefreshToken struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    string    `gorm:"type:uuid;not null;index"`
	TokenHash string    `gorm:"column:token_hash;type:text;not null"`
	ExpiresAt time.Time `gorm:"column:expires_at;not null"`
	Revoked   bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

type TokenPair struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type Service struct {
	db     *gorm.DB
	secret []byte
}

func NewService(db *gorm.DB, jwtSecret string) *Service {
	return &Service{db: db, secret: []byte(jwtSecret)}
}

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func (s *Service) Register(in RegisterInput) (*TokenPair, error) {
	if err := utils.ValidatePasswordStrength(in.Password); err != nil {
		return nil, err
	}

	var existing user.User
	if err := s.db.Where("email = ?", in.Email).First(&existing).Error; err == nil {
		return nil, ErrEmailExists
	}

	t := tenant.Tenant{
		Name:   in.Name,
		Email:  in.Email,
		Status: "active",
	}
	if err := s.db.Create(&t).Error; err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcryptCost)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	plan, err := s.getOrCreateTrialPlan()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	end := now.AddDate(0, 1, 0)
	sub := subscription.Subscription{
		TenantID:  t.ID,
		PlanID:    plan.ID,
		Status:    "trial",
		StartDate: &now,
		EndDate:   &end,
	}
	if err := s.db.Create(&sub).Error; err != nil {
		return nil, err
	}

	return s.generateTokenPair(u.ID, t.ID, u.Role)
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

func (s *Service) Login(in LoginInput) (*TokenPair, error) {
	var u user.User
	if err := s.db.Where("email = ?", in.Email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokenPair(u.ID, u.TenantID, u.Role)
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *Service) Refresh(rawToken string) (*TokenPair, error) {
	hash := hashToken(rawToken)

	var rt RefreshToken
	err := s.db.Where("token_hash = ? AND revoked = false AND expires_at > ?", hash, time.Now()).First(&rt).Error
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	s.db.Model(&rt).Update("revoked", true)

	var u user.User
	if err := s.db.Where("id = ?", rt.UserID).First(&u).Error; err != nil {
		return nil, err
	}

	return s.generateTokenPair(u.ID, u.TenantID, u.Role)
}

func (s *Service) RevokeRefreshTokens(userID string) error {
	return s.db.Model(&RefreshToken{}).Where("user_id = ? AND revoked = false", userID).Update("revoked", true).Error
}

func (s *Service) generateAccessToken(userID, tenantID, role string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		TenantID: tenantID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.secret)
}

func (s *Service) generateTokenPair(userID, tenantID, role string) (*TokenPair, error) {
	accessToken, err := s.generateAccessToken(userID, tenantID, role)
	if err != nil {
		return nil, err
	}

	rawRefresh, err := generateRandomToken(32)
	if err != nil {
		return nil, err
	}

	rt := RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(rawRefresh),
		ExpiresAt: time.Now().Add(refreshTokenDuration),
	}
	if err := s.db.Create(&rt).Error; err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefresh,
	}, nil
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

func generateRandomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}
