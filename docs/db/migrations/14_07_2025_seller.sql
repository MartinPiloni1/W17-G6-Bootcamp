CREATE TABLE sellers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cid INT NOT NULL UNIQUE,
    company_name VARCHAR(100) NOT NULL,
    address VARCHAR(150) NOT NULL,
    telephone VARCHAR(25) NOT NULL,
    locality_id VARCHAR(20) NOT NULL,
    FOREIGN KEY (locality_id) REFERENCES localities(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);