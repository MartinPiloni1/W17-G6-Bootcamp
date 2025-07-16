INSERT INTO employees (card_number_id, first_name, last_name, warehouse_id) 
VALUES
        ('ABC001', 'Ramon', 'Diaz', 1),
        ('DEF002', 'Carlos', 'Lopez', 1),
        ('GHI003', 'Marta', 'Perez', 2),
        ('JKL004', 'Lucia', 'Romero', 2),
        ('MNO005', 'Sergio', 'Castro', 3);

INSERT INTO products (
    description,
    expiration_rate,
    freezing_rate,
    height,
    length,
    width,
    netweight,
    product_code,
    recommended_freezing_temperature,
    product_type_id,
    seller_id
)
VALUES
    ('Pechuga de pollo', 4.50, 0.75, 3.0, 12.0, 8.0, 1.20, 'POL-0001', -18.0, 1, 1),
    ('Salmón', 6.00, 0.60, 2.5, 18.0, 14.0, 0.35, 'SAL-0001', -20.0, 2,  2),
    ('Leche entera', 7.00, 0.90, 25.0, 7.0, 7.0, 1.00, 'LCH-100L',  4.0, 3,   3),
    ('Yogurt helado', 5, 3, 6.4, 4.5, 1.2, 0.5, 'YOG01', -18, 4, 4)

INSERT INTO `buyers` (`card_number_id`, `first_name`, `last_name`)
VALUES
        (12345678, 'Juan', 'Pérez'),
        (23456789, 'Ana', 'Gómez'),
        (34567890, 'Luis', 'Martínez');


INSERT INTO localities (id, locality_name, province_name, country_name) VALUES
('6700', 'Lujan', 'Buenos Aires', 'Argentina'),
('1001', 'CABA', 'CABA', 'Argentina'),
('2000', 'Rosario', 'Santa Fe', 'Argentina'),
('5000', 'Córdoba', 'Córdoba', 'Argentina'),
('10115', 'Berlin', 'Berlin', 'Alemania'),
('28001', 'Madrid', 'Madrid', 'España'),
('11000', 'Montevideo', 'Montevideo', 'Uruguay');

INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES
(1, 'Alkemy', 'Monroe 860', '47470000', '6700'),
(2, 'Globant', 'Av. Córdoba 1200', '40334000', '1001'),
(3, 'Mercado Libre', 'Alem 876', '45450000', '1001'),
(4, 'Tech Solutions', 'Av. Pellegrini 900', '3411234567', '2000'),
(5, 'Panaderia El Sol', 'Av. Colon 1000', '3511112222', '5000'),
(6, 'Bäckerei Berlin', 'Unter den Linden 77', '+4930123456', '10115'),
(7, 'Supermercado Español', 'Calle Gran Vía 1', '34911223344', '28001'),
(8, 'Chivitería El Prado', 'Av. 18 de Julio 1010', '59829123456', '11000');



INSERT INTO carries (cid, company_name, address, telephone, locality_id) VALUES
("CID1217", 'Alkemy', 'Monroe 860', '47470000', '6700'),
("CID2332", 'Globant', 'Av. Córdoba 1200', '40334000', '1001'),
("CID#232", 'Mercado Libre', 'Alem 876', '45450000', '1001'),
("CID4434", 'Tech Solutions', 'Av. Pellegrini 900', '3411234567', '2000'),
("CID2445", 'Panaderia El Sol', 'Av. Colon 1000', '3511112222', '5000'),
("CID3326", 'Bäckerei Berlin', 'Unter den Linden 77', '+4930123456', '10115');


INSERT INTO warehouse (warehouse_code, address, telephone, minimun_capacity, minimun_temperature) VALUES
    ("DHK", 'Monroe 860', '47470000', '10', '10'),
    ("CBA", 'Cordoba 1234', '3516000123', '25', '4'),
    ("DHMP", 'Monroe 860', '47470000', '10', '10'),
    ("PAT", 'Parque industrial sur', '02991543210', '30', '-2'),
    ("CBA2", 'Cordoba 5678', '3516000456', '20', '4'),
    ("DHK2", 'Monroe 1234', '47470001', '15', '10'),
    ("CBA3", 'Cordoba 91011', '3516000789', '25', '4');,

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