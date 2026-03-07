package shift

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrShiftAlreadyOpen   = errors.New("you already have an open shift")
	ErrNoActiveShift      = errors.New("no active shift")
	ErrShiftNotFound      = errors.New("shift not found")
	ErrShiftAlreadyClosed = errors.New("shift is already closed")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type OpenShiftInput struct {
	OpeningCash float64 `json:"opening_cash" binding:"required,min=0"`
}

type OpenShiftResult struct {
	ShiftID     string  `json:"shift_id"`
	OpenedAt    string  `json:"opened_at"`
	OpeningCash float64 `json:"opening_cash"`
}

func (s *Service) OpenShift(tenantID, userID string, in OpenShiftInput) (*OpenShiftResult, error) {
	active, err := GetActiveByUser(s.db, tenantID, userID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return nil, ErrShiftAlreadyOpen
	}
	now := time.Now()
	sh := &CashierShift{
		TenantID:    tenantID,
		UserID:      userID,
		OpeningCash: in.OpeningCash,
		OpenedAt:    now,
		Status:      StatusOpen,
	}
	if err := Create(s.db, sh); err != nil {
		return nil, err
	}
	return &OpenShiftResult{
		ShiftID:     sh.ID,
		OpenedAt:    now.Format(time.RFC3339),
		OpeningCash: sh.OpeningCash,
	}, nil
}

type CloseShiftInput struct {
	ClosingCash *float64 `json:"closing_cash" binding:"required,min=0"` // pointer so 0 is valid (required = present)
}

type CloseShiftResult struct {
	ExpectedCash float64 `json:"expected_cash"`
	ActualCash   float64 `json:"actual_cash"`
	Difference   float64 `json:"difference"`
}

func (s *Service) CloseShift(tenantID, userID, shiftID string, in CloseShiftInput) (*CloseShiftResult, error) {
	sh, err := GetByID(s.db, tenantID, shiftID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrShiftNotFound
		}
		return nil, err
	}
	if sh.Status == StatusClosed {
		return nil, ErrShiftAlreadyClosed
	}
	if sh.UserID != userID {
		return nil, ErrShiftNotFound
	}
	if in.ClosingCash == nil {
		return nil, errors.New("closing_cash is required")
	}
	closingCash := *in.ClosingCash
	now := time.Now()
	cashSales, err := SumCashPaymentsBetween(s.db, tenantID, sh.OpenedAt, now)
	if err != nil {
		return nil, err
	}
	expected := sh.OpeningCash + cashSales
	if err := CloseShift(s.db, shiftID, closingCash, now); err != nil {
		return nil, err
	}
	return &CloseShiftResult{
		ExpectedCash: expected,
		ActualCash:   closingCash,
		Difference:   closingCash - expected,
	}, nil
}

func (s *Service) GetCurrentShift(tenantID, userID string) (*CashierShift, error) {
	return GetActiveByUser(s.db, tenantID, userID)
}

func (s *Service) ListShifts(tenantID string, limit, offset int) ([]CashierShift, error) {
	return ListShifts(s.db, tenantID, limit, offset)
}
