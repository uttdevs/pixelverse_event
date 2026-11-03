-- DATOS DEMOSTRATIVOS. ESTE SCRIPT ELIMINA LOS REGISTROS ACTUALES DE ESTAS TABLAS.
-- USE `pixelverse_db`;

SET FOREIGN_KEY_CHECKS = 0;
-- TRUNCATE TABLE `inscripciones`;
-- TRUNCATE TABLE `participantes`;
TRUNCATE TABLE `juegos`;
TRUNCATE TABLE `carreras`;
TRUNCATE TABLE `evento_config`;
SET FOREIGN_KEY_CHECKS = 1;

START TRANSACTION;

INSERT INTO `evento_config` (`clave`, `valor`) VALUES
('fecha_evento', '2026-11-27 10:00:00'),
('registros_abiertos', '0');

INSERT INTO `carreras` (`id_carrera`, `nombre`) VALUES
(1, 'Ingeniería en Tecnologías de la Información e Innovación Digital'),
(2, 'Ingeniería Mecatrónica'),
(3, 'Administración'),
(4, 'Contaduría'),
(5, 'Otra');

INSERT INTO `juegos` (`id_juego`, `nombre`, `descripcion`, `costo_inscripcion`, `cupos_maximos`) VALUES
(1, 'King of Fighters', 'Solo los más fuertes sobreviven a la arena, perfecciona tus combos y demuestra tu dominio en combates. ¡Entra en combate y demuestra quién merece el trono!', 20.00, NULL),
(2, 'Street Fighter', '¡Un rival, una batalla y una leyenda! Demuestra tu habilidad frente a los mejores combatientes, cada ronda es una prueba de técnica y concentración.', 20.00, NULL),
(3, 'FIFA 2026', 'El fútbol se juega con talento y se gana con estrategia. ¡Controla el balón y domina la cancha!', 30.00, NULL);

COMMIT;
