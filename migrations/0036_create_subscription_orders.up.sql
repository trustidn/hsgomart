CREATE TABLE subscription_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    plan_id INT NOT NULL REFERENCES plans(id),
    amount NUMERIC NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_payment',
    payment_proof_url TEXT,
    invoice_number VARCHAR(50) NOT NULL,
    notes TEXT,
    admin_notes TEXT,
    reviewed_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    paid_at TIMESTAMP,
    reviewed_at TIMESTAMP
);

CREATE INDEX idx_subscription_orders_tenant ON subscription_orders(tenant_id);
CREATE INDEX idx_subscription_orders_status ON subscription_orders(status);
