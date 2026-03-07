DROP INDEX IF EXISTS idx_inventory_batches_expired;
ALTER TABLE inventory_batches DROP COLUMN IF EXISTS expired_at;
