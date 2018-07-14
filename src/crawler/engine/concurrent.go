package engine

type ConcurrentEngine struct {
	Scheduler   Scheduler
	WorkerCount int
	ItemChan chan Item
	RequestProcessor Processor
}

type Processor func (r Request) (ParseResult, error)

type Scheduler interface {
	ReadyNotifier
	Submit(Request)
	WorkerChan() chan Request
	Run()
}

type ReadyNotifier interface {
	WorkerReady(chan Request)
}

func (e *ConcurrentEngine) Run(seeds ...Request) {

	out := make(chan ParseResult)
	e.Scheduler.Run()
	for i := 0; i < e.WorkerCount; i++ {
		e.createWorker(e.Scheduler.WorkerChan(),out,e.Scheduler) //里面有无限循环，channel输入输出
	}

	for _, r := range seeds {
		if isDuplicate(r.Url){
			continue
		}
		e.Scheduler.Submit(r)  //每个请求独立提交
	}

	for {
		result := <-out
		for _, item := range result.Items {
			go func() {
				e.ItemChan <- item
			}()
		}
		for _, request := range result.Request {
			if isDuplicate(request.Url){
				continue
			}
			e.Scheduler.Submit(request) //上一级请求结果中若有新的请求，则继续发起请求
		}
	}
}

func  (e *ConcurrentEngine) createWorker(in chan Request,out chan ParseResult,ready ReadyNotifier) {
	go func() {
		for {
			ready.WorkerReady(in)
			request := <-in
			result, err := e.RequestProcessor(request)
			if err!=nil{
				continue
			}
			out <- result
		}
	}()
}

var visitedUrls= make(map[string] bool)

//去重
func isDuplicate(url string) bool{
	if visitedUrls[url]{
		return true
	}
	visitedUrls[url] = true
	return false
}
