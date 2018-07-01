package fetcher

import (
	"net/http"
	"fmt"
	"golang.org/x/text/transform"
	"io/ioutil"
	"golang.org/x/text/encoding"
	"bufio"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/unicode"
	"log"
	"time"
)

var rateLimiter = time.Tick(10*time.Millisecond) //传入，延时爬取

func Fetch(url string) ([] byte, error) {
	<-rateLimiter //接收
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		//fmt.Println("Error: status code", resp.StatusCode)
		return nil,fmt.Errorf("wrong status code: %d",resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder()) //将抓取内容转为文本
	return ioutil.ReadAll(utf8Reader)
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024) //读取前1024个字节
	if err != nil {
		log.Printf("fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "") //自动判断编码类型
	return e
}
