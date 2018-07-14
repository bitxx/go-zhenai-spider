## 珍爱网爬虫-golang版
### 简介
这个是《Google资深工程师深度讲解golang》中的爬虫项目，官方没有提供源码，小编自学期间逐步完善的，前端页面小编重找了个模板，效果还是不错的。  
目前该项目源码，小编更新到了：带有前端展示以及ElasticSearch存储并发版  
<img src="img/1.png" width = "800" height = "516"/>
### 注意
ElasticSearch中文分词不准确的缘故，因此带有中文的检索并不准确。
若要提高精确度，可以给ElastciSearch安装一个中文分词插件
### 环境
docker
docker中开个elasticsearch镜像
### 依赖包
为了方便，小编已经将下面这些包下载到src中了，可以直接用
1. gopm  
为了便于从google服务器拉取某些go包
```text
go get -u github.com/gpmgo/gopm
``` 
2. golang.org/x/net/html 
用于自动发现当前html页面是哪种编码
```text
gopm get -g -v golang.org/x/net/html
```
3. golang.org/x/text
用于对文本内容进行转码
```text
gopm get -g -v golang.org/x/text
```  
4. ElasticSearch客户端
不解释，这东西搞IT的都应该懂
小编的ElasticSearch服务器端用的5.x版本，因此安装的elastic.v5，具体请到该地址对号入座：https://github.com/olivere/elastic
```text
go get -v gopkg.in/olivere/elastic.v5
```
## 本项目使用
1. 所有终端初始时候都在`go-zhenai-spider/src`这个目录下
2. 第一个终端启动itemsaver，用于保存数据到elasticsearch：
    ```text
    go run crawler_distributed/persist/server/itemsaver.go --port=1234
    ```  
3. 第2个终端启动worker1，爬取解析数据：
    ```text
    go run crawler_distributed/worker/server/worker.go --port=9000
    ```
4. 第3个终端启动worker2，爬取解析数据:
    ```text
    go run crawler_distributed/worker/server/worker.go --port=9001
    ```
5. 第4个终端启动worker3，爬取解析数据:
    ```text
    go run crawler_distributed/worker/server/worker.go --port=9002
    ```
6. 第5个终端启动main，客户端：
    ```text
    go run itemsaver_host=":1234" worker_hosts=":9000,:9001,9002"
    ```
7. 分布式爬虫运行中
