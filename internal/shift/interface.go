package shift

type ServiceInterface interface {
	OpenShift(tenantID, userID string, in OpenShiftInput) (*OpenShiftResult, error)
	CloseShift(tenantID, userID, shiftID string, in CloseShiftInput) (*CloseShiftResult, error)
	GetCurrentShift(tenantID, userID string) (*CashierShift, error)
	ListShifts(tenantID string, limit, offset int) ([]CashierShift, error)
}

var _ ServiceInterface = (*Service)(nil)
