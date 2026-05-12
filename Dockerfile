FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG SERVICE
RUN CGO_ENABLED=0 go build -o /bin/service ./cmd/${SERVICE}

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/service /bin/service
CMD ["/bin/service"]