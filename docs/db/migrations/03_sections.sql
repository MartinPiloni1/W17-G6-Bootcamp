CREATE TABLE sections (
    id INT NOT NULL AUTO_INCREMENT,
    section_number VARCHAR(255) NOT NULL UNIQUE,
    current_temperature DECIMAL(5,2) NOT NULL,
    minimum_temperature DECIMAL(5,2) NOT NULL,
    current_capacity INT NOT NULL,
    minimum_capacity INT NOT NULL,
    maximum_capacity INT NOT NULL,
    warehouse_id INT NOT NULL,
    product_type_id INT NOT NULL,
    PRIMARY KEY (id)
);

INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity, minimum_capacity, maximum_capacity, warehouse_id, product_type_id) VALUES
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
