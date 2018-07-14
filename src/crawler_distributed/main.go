package main

import (
	"crawler/engine"
	"crawler/zhenai/parser"
	"crawler/scheduler"
	"crawler/persist"
	"crawler_distributed/config"
	itemsaver "crawler_distributed/persist/client"
	worker "crawler_distributed/worker/client"
	"net/rpc"
	"crawler_distributed/rpcsupport"
	"log"
	"flag"
	"strings"
)

var (
	itemSaverHost= flag.String("itemsaver_host","","itemsaver host")
	workerHosts = flag.String("worker_hosts","","worker hosts(comma separated)")
	)

func main() {
	flag.Parse()

	concurrentScheduler() //并发队列版，每个worker单独用一个channel
	//simpleScheduler()  //并发非队列版，公用一个channel
	//singleCity()  //单个城市直接请求
}

func concurrentScheduler() {
	itemsChan, err := itemsaver.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}

	pool := createClientPool(strings.Split(*workerHosts, ", "))
	processor := worker.CreateProcessor(pool)

	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      10,
		ItemChan:         itemsChan,
		RequestProcessor: processor,
	}

	e.Run(engine.Request{
		Url:    "http://www.zhenai.com/zhenghun",
		Parser: engine.NewFuncParser(parser.ParserCityList, config.ParseCityList),
	})
}

func simpleScheduler() {
	itemsChan, err := persist.ItemSaver("dating_profile")
	if err != nil {
		panic(err)
	}

	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.SimpleScheduler{},
		WorkerCount: 10,
		ItemChan:    itemsChan,
	}

	e.Run(engine.Request{
		Url:    "http://www.zhenai.com/zhenghun",
		Parser: engine.NewFuncParser(parser.ParserCityList, config.ParseCityList),
	})
}

func singleCity() {
	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 10,
	}

	e.Run(engine.Request{
		Url:    "http://www.zhenai.com/zhenghun/shanghai",
		Parser: engine.NewFuncParser(parser.ParseCity, config.ParseCity),
	})

}

func createClientPool(hosts []string) chan *rpc.Client {
	var clients []*rpc.Client
	for _, h := range hosts {
		client, err := rpcsupport.NewClient(h)
		if err == nil {
			clients = append(clients, client)
			log.Printf("connected to %s", h)
		} else {
			log.Printf("error connecting to %s：%v", h, err)
		}
	}

	out := make(chan *rpc.Client)
	go func() {
		for {
			for _, client := range clients {
				out <- client
			}
		}
	}()
	return out
}
