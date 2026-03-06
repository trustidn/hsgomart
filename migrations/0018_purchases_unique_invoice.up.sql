-- Prevent duplicate purchase entry per tenant by invoice number
CREATE UNIQUE INDEX IF NOT EXISTS idx_purchases_tenant_invoice ON purchases(tenant_id, invoice_number)
WHERE invoice_number IS NOT NULL AND invoice_number != '';
