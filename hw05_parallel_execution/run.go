package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	taskChanel := make(chan Task)
	wg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}
	errorsCount := 0

	go func() {
		defer close(taskChanel)
		for _, v := range tasks {
			mutex.Lock()
			curErrorsCount := errorsCount
			mutex.Unlock()

			if curErrorsCount >= m {
				break
			}

			taskChanel <- v
		}
	}()

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range taskChanel {
				mutex.Lock()
				curErrorsCount := errorsCount
				mutex.Unlock()

				if curErrorsCount >= m {
					break
				}

				if task == nil { // Канал закрыт и кончился
					break
				}

				err := task()
				if err != nil {
					mutex.Lock()
					errorsCount++
					mutex.Unlock()
				}
			}
		}()
	}

	wg.Wait()

	if errorsCount >= m {
		return ErrErrorsLimitExceeded
	}

	return nil
}
