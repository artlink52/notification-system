package storage

import (
	"context"
	"fmt"

	notificationpb "github.com/artlink52/notification-system/pkg/pb/notification"
	storagepb "github.com/artlink52/notification-system/pkg/pb/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, req *storagepb.SaveNotificationRequest) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx,
		`INSERT INTO notifications (user_id, message, type, status, created_at, retry_count)
		 VALUES ($1, $2, $3, $4, NOW(), 0)
		 RETURNING id`,
		req.UserId, req.Message, req.Type.String(), req.Status.String(),
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert notification: %w", err)
	}
	return id, nil
}

func (r *Repository) Get(ctx context.Context, id int64) (*storagepb.GetNotificationResponse, error) {
	var n storagepb.GetNotificationResponse
	var typStr, statusStr string

	err := r.db.QueryRow(ctx,
		`SELECT id, user_id, message, type, status, created_at::text, retry_count
		 FROM notifications WHERE id = $1`, id,
	).Scan(&n.Id, &n.UserId, &n.Message, &typStr, &statusStr, &n.CreatedAt, &n.RetryCount)
	if err != nil {
		return nil, fmt.Errorf("get notification %d: %w", id, err)
	}

	n.Type = notificationpb.NotificationType(notificationpb.NotificationType_value[typStr])
	n.Status = notificationpb.NotificationStatus(notificationpb.NotificationStatus_value[statusStr])

	return &n, nil
}
