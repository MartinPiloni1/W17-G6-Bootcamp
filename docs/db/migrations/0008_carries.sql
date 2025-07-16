CREATE TABLE carries (
    `id` int PRIMARY KEY AUTO_INCREMENT,
    `cid` VARCHAR(20) NOT NULL UNIQUE,
    `company_name` VARCHAR(100) NOT NULL,
    `address` VARCHAR(100) NOT NULL,
    `telephone` VARCHAR(100) NOT NULL,
    `locality_id` VARCHAR(20) NOT NULL,
   FOREIGN KEY (locality_id) REFERENCES localities(id)
       ON UPDATE CASCADE
       ON DELETE RESTRICT
);