CREATE DATABASE IF NOT EXISTS checkout CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE checkout;

CREATE TABLE vaults (
    id int(10) NOT NULL AUTO_INCREMENT,
    amount BIGINT NOT NULL,
    reference_number char(36) NOT NULL,
    PRIMARY KEY (id)
);