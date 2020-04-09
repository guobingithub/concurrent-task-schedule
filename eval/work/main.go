package main

import (
	"elen/concurrent-task-schedule/logger"
	"time"
	"elen/concurrent-task-schedule/work"
	"sync"
)

var nameList = []string{
	"GuoBin",
	"elen",
	"郭斌",
}

// namePrint 模拟一个业务
type namePrint struct {
	name string
}

// namePrint 实现Worker接口
func (np *namePrint) Task() {
	logger.Info(np.name)
	time.Sleep(time.Second)
}

func main() {
	// 使用100个goroutine来创建工作池
	wp := work.New(100)

	var wg sync.WaitGroup
	wg.Add(100 * len(nameList))

	for i:=0; i<100; i++ {
		for _,name := range nameList {
			// 初始化一个业务实例
			np := namePrint{
				name: name,
			}

			go func() {
				//启动协程，将任务提交执行。当Run任务执行完时退出协程
				wp.Run(&np)
				wg.Done()
			}()
		}
	}

	wg.Wait()

	// 让工作池停止工作，等待所有运行的工作完成
	wp.Shutdown()
}
