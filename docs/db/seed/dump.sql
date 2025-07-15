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
    ('Pechuga de pollo', 4.50, 0.75, 3.0, 12.0, 8.0, 1.20, 'POL-0001', -18.0, 1, 1),
    ('Salmón', 6.00, 0.60, 2.5, 18.0, 14.0, 0.35, 'SAL-0001', -20.0, 2,  2),
    ('Leche entera', 7.00, 0.90, 25.0, 7.0, 7.0, 1.00, 'LCH-100L',  4.0, 3,   3),
    ('Yogurt helado', 5, 3, 6.4, 4.5, 1.2, 0.5, 'YOG01', -18, 4, 4)
