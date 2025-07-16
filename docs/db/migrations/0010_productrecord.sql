CREATE TABLE product_records (
  id               INT           NOT NULL AUTO_INCREMENT,
  last_update_date DATE          NOT NULL,
  purchase_price   DECIMAL(12,2) NOT NULL,
  sale_price       DECIMAL(12,2) NOT NULL,
  product_id       INT           NOT NULL,
  PRIMARY KEY (id),
  FOREIGN KEY (product_id)
    REFERENCES products(id)
    ON DELETE RESTRICT
    ON UPDATE CASCADE
);