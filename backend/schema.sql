-- Database schema for PixelVerse Event

CREATE DATABASE IF NOT EXISTS `pixelverse_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `pixelverse_db`;

CREATE TABLE IF NOT EXISTS `attendees` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `full_name` VARCHAR(150) NOT NULL,
    `email` VARCHAR(150) NOT NULL UNIQUE,
    `phone` VARCHAR(30) NULL,
    `ticket_type` ENUM('general', 'vip', 'speaker') DEFAULT 'general',
    `registered_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
