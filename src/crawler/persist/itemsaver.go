package persist

import (
	"log"
	"gopkg.in/olivere/elastic.v5"
	"context"
	"crawler/engine"
	"github.com/pkg/errors"
)

func ItemSaver(index string) (chan engine.Item,error) {
	client, err := elastic.NewClient(elastic.SetSniff(false)) //客户端维护集群，但集群不在本地，因此不需要维护
	if err != nil {
		return nil,err
	}

	out := make(chan engine.Item)
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver: Got item #%d: %v", itemCount, item)
			itemCount++
			err := Save(client,index, item)
			if err != nil {
				log.Printf("Item Saver:error saving item %v : %v", item, err)
			}

		}

	}()
	return out,nil
}

func Save(client *elastic.Client,index string, item engine.Item) (err error) {

	if item.Type == "" {
		return errors.New("must supply Type")
	}
	indexService := client.Index().Index(index).Type(item.Type).BodyJson(item)
	if item.Id != "" {
		indexService.Id(item.Id)
	}
	_, err = indexService.Do(context.Background())
	if err != nil {
		return err
	}
	//fmt.Printf("%+v",response)  //+v带有key
	return nil
}
