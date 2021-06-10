CREATE TABLE users (
    id char(36) PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password CHAR(60) NOT NULL,
    email VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    UNIQUE KEY unique_username (username)
);
