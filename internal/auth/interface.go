package auth

type ServiceInterface interface {
	Register(in RegisterInput) (*TokenPair, error)
	Login(in LoginInput) (*TokenPair, error)
	Refresh(rawToken string) (*TokenPair, error)
	RevokeRefreshTokens(userID string) error
	ValidateToken(tokenString string) (*JWTClaims, error)
	UserEmail(userID string) (string, error)
}

var _ ServiceInterface = (*Service)(nil)
