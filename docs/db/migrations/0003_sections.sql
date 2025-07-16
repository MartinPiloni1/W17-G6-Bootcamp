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