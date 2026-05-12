package queue

import (
	"errors"

	"github.com/artlink52/notification-system/internal/models"
)

var ErrFull = errors.New("queue is full")

type Queue struct {
	ch chan models.Task
}

func New(bufSize int) *Queue {
	return &Queue{ch: make(chan models.Task, bufSize)}
}

// Push кладёт задачу в очередь. Если очередь заполнена — возвращает ErrFull.
func (q *Queue) Push(task models.Task) error {
	select {
	case q.ch <- task:
		return nil
	default:
		return ErrFull
	}
}

func (q *Queue) Chan() <-chan models.Task {
	return q.ch
}
