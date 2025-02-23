package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("m is equal zero", func(t *testing.T) {
		err := Run(nil, 0, 0)
		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
	})

	t.Run("n is equal zero", func(t *testing.T) {
		err := Run(nil, 0, 1)
		require.Truef(t, errors.Is(err, ErrInvalidWorkersCount), "actual err - %v", err)
	})

	t.Run("check concurrent execution with eventually", func(t *testing.T) {
		const (
			tasksCount   = 500
			workersCount = 5
		)

		var (
			mu             sync.Mutex
			maxConcurrent  = 0
			currentWorkers = 0
		)

		tasks := make([]Task, tasksCount)
		for i := range tasks {
			tasks[i] = func() error {
				mu.Lock()
				currentWorkers++
				if currentWorkers > maxConcurrent {
					maxConcurrent = currentWorkers
				}
				mu.Unlock()

				// do some work + try to change goroutines
				runtime.Gosched()

				mu.Lock()
				currentWorkers--
				mu.Unlock()

				return nil
			}
		}

		err := Run(tasks, workersCount, 1)
		require.NoError(t, err)

		require.Eventually(t, func() bool {
			mu.Lock()
			defer mu.Unlock()

			return maxConcurrent == workersCount
		}, time.Second, 10*time.Millisecond,
			"expected %d concurrent tasks, got %d",
			workersCount, maxConcurrent)
	})
}
