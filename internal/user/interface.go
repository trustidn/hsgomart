package user

type ServiceInterface interface {
	CreateUser(tenantID string, in CreateUserInput) (*User, error)
	ListUsers(tenantID string) ([]User, error)
	GetUser(tenantID, userID string) (*User, error)
	UpdateUser(tenantID, userID string, in UpdateUserInput) (*User, error)
	DeleteUser(tenantID, currentUserID, targetUserID string) error
}

var _ ServiceInterface = (*Service)(nil)
