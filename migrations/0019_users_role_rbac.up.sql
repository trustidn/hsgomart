-- Restrict role to valid values and set default (existing NULL/invalid become cashier)
UPDATE users SET role = 'cashier' WHERE role IS NULL OR role NOT IN ('owner', 'cashier');
ALTER TABLE users ADD CONSTRAINT users_role_check CHECK (role IN ('owner', 'cashier'));
ALTER TABLE users ALTER COLUMN role SET DEFAULT 'cashier';
