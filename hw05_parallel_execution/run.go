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
	doneChan := make(chan struct{})
	wg := &sync.WaitGroup{}
	var errorsCount atomic.Int64

	// little optimisation when workers count (n) > tasks count
	workersCount := min(n, len(tasks))
	wg.Add(workersCount)

	for range workersCount {
		go func() {
			defer func() {
				wg.Done()
			}()

			for task := range taskChan {
				err := task()
				if err != nil {
					if errorsCount.Add(1) >= int64(m) {
						select {
						case doneChan <- struct{}{}:
						default:
						}

						return
					}
				}
			}
		}()
	}

	go func() {
		defer func() {
			close(taskChan)
		}()

		for _, task := range tasks {
			select {
			case taskChan <- task:
			case <-doneChan:
				return
			}
		}
	}()

	wg.Wait()
	close(doneChan)

	if errorsCount.Load() >= int64(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

// solution with context
// func Run(tasks []Task, n, m int) error {
// 	if m <= 0 {
// 		return ErrErrorsLimitExceeded
// 	}
//
// 	if n <= 0 {
// 		return ErrInvalidWorkersCount
// 	}
//
// 	taskChan := make(chan Task)
// 	wg := &sync.WaitGroup{}
// 	var errorsCount atomic.Int64
//
// 	ctx, cancel := context.WithCancel(context.TODO())
// 	defer cancel()
//
//  workersCount := min(n, len(tasks))
// 	wg.Add(workersCount)
//
// 	for i := 0; i < workersCount; i++ {
// 		go func() {
// 			defer wg.Done()
//
// 			for {
// 				select {
// 				case task, ok := <-taskChan:
// 					if !ok {
// 						return
// 					}
//
// 					err := task()
// 					if err != nil {
// 						if errorsCount.Add(1) >= int64(m) {
// 							cancel()
// 							return
// 						}
// 					}
//
// 				case <-ctx.Done():
// 					return
// 				}
// 			}
//
// 		}()
// 	}
//
// 	go func() {
// 		defer close(taskChan)
//
// 		for _, task := range tasks {
// 			select {
// 			case taskChan <- task:
// 			case <-ctx.Done():
// 				return
// 			}
// 		}
// 	}()
//
// 	wg.Wait()
//
// 	if errorsCount.Load() >= int64(m) {
// 		return ErrErrorsLimitExceeded
// 	}
//
// 	return nil
// }
