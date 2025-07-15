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