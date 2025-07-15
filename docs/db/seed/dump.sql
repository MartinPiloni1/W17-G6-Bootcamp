INSERT IGNORE INTO employees (card_number_id, first_name, last_name, warehouse_id) VALUES
        ('ABC001', 'Ramon', 'Diaz', 1),
        ('DEF002', 'Carlos', 'Lopez', 1),
        ('GHI003', 'Marta', 'Perez', 2),
        ('JKL004', 'Lucia', 'Romero', 2),
        ('MNO005', 'Sergio', 'Castro', 3);

INSERT IGNORE INTO `buyers` (`card_number_id`, `first_name`, `last_name`) VALUES
        (12345678, 'Juan', 'Pérez'),
        (23456789, 'Ana', 'Gómez'),
        (34567890, 'Luis', 'Martínez');

INSERT IGNORE INTO purchase_orders
(order_number, order_date, tracking_code, buyer_id, product_record_id) VALUES
        ('ORD-001', '2024-06-15 12:34:56', 'TRACK-111AAA', 1, 1),
        ('ORD-002', '2024-06-16 15:10:35', 'TRACK-222BBB', 2, 2);