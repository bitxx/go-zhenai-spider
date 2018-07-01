package scheduler

import "crawler/engine"

/**
 * 队列形式，每个worker一个channel
 * 代码更容易控制
 */


type QueuedScheduler struct {
	requestChan chan engine.Request
	workerChan  chan chan engine.Request //每个worker都是一个独立的chan
}

func (s *QueuedScheduler) WorkerChan() chan engine.Request {
	return make(chan engine.Request )
}

func (s *QueuedScheduler) Submit(r engine.Request) {
	s.requestChan <- r
}

//
func (s *QueuedScheduler) WorkerReady(w chan engine.Request) {
	s.workerChan <- w
}

func (s *QueuedScheduler) Run() {
	s.workerChan = make(chan chan engine.Request)
	s.requestChan = make(chan engine.Request)
	go func() {
		var requestQ []engine.Request
		var workerQ []chan engine.Request
		for {
			var activeRequest engine.Request
			var activeWorker chan engine.Request
			if len(requestQ)>0&&len(workerQ)>0{
				activeWorker = workerQ[0]
				activeRequest = requestQ[0]
			}
			select {
			case r := <-s.requestChan:
				requestQ = append(requestQ, r)
			case w := <-s.workerChan:
				workerQ = append(workerQ, w)
			case activeWorker<-activeRequest:
				workerQ = workerQ[1:]
				requestQ = requestQ[1:]
			}

		}
	}()
}
