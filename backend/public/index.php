<?php

header('Content-Type: application/json; charset=utf-8');

// Autoload composer if available
if (file_exists(__DIR__ . '/../vendor/autoload.php')) {
    require_once __DIR__ . '/../vendor/autoload.php';
}

require_once __DIR__ . '/../config/database.php';

$method = $_SERVER['REQUEST_METHOD'];
$uri = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH);

// Basic health check endpoint
if ($uri === '/api/health' || $uri === '/' || $uri === '/index.php') {
    echo json_encode([
        'status' => 'ok',
        'message' => 'PixelVerse Event API is running'
    ]);
    exit;
}

// 404 handler
http_response_code(404);
echo json_encode([
    'status' => 'error',
    'message' => 'Endpoint not found'
]);
