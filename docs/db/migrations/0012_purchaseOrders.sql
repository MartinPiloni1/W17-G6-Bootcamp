CREATE TABLE IF NOT EXISTS purchase_orders (
    id                INT AUTO_INCREMENT PRIMARY KEY,
    order_number      VARCHAR(255) NOT NULL UNIQUE,
    order_date        DATETIME NOT NULL,
    tracking_code     VARCHAR(255) NOT NULL,
    buyer_id          INT NOT NULL,
    product_record_id INT NOT NULL,
    CONSTRAINT fk_purchase_orders_buyer
        FOREIGN KEY (buyer_id) REFERENCES buyers(id),
    CONSTRAINT fk_purchase_orders_product
        FOREIGN KEY (product_record_id) REFERENCES product_records(id)
);