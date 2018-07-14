package engine

import (
	"log"
)

type SimpleEngine struct {}

func (SimpleEngine)Run(seeds ...Request){
	var requests []Request
	for _,r := range seeds{
		requests = append(requests,r)
	}

	for len(requests)>0{
		r :=requests[0]
		requests = requests[1:]
		parseResult, err := Worker(r)
		if err!=nil{
			continue
		}
		requests = append(requests,parseResult.Request...)

		for _,item := range parseResult.Items{
			log.Printf("Got item %v",item)
		}

	}
}


