package parser

import (
	"crawler/engine"
	"regexp"
)

const cityListRe = `<a href="(http://www.zhenai.com/zhenghun/[0-9a-z]+)"[^>]*>([^<]+)</a>`

func ParserCityList(contents []byte,_ string) engine.ParseResult {
	re := regexp.MustCompile(cityListRe)
	matches := re.FindAllSubmatch(contents, -1) //-1表示匹配所有
	result := engine.ParseResult{}

	limit := 10  //显示请求量
	for _, m := range matches {
		result.Request = append(result.Request, engine.Request{
			Url:        string(m[1]),
			ParserFunc: ParseCity,
		})
		limit--
		if limit == 0{
			break
		}
	}
	return result
}
