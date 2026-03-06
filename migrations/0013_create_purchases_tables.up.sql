-- purchases: header for each purchase
CREATE TABLE purchases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    supplier_name VARCHAR(255),
    invoice_number VARCHAR(255),
    total_amount NUMERIC NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_purchases_tenant_id ON purchases(tenant_id);
CREATE INDEX idx_purchases_created_at ON purchases(created_at);

-- purchase_items: line items per purchase
CREATE TABLE purchase_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    purchase_id UUID NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    quantity INT NOT NULL CHECK (quantity > 0),
    cost_price NUMERIC NOT NULL CHECK (cost_price >= 0),
    subtotal NUMERIC NOT NULL CHECK (subtotal >= 0)
);

CREATE INDEX idx_purchase_items_purchase_id ON purchase_items(purchase_id);
CREATE INDEX idx_purchase_items_product_id ON purchase_items(product_id);

-- inventory_batches: FIFO batches per product (linked to purchase_item)
CREATE TABLE inventory_batches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id),
    purchase_item_id UUID NOT NULL REFERENCES purchase_items(id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
    remaining_quantity INT NOT NULL CHECK (remaining_quantity >= 0),
    cost_price NUMERIC NOT NULL CHECK (cost_price >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_inventory_batches_product_id ON inventory_batches(product_id);
CREATE INDEX idx_inventory_batches_remaining ON inventory_batches(product_id, created_at);
