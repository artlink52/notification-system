PROTO_DIR = api/proto
MODULE    = github.com/artlink52/notification-system

.PHONY: proto proto-notification proto-storage

proto: proto-notification proto-storage

proto-notification:
	mkdir -p pkg/pb/notification
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=pkg/pb/notification \
		--go-grpc_out=pkg/pb/notification \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		notification.proto

proto-storage:
	mkdir -p pkg/pb/storage
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=pkg/pb/storage \
		--go-grpc_out=pkg/pb/storage \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--go_opt=Mnotification.proto=$(MODULE)/pkg/pb/notification \
		--go-grpc_opt=Mnotification.proto=$(MODULE)/pkg/pb/notification \
		storage.proto