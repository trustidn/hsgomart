-- Documentation articles (managed by superadmin, read by all)
CREATE TABLE IF NOT EXISTS documentation (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Platform updates / changelog (managed by superadmin)
CREATE TABLE IF NOT EXISTS platform_updates (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Add WhatsApp number to saas_settings
ALTER TABLE saas_settings ADD COLUMN IF NOT EXISTS whatsapp_number VARCHAR(50) DEFAULT '';
