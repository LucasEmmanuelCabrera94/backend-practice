CREATE DATABASE IF NOT EXISTS backend;

CREATE USER IF NOT EXISTS 'backend'@'localhost' IDENTIFIED WITH mysql_native_password BY 'backendpw';
CREATE USER IF NOT EXISTS 'backend'@'%' IDENTIFIED WITH mysql_native_password BY 'backendpw';

GRANT ALL PRIVILEGES ON backend.* TO 'backend'@'localhost';
GRANT ALL PRIVILEGES ON backend.* TO 'backend'@'%';

FLUSH PRIVILEGES;

USE backend;

CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  passwordhash VARCHAR(255) NOT NULL
);
