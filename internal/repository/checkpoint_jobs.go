package repository

import (
	"sync"
	"time"
)

type CheckpointJob struct {
	ID, RuntimeID, State string
	UpdatedAt            time.Time
}

type CheckpointJobs struct {
	mu   sync.RWMutex
	jobs []CheckpointJob
}

func NewCheckpointJobs() *CheckpointJobs { return &CheckpointJobs{} }

func (s *CheckpointJobs) Begin(id, runtimeID string) {
	s.jobs = append(s.jobs, CheckpointJob{ID: id, RuntimeID: runtimeID, State: "saving", UpdatedAt: time.Now()})
}

func (s *CheckpointJobs) Complete(id string) {
	for i := range s.jobs {
		if s.jobs[i].ID == id {
			s.jobs[i].State = "completed"
			s.jobs[i].UpdatedAt = time.Now()
			return
		}
	}
}

func (s *CheckpointJobs) Get(id string) (CheckpointJob, bool) {
	for _, job := range s.jobs {
		if job.ID == id {
			return job, true
		}
	}
	return CheckpointJob{}, false
}
