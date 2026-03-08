ALTER TABLE plans ADD COLUMN IF NOT EXISTS duration_days INT NOT NULL DEFAULT 30;

UPDATE plans SET duration_days = 7 WHERE LOWER(name) = 'trial';
UPDATE plans SET duration_days = 30 WHERE LOWER(name) = 'basic';
UPDATE plans SET duration_days = 90 WHERE LOWER(name) = 'professional';
UPDATE plans SET duration_days = 365 WHERE LOWER(name) = 'enterprise';
