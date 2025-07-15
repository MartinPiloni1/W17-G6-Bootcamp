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

INSERT IGNORE INTO purchase_orders
(order_number, order_date, tracking_code, buyer_id, product_record_id) VALUES
        ('ORD-001', '2024-06-15 12:34:56', 'TRACK-111AAA', 1, 1),
        ('ORD-002', '2024-06-16 15:10:35', 'TRACK-222BBB', 2, 2);