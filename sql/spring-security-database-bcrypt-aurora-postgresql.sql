-- Drop tables if they exist
DROP TABLE IF EXISTS authorities;
DROP TABLE IF EXISTS users;

-- Create the "users" table
CREATE TABLE users (
    username VARCHAR(50) NOT NULL,
    password CHAR(68) NOT NULL,
    enabled BOOLEAN NOT NULL,
    PRIMARY KEY (username)
);

-- Create the "authorities" table without a foreign key constraint
CREATE TABLE authorities (
    username VARCHAR(50) NOT NULL,
    authority VARCHAR(50) NOT NULL,
    UNIQUE (username, authority)
);