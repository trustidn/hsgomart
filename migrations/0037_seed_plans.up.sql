INSERT INTO plans (name, price, max_users, max_products)
SELECT 'Basic', 99000, 3, 50
WHERE NOT EXISTS (SELECT 1 FROM plans WHERE name = 'Basic');

INSERT INTO plans (name, price, max_users, max_products)
SELECT 'Professional', 199000, 10, 500
WHERE NOT EXISTS (SELECT 1 FROM plans WHERE name = 'Professional');

INSERT INTO plans (name, price, max_users, max_products)
SELECT 'Enterprise', 499000, 50, 5000
WHERE NOT EXISTS (SELECT 1 FROM plans WHERE name = 'Enterprise');
