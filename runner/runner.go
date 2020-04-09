package runner

import (
	"os"
	"time"
	"errors"
	"os/signal"
	"elen/concurrent-task-schedule/logger"
	"fmt"
)

// 定义错误码
var (
	ErrTimeout 		= errors.New("recv timeout")
	ErrInterrupt 	= errors.New("recv interrupt")
)

type taskFunc func(int)

// Runner 在给定的超时时间内执行一组任务
type Runner struct {
	// 接收系统的信号
	interrupt chan os.Signal

	// 任务完成情况
	over chan error

	// 任务超时与否
	timeout <-chan time.Time

	// 任务列表
	tasks []taskFunc
}

func New(d time.Duration) *Runner {
	return &Runner{
		interrupt: 	make(chan os.Signal, 1),
		over: 		make(chan error),
		timeout:	time.After(d),
	}
}

func (r *Runner) AddTask(tasks ...taskFunc) {
	r.tasks = append(r.tasks, tasks...)
}

func (r *Runner) Start() error {
	signal.Notify(r.interrupt, os.Interrupt)

	go func() {
		r.over <- r.run()
	}()

	select {
	case err := <-r.over:
		return err
	case <-r.timeout:
		return ErrTimeout
	}
}

func (r *Runner) run() error {
	for id, task := range r.tasks {
		//whether interrupt occur
		if r.hasInterrupt() {
			return ErrInterrupt
		}

		//process task
		task(id)

		logger.Info(fmt.Sprintf("Task #%d end.",id))
	}

	return nil
}

func (r *Runner) hasInterrupt() bool {
	select {
	case <-r.interrupt:
		signal.Stop(r.interrupt)
		return true

	default:
		return false
	}
}
