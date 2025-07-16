DROP DATABASE IF EXISTS fresh;

CREATE DATABASE fresh DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE fresh;

CREATE TABLE IF NOT EXISTS `buyers` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `card_number_id` INT NOT NULL UNIQUE,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name`  VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS `warehouses` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `warehouse_code` VARCHAR(50) NOT NULL UNIQUE,
    `address` VARCHAR(100) NOT NULL,
    `telephone` VARCHAR(100) NOT NULL,
    `minimun_capacity` INT NOT NULL,
    `minimun_temperature` FLOAT NOT NULL,
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS localities (
    id VARCHAR(20) PRIMARY KEY,
    locality_name VARCHAR(100) NOT NULL,
    province_name VARCHAR(100) NOT NULL,
    country_name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS employees (
    id INT NOT NULL AUTO_INCREMENT,
    card_number_id VARCHAR(50) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    warehouse_id INT NOT NULL,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS sellers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cid INT NOT NULL UNIQUE,
    company_name VARCHAR(100) NOT NULL,
    address VARCHAR(150) NOT NULL,
    telephone VARCHAR(25) NOT NULL,
    locality_id VARCHAR(20) NOT NULL,
    FOREIGN KEY (locality_id) REFERENCES localities(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS products (
    id                               INT AUTO_INCREMENT PRIMARY KEY,
    description                      TEXT NOT NULL,
    expiration_rate                  INTEGER   NOT NULL,  
    freezing_rate                    INTEGER   NOT NULL, 
    height                           DECIMAL(10,2) NOT NULL,  
    length                           DECIMAL(10,2) NOT NULL,  
    width                            DECIMAL(10,2) NOT NULL,  
    netweight                        DECIMAL(10,3) NOT NULL,
    product_code                     VARCHAR(50)  NOT NULL UNIQUE,
    recommended_freezing_temperature DECIMAL(5,2) NOT NULL,
    product_type_id                  INT NOT NULL,
    seller_id                        INT DEFAULT NULL,
    CONSTRAINT fk_products_seller
        FOREIGN KEY (seller_id)
        REFERENCES sellers(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS sections (
    id INT NOT NULL AUTO_INCREMENT,
    section_number VARCHAR(255) NOT NULL UNIQUE,
    current_temperature DECIMAL(5,2) NOT NULL,
    minimum_temperature DECIMAL(5,2) NOT NULL,
    current_capacity INT NOT NULL,
    minimum_capacity INT NOT NULL,
    maximum_capacity INT NOT NULL,
    warehouse_id INT NOT NULL,
    product_type_id INT NOT NULL,
    FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS product_records (
    id               INT           NOT NULL AUTO_INCREMENT,
    last_update_date DATE          NOT NULL,
    purchase_price   DECIMAL(12,2) NOT NULL,
    sale_price       DECIMAL(12,2) NOT NULL,
    product_id       INT           NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (product_id)
        REFERENCES products(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS purchase_orders (
    id                INT AUTO_INCREMENT PRIMARY KEY,
    order_number      VARCHAR(255) NOT NULL UNIQUE,
    order_date        DATETIME NOT NULL,
    tracking_code     VARCHAR(255) NOT NULL,
    buyer_id          INT NOT NULL,
    product_record_id INT NOT NULL,
    CONSTRAINT fk_purchase_orders_buyer
        FOREIGN KEY (buyer_id) REFERENCES buyers(id),
    CONSTRAINT fk_purchase_orders_product
        FOREIGN KEY (product_record_id) REFERENCES product_records(id)
);

CREATE TABLE IF NOT EXISTS carries (
    `id` int PRIMARY KEY AUTO_INCREMENT,
    `cid` VARCHAR(20) NOT NULL UNIQUE,
    `company_name` VARCHAR(100) NOT NULL,
    `address` VARCHAR(100) NOT NULL,
    `telephone` VARCHAR(100) NOT NULL,
    `locality_id` VARCHAR(20) NOT NULL,
    FOREIGN KEY (locality_id) REFERENCES localities(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS product_batches (
    id INT NOT NULL AUTO_INCREMENT,
    batch_number INT NOT NULL UNIQUE,
    current_quantity INT NOT NULL,
    current_temperature DECIMAL(5,2) NOT NULL,
    due_date DATE NOT NULL,
    initial_quantity INT NOT NULL,
    manufacturing_date DATE NOT NULL,
    manufacturing_hour INT NOT NULL,
    minimum_temperature DECIMAL(5,2) NOT NULL,
    product_id INT NOT NULL,
    section_id INT NOT NULL,
    FOREIGN KEY (section_id) REFERENCES sections(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS inbound_orders (
        id SERIAL PRIMARY KEY,
        order_number VARCHAR(50) NOT NULL UNIQUE,
        order_date TIMESTAMP NOT NULL DEFAULT NOW(),
        employee_id INTEGER NOT NULL,
        warehouse_id INTEGER NOT NULL,
        product_batch_id INTEGER NOT NULL,
        FOREIGN KEY (employee_id) REFERENCES employees(id),
        FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
        FOREIGN KEY (product_batch_id) REFERENCES product_batches(id)
);

INSERT IGNORE INTO warehouses (warehouse_code, address, telephone, minimun_capacity, minimun_temperature) VALUES
    ("DHK", 'Monroe 860', '47470000', '10', '10'),
    ("CBA", 'Cordoba 1234', '3516000123', '25', '4'),
    ("DHMP", 'Monroe 860', '47470000', '10', '10'),
    ("PAT", 'Parque industrial sur', '02991543210', '30', '-2'),
    ("CBA2", 'Cordoba 5678', '3516000456', '20', '4'),
    ("DHK2", 'Monroe 1234', '47470001', '15', '10'),
    ("CBA3", 'Cordoba 91011', '3516000789', '25', '4');

INSERT IGNORE INTO localities (id, locality_name, province_name, country_name) VALUES
    ('6700', 'Lujan', 'Buenos Aires', 'Argentina'),
    ('1001', 'CABA', 'CABA', 'Argentina'),
    ('2000', 'Rosario', 'Santa Fe', 'Argentina'),
    ('5000', 'Córdoba', 'Córdoba', 'Argentina'),
    ('10115', 'Berlin', 'Berlin', 'Alemania'),
    ('28001', 'Madrid', 'Madrid', 'España'),
    ('11000', 'Montevideo', 'Montevideo', 'Uruguay');

INSERT IGNORE INTO buyers (card_number_id, first_name, last_name) VALUES
    (12345678, 'Juan', 'Pérez'),
    (23456789, 'Ana', 'Gómez'),
    (34567890, 'Luis', 'Martínez');

INSERT IGNORE INTO employees (card_number_id, first_name, last_name, warehouse_id) VALUES
    ('ABC001', 'Ramon', 'Diaz', 1),
    ('DEF002', 'Carlos', 'Lopez', 1),
    ('GHI003', 'Marta', 'Perez', 2),
    ('JKL004', 'Lucia', 'Romero', 2),
    ('MNO005', 'Sergio', 'Castro', 3);

INSERT IGNORE INTO sellers (cid, company_name, address, telephone, locality_id) VALUES
    (1, 'Alkemy', 'Monroe 860', '47470000', '6700'),
    (2, 'Globant', 'Av. Córdoba 1200', '40334000', '1001'),
    (3, 'Mercado Libre', 'Alem 876', '45450000', '1001'),
    (4, 'Tech Solutions', 'Av. Pellegrini 900', '3411234567', '2000'),
    (5, 'Panaderia El Sol', 'Av. Colon 1000', '3511112222', '5000'),
    (6, 'Bäckerei Berlin', 'Unter den Linden 77', '+4930123456', '10115'),
    (7, 'Supermercado Español', 'Calle Gran Vía 1', '34911223344', '28001'),
    (8, 'Chivitería El Prado', 'Av. 18 de Julio 1010', '59829123456', '11000');

INSERT IGNORE INTO carries (cid, company_name, address, telephone, locality_id) VALUES
    ("CID1217", 'Alkemy', 'Monroe 860', '47470000', '6700'),
    ("CID2332", 'Globant', 'Av. Córdoba 1200', '40334000', '1001'),
    ("CID#232", 'Mercado Libre', 'Alem 876', '45450000', '1001'),
    ("CID4434", 'Tech Solutions', 'Av. Pellegrini 900', '3411234567', '2000'),
    ("CID2445", 'Panaderia El Sol', 'Av. Colon 1000', '3511112222', '5000'),
    ("CID3326", 'Bäckerei Berlin', 'Unter den Linden 77', '+4930123456', '10115');

INSERT IGNORE INTO products (description, expiration_rate, freezing_rate, height, length, width, netweight, product_code, recommended_freezing_temperature, product_type_id, seller_id) VALUES
    ('Pechuga de pollo', 4.50, 0.75, 3.0, 12.0, 8.0, 1.20, 'POL-0001', -18.0, 1, 1),
    ('Salmón', 6.00, 0.60, 2.5, 18.0, 14.0, 0.35, 'SAL-0001', -20.0, 2,  2),
    ('Leche entera', 7.00, 0.90, 25.0, 7.0, 7.0, 1.00, 'LCH-100L',  4.0, 3,   3),
    ('Yogurt helado', 5, 3, 6.4, 4.5, 1.2, 0.5, 'YOG01', -18, 4, 4);

INSERT IGNORE INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, product_type_id) VALUES
    ('SEC-001', 5.00, 0.00, 50, 0, 100, 1, 1),
    ('SEC-002', -18.00, -20.00, 20, 0, 50, 1, 2),
    ('SEC-003', 20.00, 18.00, 80, 50, 150, 2, 3),
    ('SEC-004', 8.50, 5.00, 10, 0, 30, 2, 1),
    ('SEC-005', 1.00, -2.00, 75, 20, 100, 3, 2),
    ('SEC-006', 22.00, 20.00, 120, 100, 200, 3, 3),
    ('SEC-007', 0.00, -5.00, 30, 10, 50, 1, 3),
    ('SEC-008', 15.50, 10.00, 60, 40, 80, 2, 2),
    ('SEC-009', -1.00, -3.00, 5, 0, 10, 3, 1),
    ('SEC-010', 10.00, 8.00, 90, 70, 110, 1, 1);

INSERT IGNORE INTO product_records (last_update_date, purchase_price, sale_price, product_id) VALUES
    ('2024-01-10',  10.00,  12.50,  1),
    ('2024-01-12',  20.00,  20.00,  2),
    ('2024-02-05',   5.75,   8.99,  3),
    ('2024-03-20', 100.00, 120.00,  4),
    ('2024-04-01',  9.00,  11.5,  1);

INSERT IGNORE INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minimum_temperature, product_id, section_id) VALUES
    (111, 500, 4.5, '2025-08-15', 500, '2025-07-15', 8, 2.0, 1, 1),
    (112, 1000, -18.0, '2026-01-20', 1000, '2025-07-15', 10, -22.0, 2, 2),
    (113, 300, 5.0, '2025-11-10', 400, '2025-07-14', 14, 3.0, 3, 1),
    (114, 450, 4.2, '2025-08-22', 500, '2025-07-22', 9, 2.0, 1, 3),
    (115, 800, -19.5, '2026-02-01', 800, '2025-07-20', 11, -22.0, 2, 2);

INSERT IGNORE INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id) VALUES
        ('ORD-001', '2024-06-15 12:34:56', 'TRACK-111AAA', 1, 1),
        ('ORD-002', '2024-06-16 15:10:35', 'TRACK-222BBB', 2, 2);

INSERT IGNORE INTO inbound_orders (order_number, order_date, employee_id, warehouse_id, product_batch_id) VALUES
    ('INB-1001', '2024-06-01 09:00:00', 1, 1, 1),
    ('INB-1002', '2024-06-02 10:30:00', 2, 2, 2),
    ('INB-1003', '2024-06-03 11:00:00', 3, 3, 3);