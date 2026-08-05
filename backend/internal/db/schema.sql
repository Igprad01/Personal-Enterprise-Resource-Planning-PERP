-- Schema for the Personal ERP Daily Activity module.
-- Run automatically by the server on startup (idempotent).

CREATE TABLE IF NOT EXISTS categories (
    id         BIGINT       NOT NULL AUTO_INCREMENT,
    name       VARCHAR(100) NOT NULL,
    color      VARCHAR(7)   NOT NULL DEFAULT '#2563eb',
    icon       VARCHAR(32)  NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uq_categories_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS activities (
    id            BIGINT       NOT NULL AUTO_INCREMENT,
    title         VARCHAR(255) NOT NULL,
    description   TEXT         NULL,
    category_id   BIGINT       NULL,
    activity_date DATE         NOT NULL,
    start_time    TIME         NULL,
    end_time      TIME         NULL,
    duration_min  INT          NULL,
    status        ENUM('planned','in_progress','done','cancelled') NOT NULL DEFAULT 'planned',
    priority      ENUM('low','medium','high') NOT NULL DEFAULT 'medium',
    notes         TEXT         NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    KEY idx_activities_date (activity_date),
    KEY idx_activities_status (status),
    KEY idx_activities_category (category_id),
    CONSTRAINT fk_activities_category FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT IGNORE INTO categories (name, color, icon) VALUES
    ('Work', '#2563eb', 'briefcase'),
    ('Personal', '#16a34a', 'user'),
    ('Health', '#dc2626', 'heart'),
    ('Learning', '#9333ea', 'book'),
    ('Errands', '#d97706', 'shopping-cart');
