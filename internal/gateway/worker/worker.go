package worker

import (
	"context"
	"log"
	"time"

	"github.com/artlink52/notification-system/internal/gateway/queue"
	"github.com/artlink52/notification-system/internal/models"
	pb "github.com/artlink52/notification-system/pkg/pb/notification"
)

type Pool struct {
	queue  *queue.Queue
	client pb.NotificationServiceClient
	size   int
}

func NewPool(q *queue.Queue, client pb.NotificationServiceClient, size int) *Pool {
	return &Pool{queue: q, client: client, size: size}
}

// Run запускает size горутин-воркеров. Блокирует до отмены ctx.
func (p *Pool) Run(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	for i := range p.size {
		go p.run(ctx, i)
	}
}

func (p *Pool) run(ctx context.Context, id int) {
	log.Printf("[worker %d] started", id)
	for {
		select {
		case <-ctx.Done():
			log.Printf("[worker %d] stopped", id)
			return
		case task := <-p.queue.Chan():
			p.process(ctx, task)
		}
	}
}

func (p *Pool) process(ctx context.Context, task models.Task) {
	resp, err := p.client.SendNotification(ctx, &pb.SendNotificationRequest{
		UserId:  task.UserID,
		Message: task.Message,
		Type:    toProtoType(task.Type),
	})
	if err != nil {
		log.Printf("[worker] send error: user_id=%d err=%v", task.UserID, err)
		return
	}
	log.Printf("[worker] sent: notification_id=%d status=%v", resp.NotificationId, resp.Status)
}

func toProtoType(t models.NotificationType) pb.NotificationType {
	if t == models.TypeSMS {
		return pb.NotificationType_SMS
	}
	return pb.NotificationType_EMAIL
}
