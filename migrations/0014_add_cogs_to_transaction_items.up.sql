-- Store COGS (from FIFO batch cost) per transaction item for accurate profit reporting
ALTER TABLE transaction_items ADD COLUMN IF NOT EXISTS cogs NUMERIC NOT NULL DEFAULT 0;
