package core

import (
	"sync"
	"time"

	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"go.uber.org/zap"
)

// TaskScheduler manages job queueing and dispatching
type TaskScheduler struct {
	queue      []*entities.JobModel // Task queue, can implement a priority queue if needed
	workerPool chan struct{}        // Manages number of concurrent workers
	mu         sync.Mutex
}

var Scheduler *TaskScheduler

// NewTaskScheduler initializes TaskScheduler
func NewTaskScheduler(workerCount int) *TaskScheduler {
	return &TaskScheduler{
		queue:      []*entities.JobModel{},
		workerPool: make(chan struct{}, workerCount),
	}
}

// AddJob adds a new task to the queue
func (s *TaskScheduler) AddJob(job *entities.JobModel) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job.Status = entities.Queued
	s.queue = append(s.queue, job)
}

// Start continuously processes tasks from the queue
func (s *TaskScheduler) Start(inspector *Inspector) {
	go func() {
		for {
			s.mu.Lock()
			if len(s.queue) > 0 {
				job := s.queue[0]
				s.queue = s.queue[1:] // Remove the job from the queue
				s.mu.Unlock()

				s.workerPool <- struct{}{} // Block if max workers are busy
				go func(j *entities.JobModel) {
					defer func() { <-s.workerPool }()
					s.processJob(j, inspector) // Process the job
				}(job)
			} else {
				s.mu.Unlock()
				time.Sleep(1 * time.Second) // Wait a second before checking the queue again
			}
		}
	}()
}

// processJob processes individual tasks and updates job status
func (s *TaskScheduler) processJob(job *entities.JobModel, inspector *Inspector) {
	job.Status = entities.Processing
	err := inspector.Execute(job) // Call the relevant module based on Job.Type

	// Update job status based on execution result
	if err != nil {
		job.Status = entities.Failed
		logger.Log.Error("Job execution failed", zap.Error(err))
	} else {
		job.Status = entities.Completed
	}

	s.updateJobStatus(job) // Update job status in the datastore
}

// updateJobStatus updates the status in DataStore
func (s *TaskScheduler) updateJobStatus(job *entities.JobModel) {
	// TODO: Add code to update job status in your database
	logger.Log.Info("Job status updated", zap.Uint("jobID", job.ID), zap.String("status", string(job.Status)))
}
