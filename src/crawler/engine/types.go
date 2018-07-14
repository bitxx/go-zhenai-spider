package engine

type ParserFunc func(contents []byte,url string) ParseResult

type Parser interface {
	Parser(contents []byte,url string) ParseResult
	Serialize() (name string,args interface{})
}

//一条请求
type Request struct{
	Url string
	Parser Parser
}


//{"ParserCityList",nil}{"ProfileParse",userName}

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

type NilParser struct{

}

func (NilParser) Parser(_ []byte, _ string) ParseResult {
	return ParseResult{}
}

func (NilParser) Serialize() (name string, args interface{}) {
	return "NilParser",nil
}

type FuncParser struct{
	parser ParserFunc
	Name string
}

func (f *FuncParser) Parser(contents []byte, url string) ParseResult {
	return f.parser(contents,url)
}

func (f *FuncParser) Serialize() (name string, args interface{}) {
	return f.Name,nil
}

func NewFuncParser(p ParserFunc,name string) *FuncParser{
	return &FuncParser{
		parser:p,
		Name:name,
	}
}