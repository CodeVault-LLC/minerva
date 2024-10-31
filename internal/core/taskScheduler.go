package core

import (
	"sync"

	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
)

// TaskScheduler manages job queueing and dispatching
type TaskScheduler struct {
	queue      []*entities.Job // Task queue, can implement a priority queue
	workerPool chan struct{}   // Manages number of concurrent workers
	mu         sync.Mutex
}

var Scheduler *TaskScheduler

// NewTaskScheduler initializes TaskScheduler
func NewTaskScheduler(workerCount int) *TaskScheduler {
	return &TaskScheduler{
		queue:      []*entities.Job{},
		workerPool: make(chan struct{}, workerCount),
	}
}

// AddJob adds a new task to the queue
func (s *TaskScheduler) AddJob(job *entities.Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue = append(s.queue, job)
}

// DispatchTasks processes queued tasks
func (s *TaskScheduler) DispatchTasks(inspector *Inspector) {
	for _, job := range s.queue {
		s.workerPool <- struct{}{}
		go func(j *entities.Job) {
			defer func() { <-s.workerPool }()
			if err := inspector.Execute(j); err != nil {
				j.Status = entities.Failed
			} else {
				j.Status = entities.Completed
			}
			s.updateJobStatus(j)
		}(job)
	}
}

// processJob processes individual tasks and updates DataStore
func (s *TaskScheduler) processJob(job *entities.Job, inspector *Inspector) {
	job.Status = entities.Processing
	inspector.Execute(job) // Call relevant module based on Job.Type
	job.Status = entities.Completed
	// Save results to DataStore via repository
}

// updateJobStatus updates the status in DataStore
func (s *TaskScheduler) updateJobStatus(job *entities.Job) {
	// Logic to update job status in database

}
