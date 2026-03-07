ALTER TABLE transaction_items DROP COLUMN IF EXISTS discount_type;
ALTER TABLE transaction_items DROP COLUMN IF EXISTS discount_value;
ALTER TABLE transaction_items DROP COLUMN IF EXISTS discount_amount;
ALTER TABLE transactions DROP COLUMN IF EXISTS discount_amount;
