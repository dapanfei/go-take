# go-take golang爬虫项目
## 项目简介
    使用golang编写的小说爬虫,可指定并发检测线程数量
## 使用方式

### 命令
    windows:(例)
    go-take.exe -listmain https://www.biqiuge.com/book/4772/  -p 50 -save  text/xxx.txt

    linux:(例)
    ./go-take -listmain https://www.biqiuge.com/book/4772/  -p 50 -save  text/xxx.txt
##### 参数说明：
    -listmain string
	Enter the URL of the novel catalogue (default "null") (指定目录URL)
    -p int
        Number of threads (default 1) (指定并发爬取线程数量)
    -save string
	Address of Storage Documents for Novels (default "text/1.txt") (指定结果存放文件)
