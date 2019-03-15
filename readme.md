# go-take golang爬虫项目
## 项目简介
    使用golang编写的小说爬虫,可指定并发检测线程数量
## 使用方式

### 命令
    windows:(例)
    go-take.exe -listtmain https://www.biqiuge.com/book/4772/  -p 50 -save  text/xxx.txt
    
    linux:(例)
    ./go-take -listtmain https://www.biqiuge.com/book/4772/  -p 50 -save  text/xxx.txt
##### 参数说明：
    -listtmain   指定目录URL
    -proto  指定检测服务协议
    -p      指定并发爬取线程数量
    -save   指定结果存放文件