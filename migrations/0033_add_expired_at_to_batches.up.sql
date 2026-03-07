ALTER TABLE inventory_batches ADD COLUMN IF NOT EXISTS expired_at DATE;
CREATE INDEX IF NOT EXISTS idx_inventory_batches_expired ON inventory_batches(expired_at);
