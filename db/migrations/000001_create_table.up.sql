CREATE TABLE users (
     id SERIAL PRIMARY KEY,
     username VARCHAR(50) UNIQUE NOT NULL,
     email VARCHAR(100) UNIQUE NOT NULL,
     password_hash VARCHAR(100) NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);