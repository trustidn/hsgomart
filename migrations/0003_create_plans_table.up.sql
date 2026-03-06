CREATE TABLE plans (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    price NUMERIC,
    max_users INT,
    max_products INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
