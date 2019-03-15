package main

import (
	_"reflect"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strings"
	"mahonia"
	"io"
	"sync"
	"time"
)


var f *os.File
var nexturl string
var name = "text/圣墟.txt"

var URLSTR = "https://www.biqiuge.com"


type Workdist struct {
	Url	string
}

const (
	taskload		    = 100
)
var tasknum = 100
var wg sync.WaitGroup



//获取目录列表
func getlist (listurl string){
	tasks := make(chan Workdist,taskload)

	wg.Add(tasknum)
	for gr:=1;gr<=tasknum;gr++ {
		go	worker(tasks)
	}

	list, err := goquery.NewDocument(listurl)
	if err != nil {
		log.Fatal(err)
	}
	list.Find(".listmain a").Each(func(index int, s *goquery.Selection) {
		linkTag := s
		link, _ := linkTag.Attr("href")
		//linkText := linkTag.Text()
		task := Workdist{
			Url:link,
		}
		time.Sleep(time.Duration(10)*time.Millisecond)

		tasks <- task
		//fmt.Printf("Link #%d: '%s'\n", index, link)
	})

	close(tasks)
	wg.Wait()

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func WriteWithIoutil(content string) {
	//mutex.Lock()
	//defer mutex.Unlock()
	_, err := io.WriteString(f, content) //写入文件(字符串)
	check(err)

}

func worker(tasks chan Workdist) {
	defer wg.Done()

	for{
		task,ok := <- tasks
		if !ok {
			//log.Print("通道关闭")
			return
		}
		zhengwen(URLSTR + task.Url)
	}


}
var mutex sync.Mutex
func zhengwen(newurl string){
	enc := mahonia.NewDecoder("gbk")
	doc, err := goquery.NewDocument(newurl)
	if err != nil {
		log.Fatal(err)
	}
	doc.Find(".content h1").Each(func(i int, s *goquery.Selection) {
		centstr := s.Text()
		WriteWithIoutil("\r\n"+strings.Replace(enc.ConvertString(centstr), "聽聽聽聽", "", -1)+"\r\n")
		log.Print(strings.Replace(enc.ConvertString(centstr), "聽聽聽聽", "", -1))
	})
	//正文获取
	doc.Find(".showtxt").Each(func(i int, s *goquery.Selection) {
		centstr := s.Text()
		WriteWithIoutil(strings.Replace(enc.ConvertString(centstr), "聽聽聽聽", "", -1))
		//log.Print(strings.Replace(enc.ConvertString(centstr), "聽聽聽聽", "", -1))
	})
}



func main() {
	timestrat:=time.Now().Unix()
	var err1 error
	f, err1 = os.Create(name)
	check(err1)
	getlist("https://www.biqiuge.com/book/4772/")
	log.Print("爬取使用时间为：",time.Now().Unix() - timestrat)
}
