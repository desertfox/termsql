CREATE DATABASE IF NOT EXISTS test_db;
USE test_db;

CREATE TABLE IF NOT EXISTS test_table (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE USER 'test_user'@'%' IDENTIFIED WITH mysql_native_password BY 'test_password';
GRANT ALL PRIVILEGES ON test_db.* TO 'test_user'@'%';

CREATE USER 'test_user_ssl'@'%' IDENTIFIED BY 'test_password';
GRANT ALL PRIVILEGES ON test_db.* TO 'test_user_ssl'@'%';
ALTER USER 'test_user_ssl'@'%' REQUIRE SSL;

INSERT INTO test_table (name, email) VALUES ('John Doe', 'john.doe@example.com');