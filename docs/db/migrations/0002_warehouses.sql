CREATE TABLE warehouses (
     id INT NOT NULL AUTO_INCREMENT,
     warehouse_code VARCHAR(50) NOT NULL UNIQUE,
     address VARCHAR(100) NOT NULL,
     telephone VARCHAR(100) NOT NULL,
     minimun_capacity INT NOT NULL,
     minimun_temperature FLOAT NOT NULL,
     PRIMARY KEY (id)
);