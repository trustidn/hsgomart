CREATE TABLE cashier_shifts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),

    opening_cash NUMERIC NOT NULL DEFAULT 0,
    closing_cash NUMERIC,

    opened_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    closed_at TIMESTAMP,

    status VARCHAR(20) NOT NULL DEFAULT 'open'
);

CREATE INDEX idx_cashier_shifts_tenant_user ON cashier_shifts(tenant_id, user_id);
CREATE INDEX idx_cashier_shifts_tenant_status ON cashier_shifts(tenant_id, status);
CREATE INDEX idx_cashier_shifts_opened_at ON cashier_shifts(opened_at);
