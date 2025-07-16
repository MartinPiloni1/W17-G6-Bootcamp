CREATE TABLE localities (
    id VARCHAR(20) PRIMARY KEY,
    locality_name VARCHAR(100) NOT NULL,
    province_name VARCHAR(100) NOT NULL,
    country_name VARCHAR(100) NOT NULL
);

ALTER TABLE sellers
    ADD CONSTRAINT fk_sellers_locality
        FOREIGN KEY (locality_id)
            REFERENCES localities(id)
            ON UPDATE CASCADE
            ON DELETE RESTRICT;