-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS gateway_services (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    base_url VARCHAR(500) NOT NULL,
    protocol VARCHAR(20) NOT NULL DEFAULT 'http',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    rate_limit_per_minute INT UNSIGNED NULL,
    health_status VARCHAR(20) NOT NULL DEFAULT 'unknown',
    health_checked_at DATETIME(3) NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    INDEX idx_gateway_services_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS gateway_routes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    service_id BIGINT UNSIGNED NOT NULL,
    method VARCHAR(10) NOT NULL,
    path_pattern VARCHAR(500) NOT NULL,
    permission_match_mode VARCHAR(10) NOT NULL DEFAULT 'any',
    rate_limit_per_minute INT UNSIGNED NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    deleted_at DATETIME(3) NULL,
    INDEX idx_gateway_routes_service_id (service_id),
    INDEX idx_gateway_routes_deleted_at (deleted_at),
    UNIQUE KEY uq_gateway_routes_service_method_path (service_id, method, path_pattern),
    FOREIGN KEY (service_id) REFERENCES gateway_services(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS gateway_route_permissions (
    route_id BIGINT UNSIGNED NOT NULL,
    permission_id BIGINT UNSIGNED NOT NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (route_id, permission_id),
    FOREIGN KEY (route_id) REFERENCES gateway_routes(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS gateway_audit_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    entity_type VARCHAR(20) NOT NULL,
    entity_id BIGINT UNSIGNED NOT NULL,
    action VARCHAR(20) NOT NULL,
    actor_user_id BIGINT UNSIGNED NOT NULL,
    changes JSON NULL,
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    INDEX idx_gateway_audit_logs_entity (entity_type, entity_id),
    INDEX idx_gateway_audit_logs_actor (actor_user_id),
    FOREIGN KEY (actor_user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS gateway_audit_logs;
DROP TABLE IF EXISTS gateway_route_permissions;
DROP TABLE IF EXISTS gateway_routes;
DROP TABLE IF EXISTS gateway_services;
-- +goose StatementEnd
