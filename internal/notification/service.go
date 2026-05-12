package notification

import (
	"context"
	"fmt"
	"log"

	pb "github.com/artlink52/notification-system/pkg/pb/notification"
	storagepb "github.com/artlink52/notification-system/pkg/pb/storage"
)

type Service struct {
	pb.UnimplementedNotificationServiceServer
	storageClient storagepb.StorageServiceClient
}

func New(storageClient storagepb.StorageServiceClient) *Service {
	return &Service{storageClient: storageClient}
}

func (s *Service) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	log.Printf("[notification] sending %s to user_id=%d", req.Type, req.UserId)

	// эмуляция отправки
	status := pb.NotificationStatus_SENT
	if err := emulate(req.Type); err != nil {
		log.Printf("[notification] failed to send: %v", err)
		status = pb.NotificationStatus_FAILED
	}

	saved, err := s.storageClient.SaveNotification(ctx, &storagepb.SaveNotificationRequest{
		UserId:  req.UserId,
		Message: req.Message,
		Type:    req.Type,
		Status:  status,
	})
	if err != nil {
		return nil, fmt.Errorf("save notification: %w", err)
	}

	return &pb.SendNotificationResponse{
		NotificationId: saved.Id,
		Status:         status,
	}, nil
}

func emulate(t pb.NotificationType) error {
	switch t {
	case pb.NotificationType_EMAIL:
		log.Println("[notification] email sent")
	case pb.NotificationType_SMS:
		log.Println("[notification] sms sent")
	}
	return nil
}
