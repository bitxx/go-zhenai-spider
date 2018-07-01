package engine

type ParserFunc func(contents []byte,url string) ParseResult

//一条请求
type Request struct{
	Url string
	ParserFunc ParserFunc
}

//解析出来的结果包含新的请求和内容
type ParseResult struct{
	 Request []Request  //所有请求，切片
	 Items []Item //所有的内容，切片
}

type Item struct {
	Url string
	Id string
	Type string
	Payload interface{}
}

func NilParser([]byte) ParseResult{
	return ParseResult{}
}
