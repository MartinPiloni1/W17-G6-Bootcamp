INSERT INTO sellers (id, cid, company_name, address, telephone, locality_id) VALUES
     (1, 1001, 'Acme Corp', 'Calle Luna 123, Madrid', '+34-911-222-333', 'LOC001'),
     (2, 1002, 'Distribuciones Sol', 'Av. Sol 456, Barcelona', '+34-922-333-444', 'LOC002'),
     (3, 1003, 'Proveedor Norte', 'Pza. Norte 789, Bilbao', '+34-933-444-555', 'LOC003');

INSERT INTO warehouses (id, warehouse_code, address, telephone, minimun_capacity, minimun_temperature) VALUES
   (1, 'WH-001', '123 Main St', '+34-111-222-333', 1000, 5.0),
   (2, 'WH-002', '456 Secondary Ave', '+34-222-333-444', 2000, 3.5),
   (3, 'WH-003', '789 Warehouse Rd', '+34-333-444-555', 1500, 2.0);

INSERT INTO employees (card_number_id, first_name, last_name, warehouse_id) VALUES
    ('10101010', 'Ramon', 'Diaz', 1),
    ('10101011', 'Carlos', 'Lopez', 1),
    ('10101012', 'Marta', 'Perez', 2),
    ('10101013', 'Lucia', 'Romero', 2),
                                                                                ('10101014', 'Sergio', 'Castro', 3);

INSERT INTO product_batch (id) VALUES
   (1),
   (2),
   (3);

INSERT INTO inbound_orders (order_number, order_date, employee_id, warehouse_id, product_batch_id) VALUES
   ('INB-1001', '2024-06-01 09:00:00', 1, 1, 1),
   ('INB-1002', '2024-06-02 10:30:00', 2, 2, 2),
   ('INB-1003', '2024-06-03 11:00:00', 3, 3, 3);