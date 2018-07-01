package scheduler

import "crawler/engine"

/**
 * 非队列，公用一个channel
 * 全都在抢一个channel，不可控制
 */

type SimpleScheduler struct {
	workerChan chan engine.Request
}

func (s *SimpleScheduler) WorkerChan() chan engine.Request {
	return s.workerChan
}

func (s *SimpleScheduler) WorkerReady(chan engine.Request) {
}

func (s *SimpleScheduler) Run() {
	s.workerChan = make(chan engine.Request)
}

func (s *SimpleScheduler) Submit(r engine.Request) {
	//send request down to worker chan
	go func() {s.workerChan <- r}()  //为什么用go，每个request用一个goroutine
}
