
CREATE TABLE products (
    id                               INT NOT NULL AUTO_INCREMENT,
    description                      TEXT NOT NULL,
    expiration_rate                  INTEGER   NOT NULL,  
    freezing_rate                    INTEGER   NOT NULL, 
    height                           NUMERIC(10,2) NOT NULL,  
    length                           NUMERIC(10,2) NOT NULL,  
    width                            NUMERIC(10,2) NOT NULL,  
    netweight                        NUMERIC(10,2) NOT NULL,
    product_code                     VARCHAR(50)  NOT NULL UNIQUE,
    recommended_freezing_temperature NUMERIC(5,2) NOT NULL,
    product_type_id                  INT NOT NULL,
    seller_id                        INT NOT NULL,
    PRIMARY KEY (id)
);
