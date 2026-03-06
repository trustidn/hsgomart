CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    plan_id INT NOT NULL REFERENCES plans(id),
    status VARCHAR(50),
    start_date TIMESTAMP,
    end_date TIMESTAMP
);
