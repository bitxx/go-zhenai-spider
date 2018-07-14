package client

import (
	"crawler/engine"
	"crawler_distributed/config"
	"crawler_distributed/worker"
	"net/rpc"
)

func CreateProcessor(clientChan chan *rpc.Client) engine.Processor{

	return func(req engine.Request) (engine.ParseResult, error) {
		sReq := worker.SerializeRequest(req)
		var sResult worker.ParseResult
		client := <-clientChan
		err := client.Call(config.CrawlServiceRpc, sReq, &sResult)
		if err!=nil{
			return engine.ParseResult{},err
		}
		return worker.DesrializeResult(sResult),nil
	}
}
