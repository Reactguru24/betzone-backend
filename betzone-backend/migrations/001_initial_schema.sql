-- Betzone Database Setup Script
-- Run this script to create the database and initial schema

-- Create database
CREATE DATABASE IF NOT EXISTS betzone;
USE betzone;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(50) PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    balance DECIMAL(10, 2) DEFAULT 0,
    currency VARCHAR(10) DEFAULT 'KES',
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_phone (phone),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create indexes for better performance
CREATE INDEX idx_updated_at ON users(updated_at);

-- Grant privileges to betzone user (if not exists)
-- Uncomment and modify as needed:
-- CREATE USER IF NOT EXISTS 'betzone'@'localhost' IDENTIFIED BY 'your_password_here';
-- GRANT ALL PRIVILEGES ON betzone.* TO 'betzone'@'localhost';
-- FLUSH PRIVILEGES;
