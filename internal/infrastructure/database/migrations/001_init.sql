-- =====================================================
-- APPLICATIONS TABLE
-- Stores high-level application metadata
-- =====================================================
CREATE TABLE IF NOT EXISTS applications
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL UNIQUE,
    description TEXT             DEFAULT NULL,
    status      INTEGER NOT NULL DEFAULT 1, -- Status (1=active, 0=inactive)
    created_at  DATETIME         DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME         DEFAULT CURRENT_TIMESTAMP
);

-- =====================================================
-- SERVICES TABLE
-- Stores services belonging to an application
-- =====================================================
CREATE TABLE IF NOT EXISTS services
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    application_id    INTEGER NOT NULL,                -- FK to applications
    name              TEXT    NOT NULL,                -- Service name
    ip                TEXT    NOT NULL,                -- Service IP
    port              INTEGER NULL,                    -- Exposed port
    replicas          INTEGER NOT NULL DEFAULT 1,
    cpu               INTEGER NOT NULL,                -- CPU Core
    memory            INTEGER NOT NULL,                -- Memory Limit
    path              TEXT    NOT NULL DEFAULT '/',    -- Service path
    type              TEXT    NOT NULL DEFAULT 'http', -- Service type (http, worker, etc.)
    deploy_mode       INTEGER NOT NULL DEFAULT 1,      -- Deployment mode (1=image, 2=build)
    current_image_tag TEXT             DEFAULT NULL,   -- Current Docker image tag
    current_version   TEXT             DEFAULT NULL,   -- Current Deployed version
    status            INTEGER NOT NULL DEFAULT 1,      -- Status (1=active, 0=inactive)
    last_deployed_at  DATETIME         DEFAULT NULL,   -- Last deployed timestamp
    created_at        DATETIME         DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME         DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (application_id) REFERENCES applications (id)
);

-- Index for faster lookup by application
CREATE INDEX IF NOT EXISTS idx_services_application_id
    ON services (application_id);
CREATE INDEX IF NOT EXISTS idx_services_type
    ON services (type);
CREATE INDEX IF NOT EXISTS idx_services_status
    ON services (status);


-- =====================================================
-- Service build config table
-- =====================================================
CREATE TABLE IF NOT EXISTS service_build_configs
(
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id          INTEGER NOT NULL,
    git_repo            TEXT    NOT NULL,
    git_branch          TEXT    NOT NULL DEFAULT 'main',
    docker_context_path TEXT    NOT NULL DEFAULT '.',
    dockerfile_name     TEXT    NOT NULL DEFAULT 'Dockerfile',
    created_at          DATETIME         DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME         DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (service_id) REFERENCES services (id)
);

CREATE INDEX IF NOT EXISTS idx_service_build_configs_service_id
    ON service_build_configs (service_id);


-- =====================================================
-- BUILD HISTORY TABLE
-- Tracks build history per service & application
-- =====================================================
CREATE TABLE IF NOT EXISTS build_history
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id     INTEGER NOT NULL,      -- FK to services
    version        TEXT    NOT NULL,      -- Build/Version
    image_tag      TEXT    NOT NULL,      -- Build/image tag
    deployed_at    DATETIME DEFAULT NULL, -- Deployed at timestamp
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (service_id) REFERENCES services (id)
);

-- Index for faster application-based queries
CREATE INDEX IF NOT EXISTS idx_build_history_service_id
    ON build_history (service_id);


-- =====================================================
-- PODS TABLE
-- Stores runtime pod instances per service
-- =====================================================
CREATE TABLE IF NOT EXISTS pods
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    service_id INTEGER NOT NULL,           -- FK to services
    name       TEXT    NOT NULL,           -- Pod name
    status     INTEGER NOT NULL DEFAULT 0, -- Status (1=running, 0=pending)
    created_at DATETIME         DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (service_id) REFERENCES services (id)
);

CREATE INDEX IF NOT EXISTS idx_pods_service_id
    ON pods (service_id);

CREATE INDEX IF NOT EXISTS idx_pods_status
    ON pods (status);
