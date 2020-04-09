package work

import "sync"

// Worker 接口定义
type Worker interface {
	Task()
}

// WorkPool 定义一个goroutine池（工作池）
type WorkPool struct {
	work chan Worker
	wg sync.WaitGroup
}

// New 创建一个新的工作池
// 开启maxGoroutines个协程来处理所有的工作任务
func New(maxGoroutines int) *WorkPool {
	p := WorkPool{
		work: make(chan Worker),
	}

	p.wg.Add(maxGoroutines)
	for i := 0; i < maxGoroutines; i++ {
		go func() {
			for w := range p.work {
				w.Task()
			}
			p.wg.Done()
		}()
	}

	return &p
}

// Run 提交工作任务到工作池
func (p *WorkPool) Run(w Worker) {
	p.work <- w
}

// Shutdown 等待所有的goroutine停止工作
func (p *WorkPool) Shutdown() {
	close(p.work)
	p.wg.Wait()
}