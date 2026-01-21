package repositories

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
)

type ClusterEventRepository struct {
	db *sqlx.DB
}

func NewClusterEventRepository(db *sqlx.DB) *ClusterEventRepository {
	return &ClusterEventRepository{db: db}
}

func (r *ClusterEventRepository) LogEvent(event *entities.ClusterEvent) error {
	event.CreatedAt = time.Now()

	_, err := r.db.NamedExec(`
		INSERT INTO cluster_events (cluster_id, pod_id, type, metadata, created_at)
		VALUES (:cluster_id, :pod_id, :type, :metadata, :created_at)
	`, event)

	return err
}

func (r *ClusterEventRepository) ListByCluster(clusterId int64) ([]entities.ClusterEvent, error) {
	var events []entities.ClusterEvent
	err := r.db.Select(&events, `
		SELECT * FROM cluster_events
		WHERE cluster_id = ?
		ORDER BY created_at DESC LIMIT 20
	`, clusterId)
	return events, err
}
