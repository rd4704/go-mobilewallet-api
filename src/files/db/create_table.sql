CREATE DATABASE IF NOT EXISTS mobilewallet;

USE mobilewallet;

-- Users

CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(100) NOT NULL,
    createdAt timestamp DEFAULT current_timestamp,
    modifiedAt timestamp DEFAULT current_timestamp
);

-- insert 3 users
INSERT INTO users (name, email) VALUES ("Rahul", "rahul@myrent.sg");
INSERT INTO users (name, email) VALUES ("Mike", "mike@myrent.sg");
INSERT INTO users (name, email) VALUES ("Pierre", "pierre@myrent.sg");

-- Wallets

CREATE TABLE wallets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    description VARCHAR(100) NOT NULL,
    balance decimal(10,2) NOT NULL,
    userId int NOT NULL,
    createdAt timestamp DEFAULT current_timestamp,
    modifiedAt timestamp DEFAULT current_timestamp,
    FOREIGN KEY (userId) REFERENCES users(id)
);

-- insert 1 wallet for each user

INSERT INTO wallets (description, userId, balance) VALUES ("Rahul's primary wallet", 1, 1000);
INSERT INTO wallets (description, userId, balance) VALUES ("Mike's primary wallet", 2, 1000);
INSERT INTO wallets (description, userId, balance) VALUES ("Pierre's primary wallet", 3, 1000);


-- Transfers

CREATE TABLE transfers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    description VARCHAR(100) NOT NULL,
    fromWallet int NOT NULL,
    toWallet int NOT NULL,
    Amount decimal(10,2) NOT NULL,
    createdAt timestamp DEFAULT current_timestamp,
    modifiedAt timestamp DEFAULT current_timestamp,
    FOREIGN KEY (fromWallet) REFERENCES wallets(id),
    FOREIGN KEY (toWallet) REFERENCES wallets(id)
);
