package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"syscall"

	"github.com/artlink52/notification-system/internal/notification"
	pb "github.com/artlink52/notification-system/pkg/pb/notification"
	storagepb "github.com/artlink52/notification-system/pkg/pb/storage"
	"google.golang.org/grpc"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	storageConn, err := grpc.NewClient("localhost:50052")
	if err != nil {
		log.Fatal(err)
	}
	defer storageConn.Close()

	storageClient := storagepb.NewStorageServiceClient(storageConn)

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, notification.New(storageClient))

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("notification-service listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
	log.Println("notification-service stopped")
}
