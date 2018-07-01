package controller

import (
	"crawler/frontend/view"
	"gopkg.in/olivere/elastic.v5"
	"net/http"
	"strings"
	"strconv"
	"crawler/frontend/view/model"
	"context"
	"reflect"
	"crawler/engine"
	"regexp"
)

type SearchResultHandler struct {
	view view.SearchResultView
	client *elastic.Client
}

func CreateSearchResultHandler(template string) SearchResultHandler{
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err!=nil{
		panic(err)
	}
	return SearchResultHandler{
		view:view.CreateSearchResultView(template),
		client:client,
	}
}

//localhost:9200/search?q=男 已购房&from=20
func (h SearchResultHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	q := strings.TrimSpace(req.FormValue("q"))
	from,err := strconv.Atoi(req.FormValue("from"))
	if err!=nil{
		from = 0
	}
	//fmt.Fprintf(w,"q=%s, from=%d",q,from)  //响应结果，写入w
	page,err := h.getSearchResult(q,from)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
	}
	err = h.view.Render(w, page)
	if err!=nil{
		http.Error(w,err.Error(),http.StatusBadRequest)
	}

}

func (h SearchResultHandler) getSearchResult(q string,from int)(model.SearchResult,error){
	var result model.SearchResult
	result.Query = q
	resp, err := h.client.Search("dating_profile").Query(elastic.NewQueryStringQuery(rewriteQueryString(q))).From(from).Do(context.Background())
	if err!=nil{
		return result,err
	}
	result.Hits = resp.TotalHits()
	result.Start = from
	result.Items = resp.Each(reflect.TypeOf(engine.Item{}))
	result.PrevFrom = result.Start - len(result.Items)
	result.NextFrom = result.Start + len(result.Items)

	return result,nil
}

//前端搜索没有输入Payload.Age<30，而是输入的Age<30,因此为了保证能在elastic中正常使用，重新添加上Payload.
func rewriteQueryString(q string) string{
	re := regexp.MustCompile(`([A-Z][a-z]*):`)//([A-Z][a-z]*)表示Height: Age:
	return re.ReplaceAllString(q,"Payload.$1:")
}



