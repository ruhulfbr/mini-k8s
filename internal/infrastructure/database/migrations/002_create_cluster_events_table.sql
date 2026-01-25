CREATE TABLE IF NOT EXISTS cluster_events (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    cluster_id  INTEGER NOT NULL,
    pod_id      INTEGER NULL,
    type        TEXT NOT NULL,
    metadata    TEXT,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cluster_events_cluster_id
    ON cluster_events(cluster_id);