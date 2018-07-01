package main

import (
	"net/http"
	"crawler/frontend/controller"
)

func main() {
	http.Handle("/",http.FileServer(http.Dir("src/crawler/frontend/view")))
	http.Handle("/search",controller.CreateSearchResultHandler("src/crawler/frontend/view/index.html"))
	err := http.ListenAndServe(":8888", nil)
	if err !=nil{
		panic(err)
	}
}
