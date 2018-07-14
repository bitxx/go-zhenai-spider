package worker

import (
	"crawler/engine"
	"crawler_distributed/config"
	"crawler/zhenai/parser"
	"github.com/pkg/errors"
	"fmt"
	"log"
)

type SerializedParser struct {
	Name string
	Args interface{}
}

type Request struct {
	Url    string
	Parser SerializedParser
}

type ParseResult struct {
	Items    []engine.Item
	Requests []Request
}

func SerializeRequest(r engine.Request) Request {
	name, args := r.Parser.Serialize()

	return Request{
		Url: r.Url,
		Parser: SerializedParser{
			Name: name,
			Args: args,
		},
	}
}

func SerializeResult(r engine.ParseResult) ParseResult {
	result := ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Request {
		result.Requests = append(result.Requests, SerializeRequest(req))
	}
	return result
}

func DesrializeRequest(r Request) (engine.Request, error) {
	parser, err := deserializeParser(r.Parser)
	if err != nil {
		return engine.Request{},err
	}
	return engine.Request{
		Url:    r.Url,
		Parser: parser,
	},nil
}

func DesrializeResult(r ParseResult) engine.ParseResult{

	result := engine.ParseResult{
		Items: r.Items,
	}
	for _, req := range r.Requests {
		engineReq, err := DesrializeRequest(req)
		if err!=nil{
			log.Printf("erro deserilizing request:%v",err)
			continue
		}
		result.Request = append(result.Request, engineReq)
	}

	return result
}

func deserializeParser(p SerializedParser) (engine.Parser, error) {
	switch p.Name {
	case config.ParseCityList:
		return engine.NewFuncParser(parser.ParserCityList, config.ParseCityList), nil
	case config.ParseCity:
		return engine.NewFuncParser(parser.ParseCity, config.ParseCity), nil
	case config.NilParser:
		return engine.NilParser{}, nil
	case config.ParseProfile:
		if userName, ok := p.Args.(string); ok {
			return parser.NewProfileParse(userName), nil
		} else {
			return nil, fmt.Errorf("invalid arg:%v", p.Args)
		}

	default:

		return nil, errors.New("unknow parser name")

	}
}
