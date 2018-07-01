package persist

import (
	"testing"
	"crawler/model"
	"gopkg.in/olivere/elastic.v5"
	"context"
	"encoding/json"
	"crawler/engine"
)

func TestSave(t *testing.T) {
	expected := engine.Item{
		Url:  "http://album.zhenai.com/u/108906739",
		Type: "zhenai",
		Id:   "108906739",
		Payload: model.Profile{
			Age:        34,
			Height:     162,
			Weight:     57,
			Income:     "3001-5000元",
			Gender:     "女",
			Name:       "安静的雪",
			Xinzuo:     "牧羊座",
			Occupation: "人事/行政",
			Marrige:    "离异",
			House:      "已够房",
			Hokou:      "山东菏泽",
			Education:  "大学本科",
			Car:        "未购车",
		},
	}



	//不想依赖外部环境，因此需要在这里单独启动docker,启动elastcisearch
	//在此酒不展示了
	client, err := elastic.NewClient(
		elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	const index = "dating_profile"
	//1.保存
	err = save(client,index,expected)
	if err != nil {
		panic(err)
	}

	//2.获取
	resp, err := client.Get().Index("dating_profile").Type(expected.Type).Id(expected.Id).Do(context.Background())
	if err != nil {
		panic(err)
	}
	t.Logf("%s", resp.Source)

	//3.比较
	var actual engine.Item
	json.Unmarshal(*resp.Source, &actual)
	actualProfile,_:=model.FromJsonObj(actual.Payload)
	actual.Payload = actualProfile
	if actual != expected {
		t.Errorf("got %v; expected %v", actual, expected)
	}
}
