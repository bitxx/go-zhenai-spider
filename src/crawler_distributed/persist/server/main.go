package main

import (
	"crawler_distributed/rpcsupport"
	"crawler_distributed/persist"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func main() {
	log.Fatal(serveRpc(":1234", "dating_profile"))
	//Fatal，若有异常，则挂了。panic还有recover的机会

}

func serveRpc(host, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}

	return rpcsupport.ServeRpc(host, &persist.ItemSaverService{
		Client: client,
		Index:  index,
	})
}
