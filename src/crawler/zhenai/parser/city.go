package parser

import (
	"crawler/engine"
	"regexp"
	"crawler_distributed/config"
)

var profilrRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*">([^<]+)</a>`)
var cityUrlRe = regexp.MustCompile(`href="(http://album.zhenai.com/zhenghun/[^"]+)"`)

func ParseCity(contents []byte, _ string) engine.ParseResult {
	matches := profilrRe.FindAllSubmatch(contents, -1) //-1表示匹配所有
	result := engine.ParseResult{}

	for _, m := range matches {
		result.Request = append(result.Request, engine.Request{
			Url: string(m[1]),
			Parser: NewProfileParse(string(m[2])),
		})
	}
	matchs := cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matchs {
		result.Request = append(result.Request, engine.Request{
			Url:        string(m[1]),
			Parser: engine.NewFuncParser(ParseCity,config.ParseCity),
		})
	}
	return result
}
