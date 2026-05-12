package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/artlink52/notification-system/internal/gateway/queue"
	"github.com/artlink52/notification-system/internal/models"
)

type Handler struct {
	queue *queue.Queue
}

func New(q *queue.Queue) *Handler {
	return &Handler{queue: q}
}

type sendRequest struct {
	UserID  int64  `json:"user_id"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type sendResponse struct {
	Status string `json:"status"`
}

func (h *Handler) SendNotification(w http.ResponseWriter, r *http.Request) {
	var req sendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 || req.Message == "" {
		writeError(w, "user_id and message are required", http.StatusBadRequest)
		return
	}

	notifType := models.NotificationType(req.Type)
	if notifType != models.TypeEmail && notifType != models.TypeSMS {
		writeError(w, "type must be 'email' or 'sms'", http.StatusBadRequest)
		return
	}

	err := h.queue.Push(models.Task{
		UserID:  req.UserID,
		Message: req.Message,
		Type:    notifType,
	})
	if errors.Is(err, queue.ErrFull) {
		writeError(w, "server is busy, try again later", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(sendResponse{Status: "queued"})
}

func writeError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
