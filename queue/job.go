package queue

import (
	"go-pipelines/config"
	"time"
)

type Job struct {
	ID         string
	Config     *config.Config
	ReceivedAt time.Time
}
