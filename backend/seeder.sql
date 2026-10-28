-- DATOS DEMOSTRATIVOS. ESTE SCRIPT ELIMINA LOS REGISTROS ACTUALES DE ESTAS TABLAS.
-- USE `pixelverse_db`;

SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE `inscripciones`;
TRUNCATE TABLE `participantes`;
TRUNCATE TABLE `juegos`;
TRUNCATE TABLE `carreras`;
TRUNCATE TABLE `evento_config`;
SET FOREIGN_KEY_CHECKS = 1;

START TRANSACTION;

INSERT INTO `evento_config` (`id`, `fecha_evento`, `registros_abiertos`) VALUES
(1, '2026-11-27 10:00:00', 0);

INSERT INTO `carreras` (`id_carrera`, `nombre`) VALUES
(1, 'Ingeniería en Tecnologías de la Información e Innovación Digital'),
(2, 'Ingeniería Mecatrónica'),
(3, 'Administración'),
(4, 'Contaduría'),
(5, 'Otra');

INSERT INTO `juegos` (`id_juego`, `nombre`) VALUES
(1, 'King of Fighters'),
(2, 'Street Fighter'),
(3, 'FIFA 1 vs 1'),
(4, 'FIFA 2 vs 2');

INSERT INTO `participantes` (`id_participante`, `nombre_completo`, `matricula`, `id_carrera`, `cuatrimestre`, `grupo`, `telefono`, `correo`, `fecha_registro`) VALUES
(1, 'Ana Torres Ruiz', 'EJ-0001', 1, 4, 'B', '2220000001', 'ana.ejemplo@correo.com', '2026-10-19 18:06:22'),
(2, 'Luis Hernández Pérez', 'EJ-0002', 1, 4, 'B', '2220000002', 'luis.ejemplo@correo.com', '2026-10-19 18:06:22'),
(3, 'María García López', 'EJ-0003', 2, 5, 'A', '2220000003', NULL, '2026-10-19 18:06:22'),
(4, 'Carlos Ramírez Soto', 'EJ-0004', 3, 3, 'C', '2220000004', 'carlos.ejemplo@correo.com', '2026-10-19 18:06:22'),
(5, 'Sofía Martínez Cruz', 'EJ-0005', 1, 4, 'A', '2220000005', NULL, '2026-10-19 18:06:22'),
(6, 'Diego Flores Vega', 'EJ-0006', 4, 2, 'B', '2220000006', 'diego.ejemplo@correo.com', '2026-10-19 18:06:22'),
(7, 'Valeria Morales Díaz', 'EJ-0007', 2, 5, 'B', '2220000007', NULL, '2026-10-19 18:06:22'),
(8, 'Jorge Castillo Rojas', 'EJ-0008', 5, 3, 'A', '2220000008', 'jorge.ejemplo@correo.com', '2026-10-19 18:06:22'),
(9, 'EDUARDO ADOLFO ROSAS ANDRADE', '49487894879', 1, 4, 'B', '2491000044', NULL, '2026-10-19 18:18:08');

INSERT INTO `inscripciones` (`id_inscripcion`, `id_participante`, `id_juego`, `fecha_inscripcion`, `asistio`) VALUES
(1, 1, 1, '2026-10-19 12:06:22', 0),
(2, 8, 1, '2026-10-19 12:06:22', 0),
(3, 2, 2, '2026-10-19 12:06:22', 0),
(4, 4, 2, '2026-10-19 12:06:22', 0),
(5, 6, 3, '2026-10-19 12:06:22', 0),
(6, 2, 4, '2026-10-19 12:06:22', 0),
(7, 6, 4, '2026-10-19 12:06:22', 0),
(14, 9, 1, '2026-10-19 12:18:08', 0),
(15, 9, 2, '2026-10-20 14:49:44', 0),
(16, 9, 3, '2026-10-20 14:49:44', 0);

COMMIT;
