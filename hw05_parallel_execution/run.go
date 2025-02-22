package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
	ErrInvalidWorkersCount = errors.New("invalid workers count")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	if n <= 0 {
		return ErrInvalidWorkersCount
	}

	taskChan := make(chan Task)
	wg := &sync.WaitGroup{}
	var errorsCount atomic.Int64

	// little optimisation when workers count (n) > tasks count
	workersCount := min(n, len(tasks))
	wg.Add(workersCount)

	for range workersCount {
		go func() {
			defer wg.Done()

			for task := range taskChan {
				if err := task(); err != nil {
					errorsCount.Add(1)
				}
			}
		}()
	}

	for _, task := range tasks {
		taskChan <- task
		if errorsCount.Load() >= int64(m) {
			break
		}
	}

	close(taskChan)
	wg.Wait()

	if errorsCount.Load() >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
