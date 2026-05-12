package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/artlink52/notification-system/internal/storage"
	storagepb "github.com/artlink52/notification-system/pkg/pb/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func dsn() string {
	if v := os.Getenv("DATABASE_URL"); v != "" {
		return v
	}
	return "postgres://postgres:postgres@localhost:5432/notifications?sslmode=disable"
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := pgxpool.New(ctx, dsn())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatal("postgres unavailable: ", err)
	}

	repo := storage.NewRepository(db)
	svc := storage.New(repo)

	grpcServer := grpc.NewServer()
	storagepb.RegisterStorageServiceServer(grpcServer, svc)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Println("storage-service listening on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	grpcServer.GracefulStop()
	log.Println("storage-service stopped")
}
