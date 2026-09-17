CREATE TABLE IF NOT EXISTS users (
    id    SERIAL PRIMARY KEY,
    name  TEXT NOT NULL,
    email TEXT NOT NULL
);

INSERT INTO users (name, email) VALUES ('Camil', 'camilemcke@gmail.com');

CREATE TABLE IF NOT EXISTS products (
    id    SERIAL PRIMARY KEY,
    sku   TEXT NOT NULL UNIQUE,
    name  TEXT NOT NULL,
    brand TEXT NOT NULL
);

INSERT INTO products (sku, name, brand) VALUES
('22P16', '12oz/350mil Brazilian Blowout Original Smoothing Solution', 'Brazilian Blowout'),
('11M71', 'BB Protective Eyeglasses', 'Brazilian Blowout'),
('91008.0', 'Curl Quencher® Moisturizing Shampoo 8.5 oz 250 ml', 'Ouidad'),
('98506.0', 'VitalCurl™ Plus Soft Defining Mousse', 'Ouidad'),
('90732.0', 'Botanical Boost Curl Energizing & Refreshing Spray 33.8 Oz. 1L', 'Ouidad'),
('2493562.0', 'Caviar Anti-Aging Smoothing Anti-Frizz Multi-Styling Air Dry Balm 3.4 oz', 'ALTERNA'),
('2683462.0', 'ALTERNA CAVIAR ANTI-AGING RESTRUCTURING BOND REPAIR SHAMPOO 8.5 OZ', 'ALTERNA'),
('2639671.0', 'ALTERNA CAVIAR ANTI-AGING REPLENISHING MOISTURE CONDITIONER 8.5 OZ', 'ALTERNA'),
('5060150182273.0', 'Dream coat Supernatural Spray 50ml', 'Color Wow'),
('L3047', 'LUSETA KERATIN OIL WEIGHTLESS SMOOTHING HAIR REPAIR SERUM 3.4OZ', 'LUSETA'),
('L3035S', 'LUSETA ARGAN OIL CONDITIONER 16.9OZ', 'LUSETA'),
('241051.0', 'PROCUTASE Ionic Hydrogel skin-protecting and repairing treatment spray - 100 ml', 'Bionike'),
('D16945', 'Mini size SHINE ON Nutri-hair restructuring conditioner', 'Bionike'),
('22111.0', 'BioNike AKNET AZELIKE plus - Azelaic Ac', 'Bionike'),
('PFSVP853', 'AGE OPTIMAL FACE AND EYES WRINKLE SMOOTHING CREAM 100ML', 'Phytomer'),
('PFESV139E', 'PHYTOMER MUESTRAS SITADINE CREME', 'Phytomer'),
('PFSVV395', 'PIONNIERE XMF RADIANCE RETEXTURING SERUM 30 ML', 'Phytomer'),
('2175B', 'Rosalieve redness reducing 30ml', 'JAN MARINI'),
('2109D', 'Luminate eye gel 15ml', 'JAN MARINI'),
('2084D', 'Age intervention retinol plus 28g', 'JAN MARINI'),
('798-CC-1', 'OLIVIA GARDEN COMB CLEANER', 'Olivia Garden'),
('721-OGMC01', 'Mini OG Brush Styler Smooth & Shine - Black', 'Olivia Garden'),
('710-T16', 'ProThermal 3/4"', 'Olivia Garden'),
('VBCOMBS', 'Peine Para Cabello', 'Viemaa'),
('VSBAGS', 'SHOPPING BAGS', 'Viemaa');

CREATE TABLE IF NOT EXISTS inventory (
    id         SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id),
    entries    INT NOT NULL DEFAULT 0,
    exits      INT NOT NULL DEFAULT 0,
    stock      INT NOT NULL DEFAULT 0
);

INSERT INTO inventory (product_id, entries, exits, stock) VALUES
(1, 0, 0, 0),
(2, 218, 4, 214),
(3, 153, 8, 145),
(4, 120, 4, 116),
(5, 4, 4, 0),
(6, 11, 11, 0),
(7, 96, 1, 95),
(8, 94, 1, 93),
(9, 63, 63, 0),
(10, 30, 30, 0),
(11, 37, 21, 16),
(12, 1, 0, 1),
(13, 105, 0, 105),
(14, 96, 0, 96),
(15, 2, 2, 0),
(16, 151, 0, 151),
(17, 54, 0, 54),
(18, 12, 2, 10),
(19, 24, 0, 24),
(20, 24, 0, 24),
(21, 1, 1, 0),
(22, 240, 2, 238),
(23, 240, 15, 225),
(24, 98, 0, 98),
(25, 75, 0, 75);