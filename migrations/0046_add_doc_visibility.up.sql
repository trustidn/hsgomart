ALTER TABLE documentation ADD COLUMN IF NOT EXISTS visibility VARCHAR(20) NOT NULL DEFAULT 'all';
UPDATE documentation SET visibility = 'admin' WHERE title = 'Panduan Deployment ke Server / VPS';
