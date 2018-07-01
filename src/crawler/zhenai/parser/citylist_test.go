package parser

import (
	"testing"
	"io/ioutil"
)

func TestParserCityList(t *testing.T) {
	/*contents, err := fetcher.Fetch("http://www.zhenai.com/zhenghun")
	if err!=nil{
		panic(err)
	}
	fmt.Printf("%s",contents)*/
	contents, err := ioutil.ReadFile("citylist.txt")
	if err!=nil{
		panic(err)
	}
	result := ParserCityList(contents,"")

	const resultSize = 470

	expectedUrls := []string{
		"http://www.zhenai.com/zhenghun/aba","http://www.zhenai.com/zhenghun/akesu","http://www.zhenai.com/zhenghun/alashanmeng",
	}

	if len(result.Request)!=resultSize{
		t.Errorf("result should have %d"+",but had %d",resultSize,len(result.Request))
	}
	for i,url := range expectedUrls{
		if result.Request[i].Url != url{
			t.Errorf("expected url #%d: %s but was %s",i,url,result.Request[i].Url)
		}
	}

	if len(result.Items)!=resultSize{
		t.Errorf("result should have %d"+",but had %d",resultSize,len(result.Items))
	}

	//fmt.Printf("%s",contents)
}
