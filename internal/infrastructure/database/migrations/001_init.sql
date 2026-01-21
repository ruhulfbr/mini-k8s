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
-- clusters TABLE
-- Stores clusters belonging to an application
-- =====================================================
CREATE TABLE IF NOT EXISTS clusters
(
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    application_id    INTEGER NOT NULL,                -- FK to applications
    name              TEXT    NOT NULL,                -- Cluster name
    ip                TEXT    NOT NULL,                -- Cluster IP
    port              INTEGER NULL,                    -- Exposed port
    replicas          INTEGER NOT NULL DEFAULT 1,
    cpu               INTEGER NOT NULL,                -- CPU Core
    memory            INTEGER NOT NULL,                -- Memory Limit
    path              TEXT    NOT NULL DEFAULT '/',    -- Cluster path
    type              TEXT    NOT NULL DEFAULT 'http', -- Cluster type (http, worker, etc.)
    deploy_mode       INTEGER NOT NULL DEFAULT 1,      -- Deployment mode (1=image, 2=build)
    image             TEXT             DEFAULT NULL,   -- Docker image for build mode image
    envs              TEXT             DEFAULT NULL,   -- Cluster env variables
    current_image_tag TEXT             DEFAULT NULL,   -- Current Docker image tag
    current_version   TEXT             DEFAULT NULL,   -- Current Deployed version
    status            INTEGER NOT NULL DEFAULT 1,      -- Status (1=active, 0=inactive)
    last_deployed_at  DATETIME         DEFAULT NULL,   -- Last deployed timestamp
    created_at        DATETIME         DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME         DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (application_id) REFERENCES applications (id)
);

CREATE INDEX IF NOT EXISTS idx_clusters_application_id
    ON clusters (application_id);
CREATE INDEX IF NOT EXISTS idx_clusters_type
    ON clusters (type);
CREATE INDEX IF NOT EXISTS idx_clusters_status
    ON clusters (status);


-- =====================================================
-- Cluster build config table
-- =====================================================
CREATE TABLE IF NOT EXISTS cluster_build_configs
(
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    cluster_id          INTEGER NOT NULL,
    git_repo            TEXT    NOT NULL,
    git_branch          TEXT    NOT NULL DEFAULT 'main',
    docker_context_path TEXT    NOT NULL DEFAULT '.',
    dockerfile_name     TEXT    NOT NULL DEFAULT 'Dockerfile',
    created_at          DATETIME         DEFAULT CURRENT_TIMESTAMP,
    updated_at          DATETIME         DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (cluster_id) REFERENCES clusters (id)
);

CREATE INDEX IF NOT EXISTS idx_cluster_build_configs_cluster_id
    ON cluster_build_configs (cluster_id);


-- =====================================================
-- cluster_builds TABLE
-- Tracks cluster_builds per cluster & application
-- =====================================================
CREATE TABLE IF NOT EXISTS cluster_builds
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    cluster_id  INTEGER NOT NULL,      -- FK to clusters
    version     TEXT    NOT NULL,      -- Build/Version
    image_tag   TEXT    NOT NULL,      -- Build/image tag
    deployed_at DATETIME DEFAULT NULL, -- Deployed at timestamp
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (cluster_id) REFERENCES clusters (id)
);

CREATE INDEX IF NOT EXISTS idx_cluster_builds_cluster_id
    ON cluster_builds (cluster_id);


-- =====================================================
-- pods TABLE
-- Stores runtime pod instances per cluster
-- =====================================================
CREATE TABLE IF NOT EXISTS pods
(
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    cluster_id     INTEGER NOT NULL,   -- FK to clusters
    container_id   TEXT    NOT NULL,   -- Container id
    container_name TEXT    NOT NULL,   -- Container name
    ip_address     TEXT    NOT NULL,   -- Container IP Address
    status         INTEGER  DEFAULT 0, -- Pod status 0=pending,1= running
    created_at     DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (cluster_id) REFERENCES clusters (id)
);

CREATE INDEX IF NOT EXISTS idx_pods_cluster_id
    ON pods (cluster_id);
