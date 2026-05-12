package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/artlink52/notification-system/internal/gateway/handler"
	"github.com/artlink52/notification-system/internal/gateway/queue"
	"github.com/artlink52/notification-system/internal/gateway/worker"
	pb "github.com/artlink52/notification-system/pkg/pb/notification"
	"google.golang.org/grpc"
)

func notificationAddr() string {
	if addr := os.Getenv("NOTIFICATION_SERVICE_ADDR"); addr != "" {
		return addr
	}
	return "localhost:50051"
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	conn, err := grpc.NewClient(notificationAddr())
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewNotificationServiceClient(conn)

	q := queue.New(100)
	pool := worker.NewPool(q, client, 4)

	pool.Run(ctx)

	h := handler.New(q)
	http.HandleFunc("POST /send", h.SendNotification)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
