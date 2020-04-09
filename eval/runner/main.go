package main

import (
	"time"
	"elen/concurrent-task-schedule/logger"
	"elen/concurrent-task-schedule/runner"
	"os"
	"fmt"
)

const timeout = 3 * time.Second

func main() {
	logger.Info("start working ...")

	// Runner实例化
	r := runner.New(timeout)

	// 添加任务列表
	r.AddTask(createTask(), createTask(), createTask())

	if err := r.Start(); err != nil {
		switch err {
		case runner.ErrTimeout:
			logger.Error("Terminate due to timeout.")
			os.Exit(1)
		case runner.ErrInterrupt:
			logger.Error("Terminate due to interrupt.")
			os.Exit(2)
		}
	}

	logger.Info("task list all process ok!")
}

// 模拟创建一个任务
func createTask() func(int) {
	return func(id int) {
		logger.Info(fmt.Sprintf("Processor - Task #%d.",id))

		//time.Sleep(time.Duration(id) * time.Second)
		time.Sleep(time.Duration(1/2) * time.Second)
	}
}
