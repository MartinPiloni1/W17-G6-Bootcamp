CREATE TABLE buyers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    card_number_id INT NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL
);