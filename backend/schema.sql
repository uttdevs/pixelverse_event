-- Esquema de base de datos para PixelVerse Event - Torneo de videojuegos
CREATE DATABASE IF NOT EXISTS `pixelverse_db` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE `pixelverse_db`;

CREATE TABLE IF NOT EXISTS `carreras` (
    `id_carrera` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `nombre` VARCHAR(120) NOT NULL,
    PRIMARY KEY (`id_carrera`),
    UNIQUE KEY `uq_carreras_nombre` (`nombre`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `juegos` (
    `id_juego` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `nombre` VARCHAR(60) NOT NULL,
    PRIMARY KEY (`id_juego`),
    UNIQUE KEY `uq_juegos_nombre` (`nombre`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `participantes` (
    `id_participante` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `nombre` VARCHAR(60) NOT NULL,
    `apellidos` VARCHAR(80) NOT NULL,
    `matricula` VARCHAR(20) NOT NULL,
    `id_carrera` INT(10) UNSIGNED NOT NULL,
    `cuatrimestre` TINYINT(3) UNSIGNED NOT NULL,
    `grupo` VARCHAR(5) NOT NULL,
    `telefono` VARCHAR(15) NOT NULL,
    `correo` VARCHAR(100) DEFAULT NULL,
    `fecha_registro` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id_participante`),
    UNIQUE KEY `uq_participantes_matricula` (`matricula`),
    KEY `idx_participantes_carrera` (`id_carrera`),
    CONSTRAINT `fk_participantes_carrera` FOREIGN KEY (`id_carrera`) REFERENCES `carreras` (`id_carrera`) ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `inscripciones` (
    `id_inscripcion` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
    `id_participante` INT(10) UNSIGNED NOT NULL,
    `id_juego` INT(10) UNSIGNED NOT NULL,
    `fecha_inscripcion` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `asistio` TINYINT(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (`id_inscripcion`),
    UNIQUE KEY `uq_inscripcion_persona_juego` (`id_participante`, `id_juego`),
    KEY `idx_inscripciones_juego` (`id_juego`),
    CONSTRAINT `fk_inscripciones_juego` FOREIGN KEY (`id_juego`) REFERENCES `juegos` (`id_juego`) ON UPDATE CASCADE,
    CONSTRAINT `fk_inscripciones_participante` FOREIGN KEY (`id_participante`) REFERENCES `participantes` (`id_participante`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE OR REPLACE VIEW `vista_inscripciones` AS
SELECT
    `i`.`id_inscripcion` AS `id_inscripcion`,
    `j`.`nombre` AS `juego`,
    CONCAT(`p`.`nombre`, ' ', `p`.`apellidos`) AS `participante`,
    `p`.`matricula` AS `matricula`,
    `c`.`nombre` AS `carrera`,
    `p`.`cuatrimestre` AS `cuatrimestre`,
    `p`.`grupo` AS `grupo`,
    `p`.`telefono` AS `telefono`,
    `p`.`correo` AS `correo`,
    `i`.`asistio` AS `asistio`,
    `i`.`fecha_inscripcion` AS `fecha_inscripcion`
FROM `inscripciones` `i`
JOIN `participantes` `p` ON `p`.`id_participante` = `i`.`id_participante`
JOIN `juegos` `j` ON `j`.`id_juego` = `i`.`id_juego`
JOIN `carreras` `c` ON `c`.`id_carrera` = `p`.`id_carrera`;
