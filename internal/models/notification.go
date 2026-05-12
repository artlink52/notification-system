package models

type NotificationType string
type NotificationStatus string

const (
	TypeEmail NotificationType = "email"
	TypeSMS   NotificationType = "sms"
)

const (
	StatusPending NotificationStatus = "pending"
	StatusSent    NotificationStatus = "sent"
	StatusFailed  NotificationStatus = "failed"
)

type Notification struct {
	ID         int64
	UserID     int64
	Message    string
	Type       NotificationType
	Status     NotificationStatus
	CreatedAt  string
	RetryCount int32
}

// Task — единица работы, которая кладётся в очередь.
type Task struct {
	UserID  int64
	Message string
	Type    NotificationType
}
