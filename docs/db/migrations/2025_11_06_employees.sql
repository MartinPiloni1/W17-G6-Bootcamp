CREATE TABLE `employees` (
     `id` INT NOT NULL AUTO_INCREMENT,
     `card_number_id` VARCHAR(50) NOT NULL,
     `first_name` VARCHAR(100) NOT NULL,
     `last_name` VARCHAR(100) NOT NULL,
     `warehouse_id` INT NOT NULL,
     PRIMARY KEY (`id`)
);