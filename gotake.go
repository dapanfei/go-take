package main

import (
	_"reflect"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strings"
	"github.com/axgle/mahonia"
	"io"
	"sync"
	"time"
	"flag"
	"fmt"
	"sort"
)


type Workdist struct {
	Index int
	Url	string
}

const (
	taskload		    = 100
)
var (
	mutex *sync.Mutex
	f *os.File
	URLSTR = "https://www.biqiuge.com"
	wg sync.WaitGroup
	txtstr string
)


func check(e error) {
	if e != nil {
		log.Print(e)
		//panic(e)
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
	mutex.Lock()
	defer mutex.Unlock()
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
			Index:index,
			Url:link,
		}
		time.Sleep(time.Duration(10)*time.Millisecond)
		tasks <- task

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
		getContent(URLSTR + task.Url,task.Index)
	}
}
var mapText  map[int]string
func getContent(newurl string,index int){
	enc := mahonia.NewDecoder("gbk")
	log.Print("Link #",index,"  ","URL # ",newurl)
	doc, err := goquery.NewDocument(newurl)
	if err != nil {
		log.Fatal(err)
	}
	//目录获取
	var txt string
	doc.Find(".content h1").Each(func(i int, s *goquery.Selection) {
		centstr := s.Text()
		txt +=  "\r\n"+strings.Replace(enc.ConvertString(centstr), "聽聽聽聽", "", -1)+"\r\n"
	})
	//正文获取
	doc.Find(".showtxt").Each(func(i int, s *goquery.Selection) {
		centstr := s.Text()
		txt += strings.Replace(enc.ConvertString(centstr), "聽聽聽聽", "", -1)
		//log.Print(txt)
	})
	mapText[index] = txt
}

var listmain = flag.String("listmain", "null", "Enter the URL of the novel catalogue")
var tasknum = flag.Int("p", 1, "Number of threads")
var savepath = flag.String("save", "text/1.txt", "Address of Storage Documents for Novels")

func main() {

	flag.Parse()
	mutex = new(sync.Mutex)
	mapText = make(map[int]string)
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


	//map排序存储文件
	type kv struct {
		Key   int
		Value string
	}
	var ss []kv
	for k, v := range mapText {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Key < ss[j].Key  // 升序
	})

	for _, kv := range ss {
		txtstr += kv.Value
	}
	WriteWithIoutil(txtstr)
	log.Print("爬取使用时间为：",time.Now().Unix() - timestrat)
}