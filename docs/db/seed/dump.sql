INSERT INTO employees (card_number_id, first_name, last_name, warehouse_id) VALUES
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
    ('Pechuga de pollo', 4.50, 0.75, 3.0, 12.0, 8.0, 1.20, 'POL-0001', -18.0, 1, 10),
    ('Salmón', 6.00, 0.60, 2.5, 18.0, 14.0, 0.35, 'SAL-0001', -20.0, 2,  3),
    ('Leche entera', 7.00, 0.90, 25.0, 7.0, 7.0, 1.00, 'LCH-100L',  4.0, 3,   5),
    ('Yogurt helado', 5, 3, 6.4, 4.5, 1.2, 0.5, 'YOG01', -18, 4, 2);

INSERT INTO product_records (
  last_update_date, 
  purchase_price, 
  sale_price, 
  product_id
)
VALUES
  ('2024-01-10',  10.00,  12.50,  1),
  ('2024-01-12',  20.00,  20.00,  2),  
  ('2024-02-05',   5.75,   8.99,  3),
  ('2024-03-20', 100.00, 120.00,  4),
  ('2024-04-01',  9.00,  11.5,  1);

INSERT INTO `buyers` (`card_number_id`, `first_name`, `last_name`) VALUES
        (12345678, 'Juan', 'Pérez'),
        (23456789, 'Ana', 'Gómez'),
        (34567890, 'Luis', 'Martínez');
