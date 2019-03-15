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
	"flag"
	"fmt"
)


var f *os.File
var URLSTR = "https://www.biqiuge.com"
var wg sync.WaitGroup

type Workdist struct {
	Url	string
}

const (
	taskload		    = 100
)


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
	_, err := io.WriteString(f, content) //写入文件(字符串)
	check(err)
}


//获取目录列表
func getListMain (listurl string,tasknum int){
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

func worker(tasks chan Workdist) {
	defer wg.Done()

	for{
		task,ok := <- tasks
		if !ok {
			//log.Print("通道关闭")
			return
		}
		getContent(URLSTR + task.Url)
	}
}

func getContent(newurl string){
	enc := mahonia.NewDecoder("gbk")
	doc, err := goquery.NewDocument(newurl)
	if err != nil {
		log.Fatal(err)
	}
	//目录获取
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





var listmain = flag.String("listmain", "null", "Enter the URL of the novel catalogue")
var tasknum = flag.Int("p", 1, "Number of threads")
var savepath = flag.String("save", "text/1.txt", "Address of Storage Documents for Novels")

func main() {

	flag.Parse()

	timestrat:=time.Now().Unix()
	var err1 error

	//检查文件存在，删除重建
	if checkFileIsExist(*savepath) {
		del := os.Remove(*savepath);
		if del != nil {
			fmt.Println(del);
		}
	}
	f, err1 = os.Create(*savepath)
	check(err1)
	getListMain(*listmain,*tasknum)
	log.Print("爬取使用时间为：",time.Now().Unix() - timestrat)
}
