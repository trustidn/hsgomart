DROP INDEX IF EXISTS idx_products_deleted_at;
DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_categories_deleted_at;

ALTER TABLE products DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE products DROP COLUMN IF EXISTS updated_at;
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE users DROP COLUMN IF EXISTS updated_at;
ALTER TABLE categories DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE categories DROP COLUMN IF EXISTS updated_at;
