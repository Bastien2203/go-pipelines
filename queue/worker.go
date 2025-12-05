package queue

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	id         string
	jobTimeout int
	Queue      chan Job
	processJob func(ctx context.Context, j Job)
}

func NewWorker(jobTimeout, queueSize int, processJob func(ctx context.Context, j Job)) *Worker {
	return &Worker{
		id:         uuid.New().String(),
		jobTimeout: jobTimeout,
		Queue:      make(chan Job, queueSize),
		processJob: processJob,
	}
}

func (w *Worker) Start() {
	log.Printf("Worker %s started processing jobs", w.id)
	for job := range w.Queue {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(w.jobTimeout)*time.Minute)
		w.processJob(ctx, job)
		cancel()
	}
}
