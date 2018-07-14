package main

import (
	"crawler/engine"
	"crawler/zhenai/parser"
	"crawler/scheduler"
	"crawler/persist"
)

func main() {
	concurrentScheduler() //并发队列版，每个worker单独用一个channel
	//simpleScheduler()  //并发非队列版，公用一个channel
	//singleCity()  //单个城市直接请求
}

func concurrentScheduler(){
	itemsChan, err := persist.ItemSaver("dating_profile")
	if err!=nil{
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler: &scheduler.QueuedScheduler{},
		WorkerCount:10,
		ItemChan:itemsChan,
		RequestProcessor:engine.Worker,
	}

	e.Run(engine.Request{
		Url: "http://www.zhenai.com/zhenghun",
		Parser: engine.NewFuncParser(parser.ParserCityList,"ParserCityList"),
	})
}

func simpleScheduler(){
	itemsChan, err := persist.ItemSaver("dating_profile")
	if err!=nil{
		panic(err)
	}

	e := engine.ConcurrentEngine{
		Scheduler: &scheduler.SimpleScheduler{},
		WorkerCount:10,
		ItemChan:itemsChan,
	}

	e.Run(engine.Request{
		Url: "http://www.zhenai.com/zhenghun",
		Parser: engine.NewFuncParser(parser.ParserCityList,"ParserCityList"),
	})
}

func singleCity(){
	e := engine.ConcurrentEngine{
		Scheduler: &scheduler.QueuedScheduler{},
		WorkerCount:10,
	}

	e.Run(engine.Request{
		Url: "http://www.zhenai.com/zhenghun/shanghai",
		Parser: engine.NewFuncParser(parser.ParserCityList,"ParserCityList"),
	})

}

