CREATE TABLE stock_opnames (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'draft',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE TABLE stock_opname_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    opname_id UUID NOT NULL REFERENCES stock_opnames(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    system_stock INT NOT NULL DEFAULT 0,
    actual_stock INT NOT NULL DEFAULT 0,
    difference INT NOT NULL DEFAULT 0,
    notes TEXT
);

CREATE INDEX idx_stock_opnames_tenant ON stock_opnames(tenant_id);
CREATE INDEX idx_stock_opname_items_opname ON stock_opname_items(opname_id);
