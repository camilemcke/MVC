CREATE TABLE IF NOT EXISTS users (
    id    SERIAL PRIMARY KEY,
    name  TEXT NOT NULL,
    email TEXT NOT NULL
);

INSERT INTO users (name, email) VALUES ('Camil', 'camilemcke@gmail.com');

CREATE TABLE IF NOT EXISTS products (
    id       SERIAL PRIMARY KEY,
    sku      TEXT NOT NULL UNIQUE,
    name     TEXT NOT NULL,
    brand    TEXT NOT NULL,
    category TEXT NOT NULL
);

INSERT INTO products (sku, name, brand, category) VALUES
('OLAPLEX-001', 'No. 4 Bond Maintenance Shampoo', 'Olaplex', 'Shampoo'),
('K18-001', 'Leave-In Molecular Repair Hair Mask', 'K18', 'Tratamiento'),
('WOW-001', 'Dream Coat Supernatural Spray', 'Color Wow', 'Tratamiento');

CREATE TABLE IF NOT EXISTS inventory (
    id         SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    stock      INT NOT NULL DEFAULT 0
);

INSERT INTO inventory (product_id, stock) VALUES
(1, 24),
(2, 12),
(3, 5);