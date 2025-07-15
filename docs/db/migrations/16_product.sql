CREATE TABLE IF NOT EXISTS products (
    id                               INT AUTO_INCREMENT PRIMARY KEY,
    description                      TEXT NOT NULL,
    expiration_rate                  INTEGER   NOT NULL,  
    freezing_rate                    INTEGER   NOT NULL, 
    height                           DECIMAL(10,2) NOT NULL,  
    length                           DECIMAL(10,2) NOT NULL,  
    width                            DECIMAL(10,2) NOT NULL,  
    netweight                        DECIMAL(10,3) NOT NULL,
    product_code                     VARCHAR(50)  NOT NULL UNIQUE,
    recommended_freezing_temperature DECIMAL(5,2) NOT NULL,
    product_type_id                  INT NOT NULL,
    seller_id                        INT NULL,
    CONSTRAINT fk_products_seller
        FOREIGN KEY (seller_id)
        REFERENCES sellers(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);
