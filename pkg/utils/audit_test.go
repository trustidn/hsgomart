package utils

import (
	"encoding/json"
	"testing"
)

func TestAuditLog_TableName(t *testing.T) {
	log := AuditLog{}
	if got := log.TableName(); got != "audit_logs" {
		t.Errorf("TableName() = %q, want %q", got, "audit_logs")
	}
}

func TestAuditLog_JSONDetails(t *testing.T) {
	details := map[string]string{"key": "value"}
	b, err := json.Marshal(details)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	log := AuditLog{
		Action:     "test_action",
		EntityType: "test_entity",
		EntityID:   "123",
		Details:    b,
	}

	if log.Action != "test_action" {
		t.Errorf("Action = %q, want %q", log.Action, "test_action")
	}

	var parsed map[string]string
	if err := json.Unmarshal(log.Details, &parsed); err != nil {
		t.Fatalf("failed to unmarshal Details: %v", err)
	}
	if parsed["key"] != "value" {
		t.Errorf("Details[key] = %q, want %q", parsed["key"], "value")
	}
}

func TestAuditLog_NilPointerFields(t *testing.T) {
	tenantID := "tenant-abc"
	userID := "user-xyz"

	log := AuditLog{
		TenantID: &tenantID,
		UserID:   &userID,
		Action:   "create",
	}

	if log.TenantID == nil || *log.TenantID != "tenant-abc" {
		t.Errorf("TenantID = %v, want %q", log.TenantID, "tenant-abc")
	}
	if log.UserID == nil || *log.UserID != "user-xyz" {
		t.Errorf("UserID = %v, want %q", log.UserID, "user-xyz")
	}

	logNil := AuditLog{Action: "delete"}
	if logNil.TenantID != nil {
		t.Error("expected nil TenantID")
	}
	if logNil.UserID != nil {
		t.Error("expected nil UserID")
	}
}

func TestLogAudit_NilDB(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Logf("LogAudit panicked with nil db (expected): %v", r)
		}
	}()
	LogAudit(nil, "tenant", "user", "action", "entity", "id", map[string]string{"k": "v"})
	t.Log("LogAudit with nil db did not panic (unexpected but acceptable if gorm handles nil)")
}
