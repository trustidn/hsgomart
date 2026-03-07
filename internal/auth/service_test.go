package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const testSecret = "test-secret-key-for-unit-tests"

func newTestService() *Service {
	return &Service{secret: []byte(testSecret)}
}

func TestGenerateAccessToken(t *testing.T) {
	svc := newTestService()
	userID := "user-123"
	tenantID := "tenant-456"
	role := "owner"

	token, err := svc.generateAccessToken(userID, tenantID, role)
	if err != nil {
		t.Fatalf("generateAccessToken returned error: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestValidateToken_Valid(t *testing.T) {
	svc := newTestService()
	userID := "user-abc"
	tenantID := "tenant-xyz"
	role := "admin"

	token, err := svc.generateAccessToken(userID, tenantID, role)
	if err != nil {
		t.Fatalf("generateAccessToken returned error: %v", err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken returned error: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("UserID = %q, want %q", claims.UserID, userID)
	}
	if claims.TenantID != tenantID {
		t.Errorf("TenantID = %q, want %q", claims.TenantID, tenantID)
	}
	if claims.Role != role {
		t.Errorf("Role = %q, want %q", claims.Role, role)
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	svc := newTestService()

	claims := JWTClaims{
		UserID:   "user-1",
		TenantID: "tenant-1",
		Role:     "owner",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString(svc.secret)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = svc.ValidateToken(signed)
	if err == nil {
		t.Fatal("expected error for expired token, got nil")
	}
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	svc := newTestService()

	claims := JWTClaims{
		UserID:   "user-1",
		TenantID: "tenant-1",
		Role:     "owner",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tok.SignedString([]byte("wrong-secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	_, err = svc.ValidateToken(signed)
	if err == nil {
		t.Fatal("expected error for invalid signature, got nil")
	}
}

func TestValidateToken_MalformedToken(t *testing.T) {
	svc := newTestService()

	_, err := svc.ValidateToken("not.a.valid.jwt")
	if err == nil {
		t.Fatal("expected error for malformed token, got nil")
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	svc := newTestService()

	_, err := svc.ValidateToken("")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestValidateToken_ClaimsPreserved(t *testing.T) {
	svc := newTestService()

	tests := []struct {
		name     string
		userID   string
		tenantID string
		role     string
	}{
		{"owner", "u-1", "t-1", "owner"},
		{"admin", "u-2", "t-2", "admin"},
		{"cashier", "u-3", "t-3", "cashier"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := svc.generateAccessToken(tt.userID, tt.tenantID, tt.role)
			if err != nil {
				t.Fatalf("generateAccessToken error: %v", err)
			}
			claims, err := svc.ValidateToken(token)
			if err != nil {
				t.Fatalf("ValidateToken error: %v", err)
			}
			if claims.UserID != tt.userID {
				t.Errorf("UserID = %q, want %q", claims.UserID, tt.userID)
			}
			if claims.TenantID != tt.tenantID {
				t.Errorf("TenantID = %q, want %q", claims.TenantID, tt.tenantID)
			}
			if claims.Role != tt.role {
				t.Errorf("Role = %q, want %q", claims.Role, tt.role)
			}
		})
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "MySecure1Pass"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		t.Fatalf("bcrypt.GenerateFromPassword error: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword(hash, []byte(password)); err != nil {
		t.Errorf("correct password should match hash, got error: %v", err)
	}

	if err := bcrypt.CompareHashAndPassword(hash, []byte("WrongPassword1")); err == nil {
		t.Error("wrong password should not match hash")
	}
}

func TestPasswordHashing_DifferentHashesForSamePassword(t *testing.T) {
	password := "SamePassword1"
	hash1, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		t.Fatalf("first hash error: %v", err)
	}
	hash2, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		t.Fatalf("second hash error: %v", err)
	}

	if string(hash1) == string(hash2) {
		t.Error("bcrypt should produce different hashes for the same password (salt)")
	}
}

func TestHashToken(t *testing.T) {
	token := "abc123"
	h1 := hashToken(token)
	h2 := hashToken(token)
	if h1 != h2 {
		t.Error("hashToken should be deterministic")
	}
	if h1 == "" {
		t.Error("hashToken should not return empty string")
	}

	different := hashToken("different-token")
	if h1 == different {
		t.Error("different tokens should produce different hashes")
	}
}

func TestGenerateRandomToken(t *testing.T) {
	tok1, err := generateRandomToken(32)
	if err != nil {
		t.Fatalf("generateRandomToken error: %v", err)
	}
	if len(tok1) != 64 {
		t.Errorf("expected 64 hex chars for 32 bytes, got %d", len(tok1))
	}

	tok2, err := generateRandomToken(32)
	if err != nil {
		t.Fatalf("generateRandomToken error: %v", err)
	}
	if tok1 == tok2 {
		t.Error("two random tokens should differ")
	}
}

func TestTokenExpiry(t *testing.T) {
	svc := newTestService()

	token, err := svc.generateAccessToken("u", "t", "r")
	if err != nil {
		t.Fatalf("generateAccessToken error: %v", err)
	}

	claims, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken error: %v", err)
	}

	if claims.ExpiresAt == nil {
		t.Fatal("token should have an expiry")
	}
	expiry := claims.ExpiresAt.Time
	now := time.Now()
	if expiry.Before(now) {
		t.Error("fresh token should not be expired")
	}
	if expiry.After(now.Add(accessTokenDuration + time.Minute)) {
		t.Error("token expiry should be within accessTokenDuration")
	}
}
