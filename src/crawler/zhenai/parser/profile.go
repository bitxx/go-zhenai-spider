package parser

import (
	"crawler/engine"
	"regexp"
	"strconv"
	"crawler/model"
	"crawler_distributed/config"
)

var ageRe = regexp.MustCompile(`<td><span class="label">年龄：</span>([\d]+)岁</td>`)
var hrightRe = regexp.MustCompile(`<td><span class="label">身高：</span>([\d]+)CM</td>`)
var marriageRe = regexp.MustCompile(`<td><span class="label">婚况：</span>([^<]+)</td>`)
var incomeRe = regexp.MustCompile(`<td><span class="label">月收入：</span>([^<]+)</td>`)
var weightRe = regexp.MustCompile(`<td><span class="label">体重：</span><span field="">([^<]+)</span></td>`)
var genderRe = regexp.MustCompile(`<td><span class="label">性别：</span><span field="">([^<]+)</span></td>`)
var xinzuoRe = regexp.MustCompile(`<td><span class="label">星座：</span><span field="">([^<]+)</span></td>`)
var educationRe = regexp.MustCompile(`<td><span class="label">学历：</span><span field="">([^<]+)</span></td>`)
var hokouRe = regexp.MustCompile(`<td><span class="label">籍贯：</span>([^<]+)</td>`)
var houseRe = regexp.MustCompile(`<td><span class="label">住房条件：</span><span field="">([^<]+)</span></td>`)
var carRe = regexp.MustCompile(`<td><span class="label">是否购车：</span><span field="">([^<]+)</span></td>`)
var occupationRe = regexp.MustCompile(`<td><span class="label">职业： </span>([^<]+)</td>`)
//猜你喜欢的
var guessRe = regexp.MustCompile(`<a class="exp-user-name"[^>]*href="(http://album.zhenai.com/u/[\d]+)">([^<])`)
var idUrlRe = regexp.MustCompile(`http://album.zhenai.com/u/([\d]+)`)

func parseProfile(contents [] byte,url string, name string) engine.ParseResult{
	profile := model.Profile{}
	profile.Name = name

	age,err := strconv.Atoi(extractString(contents,ageRe))
	if err == nil{
		profile.Age = age
	}

	height,err := strconv.Atoi(extractString(contents,hrightRe))
	if err == nil{
		profile.Weight = height
	}

	weight,err := strconv.Atoi(extractString(contents,weightRe))
	if err == nil{
		profile.Weight = weight
	}

	profile.Marrige =extractString(contents,marriageRe)
	profile.Income =extractString(contents,incomeRe)
	profile.Gender =extractString(contents,genderRe)
	profile.Car =extractString(contents,carRe)
	profile.Education =extractString(contents,educationRe)
	profile.Hokou =extractString(contents,hokouRe)
	profile.House =extractString(contents,houseRe)
	profile.Xinzuo =extractString(contents,xinzuoRe)
	profile.Occupation =extractString(contents,occupationRe)

	result := engine.ParseResult{
		Items:[] engine.Item {
			{
				Url:url,
				Type:"zhenai",
				Id:extractString([]byte(url),idUrlRe),
				Payload:profile,
			},
		},
	} //[]interface{}{profile}加入一个元素

	//猜你喜欢的人
	matchs := guessRe.FindAllSubmatch(contents,-1) //-1表示匹配所有
	for _,m := range matchs{
		result.Request = append(result.Request,engine.Request{
			Url:string(m[1]),
			Parser: NewProfileParse(string(m[2])),
		})
	}
	return result
}

func extractString(contents []byte, re *regexp.Regexp) string{
	match := re.FindSubmatch(contents)
	if len(match)>=2{
		return string(match[1])
	}else {
		return ""
	}
}

type ProfileParse struct{
	userName string
}

func (p *ProfileParse) Parser(contents []byte, url string) engine.ParseResult {
	return parseProfile(contents,url,p.userName)
}

func (p *ProfileParse) Serialize() (name string, args interface{}) {
	return config.ParseProfile,p.userName
}

func NewProfileParse(name string) *ProfileParse{
	return &ProfileParse{
		userName:name,

	}
}