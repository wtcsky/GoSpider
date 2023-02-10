package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	var start, end int // 爬取的起始页和终止页
	fmt.Println("请输入起始页(≥1)：")
	fmt.Scan(&start)
	fmt.Println("请输入终止页(≥起始页)：")
	fmt.Scan(&end)

	SpiderRun(start, end) // 启动爬虫
}

func SpiderRun(start, end int) {
	fmt.Printf("正在爬取%d - %d页的内容\n", start, end)
	f, err := os.Create("古诗文.txt") // 创建文件存放爬取的内容
	if err != nil {
		return
	}
	defer f.Close() // 程序结束后关闭文件

	page := make(chan int)
	for i := start; i <= end; i++ { // 单独的爬取每一页
		go SpiderPage(i, f, page)
	}

	for i := start; i <= end; i++ {
		fmt.Printf("第%d页已经完成\n", <-page) // 防止主程序完成直接退出。
	}

}

func SpiderPage(i int, f *os.File, page chan int) {
	var txtContent string
	url := "https://so.gushiwen.cn/shiwens/default.aspx?page=" + strconv.Itoa(i) // 获取这一页的url
	fmt.Printf("正在爬取第%d页：%s\n", i, url)
	pageHtml, err := SpiderGet(url) // 获取这一页的源代码
	if err != nil {
		return
	}
	re1 := regexp.MustCompile(`<b>(.+?)</b>`)
	title := re1.FindAllStringSubmatch(pageHtml, -1)
	re2 := regexp.MustCompile(`<p class="source"><a href=".+?" target="_blank">(.+?)</a><a href=".+?">(.+?)</a></p>\s<div class="contson" id=".+?">([\s\S]+?)</div>`)
	content := re2.FindAllStringSubmatch(pageHtml, -1)
	for j := 0; j < 10; j++ {
		txtContent += title[j][1] + "\n" + content[j][1] + content[j][2] + content[j][3] + "\n"
	}
	txtContent = strings.Replace(txtContent, "<br />", "\n", -1)
	txtContent = strings.Replace(txtContent, "<p>　　", "", -1)
	txtContent = strings.Replace(txtContent, "</p>", "", -1)
	f.Write([]byte(txtContent))
	page <- i
}

func SpiderGet(url string) (pageHtml string, err error) {
	resp, err1 := http.Get(url) // 向网站发送get请求
	if err1 != nil {
		err = err1
		return
	}
	defer resp.Body.Close()     // 关闭响应体
	buf := make([]byte, 1024*8) // 用来接收返回的内容
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 { // 当n == 0时代表已经读取完毕。
			break
		}
		pageHtml += string(buf[:n]) // 将读取到的内容以string格式存放
	}
	return
}
