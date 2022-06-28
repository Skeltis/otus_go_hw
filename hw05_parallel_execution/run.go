package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded      = errors.New("errors limit exceeded")
	ErrNegativeErrorThreshold   = errors.New("errors limit can't be less than 0")
	ErrLessThanOneParallelCount = errors.New("can't execute with less than 1 parallel task")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 1 {
		return ErrLessThanOneParallelCount
	}

	if m < 0 {
		return ErrNegativeErrorThreshold
	}

	if !doesAnyTasksPresent(tasks) {
		return nil
	}

	taskHandle := newTaskHandler(tasks, n, m)
	taskHandle.start()
	return taskHandle.waitCompletion()
}

type taskHandler struct {
	tasks          []Task
	taskThreshold  int64
	errorThreshold int64
	taskPosition   int64
	totalTasks     int64
	errorCount     int64
	resultChannel  chan error
	finishLock     sync.WaitGroup
}

func doesAnyTasksPresent(tasks []Task) bool {
	nonNilTasksPresent := false
	for _, task := range tasks {
		if task != nil {
			nonNilTasksPresent = true
		}
	}
	return nonNilTasksPresent
}

func newTaskHandler(tasks []Task, parallelTasksNumber, maxErrors int) *taskHandler {
	handler := &taskHandler{
		tasks:          tasks,
		taskThreshold:  int64(parallelTasksNumber),
		errorThreshold: int64(maxErrors),
		totalTasks:     int64(len(tasks)),
		resultChannel:  make(chan error),
	}

	if maxErrors == 0 {
		handler.errorCount = -1
	}

	return handler
}

func (handler *taskHandler) start() {
	if handler.taskThreshold < handler.totalTasks {
		atomic.AddInt64(&handler.taskPosition, handler.taskThreshold)
		handler.finishLock.Add(int(handler.taskThreshold))
	} else {
		atomic.AddInt64(&handler.taskPosition, handler.totalTasks)
		handler.finishLock.Add(int(handler.totalTasks))
	}
	go handler.waitForTasksCompletion()

	var i int64
	for i = 0; i < handler.taskThreshold && i < handler.totalTasks; i++ {
		go runTask(handler.tasks[i], handler.resultChannel)
	}
}

func (handler *taskHandler) waitCompletion() error {
	handler.finishLock.Wait()
	close(handler.resultChannel)
	if handler.errorCount >= handler.errorThreshold {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func (handler *taskHandler) waitForTasksCompletion() {
	for {
		select {
		case result, ok := <-handler.resultChannel:
			if !ok {
				return
			}

			if result != nil {
				atomic.AddInt64(&handler.errorCount, 1)
			}

			if handler.errorCount < handler.errorThreshold && handler.taskPosition < handler.totalTasks {
				go runTask(handler.tasks[handler.taskPosition], handler.resultChannel)
				atomic.AddInt64(&handler.taskPosition, 1)
			} else {
				handler.finishLock.Done()
			}
		}
	}
}

func runTask(curTask Task, resultChannel chan<- error) {
	if curTask != nil {
		resultChannel <- curTask()
	} else {
		resultChannel <- nil
	}
}
