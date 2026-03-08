CREATE TABLE saas_settings (
    id INT PRIMARY KEY DEFAULT 1 CHECK (id = 1),
    saas_name VARCHAR(255) NOT NULL DEFAULT 'HSMart POS',
    logo_url TEXT,
    tagline TEXT,
    bank_name VARCHAR(100),
    bank_account VARCHAR(50),
    bank_holder VARCHAR(100),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO saas_settings (saas_name) VALUES ('HSMart POS');
