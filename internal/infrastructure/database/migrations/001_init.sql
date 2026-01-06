-- =====================================================
-- APPLICATIONS TABLE
-- Stores high-level application metadata
-- =====================================================
CREATE TABLE IF NOT EXISTS applications
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL UNIQUE, -- Application name (unique)
    git_repo   TEXT NOT NULL,        -- Git repository URL
    git_branch TEXT NOT NULL,        -- Git repository branch
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


-- =====================================================
-- SERVICES TABLE
-- Stores services belonging to an application
-- =====================================================
CREATE TABLE IF NOT EXISTS services
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    application_id INTEGER NOT NULL,                -- FK to applications
    name           TEXT    NOT NULL,                -- Service name
    ip             TEXT    NOT NULL,                -- Service IP
    port           INTEGER NULL,                    -- Exposed port
    image_tag      TEXT             DEFAULT NULL,   -- Docker image tag
    context_path   TEXT    NOT NULL,                -- Build context path
    replicas       INTEGER NOT NULL DEFAULT 1,
    resources      TEXT    NOT NULL,                -- Resource limits (JSON)
    path           TEXT    NOT NULL DEFAULT '/',    -- Service path
    type           TEXT    NOT NULL DEFAULT 'http', -- Service type (http, worker, etc.)
    status         INTEGER NOT NULL DEFAULT 1,      -- Status (1=active, 0 inactive)
    last_build_at  DATETIME         DEFAULT NULL,   -- Last build timestamp
    created_at     DATETIME         DEFAULT CURRENT_TIMESTAMP,
    updated_at     DATETIME         DEFAULT CURRENT_TIMESTAMP,

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
-- BUILD HISTORY TABLE
-- Tracks build history per service & application
-- =====================================================
CREATE TABLE IF NOT EXISTS build_history
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    application_id INTEGER NOT NULL, -- FK to applications
    service_id     INTEGER NOT NULL, -- FK to services
    tag            TEXT    NOT NULL, -- Build/image tag
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (application_id) REFERENCES applications (id),
    FOREIGN KEY (service_id) REFERENCES services (id)
);

-- Index for faster application-based queries
CREATE INDEX IF NOT EXISTS idx_build_history_application_id
    ON build_history (application_id);
CREATE INDEX IF NOT EXISTS idx_build_history_service_id
    ON build_history (service_id);


-- =====================================================
-- PODS TABLE
-- Stores runtime pod instances per service
-- =====================================================
CREATE TABLE IF NOT EXISTS pods
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    application_id INTEGER NOT NULL,           -- FK to applications
    service_id     INTEGER NOT NULL,           -- FK to services
    name           TEXT    NOT NULL,           -- Pod name
    status         INTEGER NOT NULL DEFAULT 0, -- Status (1=running, 0=pending)
    created_at     DATETIME         DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (application_id) REFERENCES applications (id),
    FOREIGN KEY (service_id) REFERENCES services (id)
);

-- Indexes for efficient pod lookups
CREATE INDEX IF NOT EXISTS idx_pods_application_id
    ON pods (application_id);

CREATE INDEX IF NOT EXISTS idx_pods_service_id
    ON pods (service_id);

CREATE INDEX IF NOT EXISTS idx_pods_status
    ON pods (status);
