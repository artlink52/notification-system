package storage

import (
	"context"
	"fmt"

	storagepb "github.com/artlink52/notification-system/pkg/pb/storage"
)

type Service struct {
	storagepb.UnimplementedStorageServiceServer
	repo *Repository
}

func New(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) SaveNotification(ctx context.Context, req *storagepb.SaveNotificationRequest) (*storagepb.SaveNotificationResponse, error) {
	id, err := s.repo.Save(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("save notification: %w", err)
	}
	return &storagepb.SaveNotificationResponse{Id: id}, nil
}

func (s *Service) GetNotification(ctx context.Context, req *storagepb.GetNotificationRequest) (*storagepb.GetNotificationResponse, error) {
	n, err := s.repo.Get(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("get notification: %w", err)
	}
	return n, nil
}
