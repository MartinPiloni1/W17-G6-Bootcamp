CREATE TABLE inbound_orders (
        id SERIAL PRIMARY KEY,
        order_number VARCHAR(50) NOT NULL UNIQUE,
        order_date TIMESTAMP NOT NULL DEFAULT NOW(),
        employee_id INTEGER NOT NULL,
        warehouse_id INTEGER NOT NULL,
        product_batch_id INTEGER NOT NULL,
        FOREIGN KEY (employee_id) REFERENCES employees(id),
        FOREIGN KEY (warehouse_id) REFERENCES warehouses(id),
        FOREIGN KEY (product_batch_id) REFERENCES product_batch(id)
);