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
    PRIMARY KEY (id)
    );