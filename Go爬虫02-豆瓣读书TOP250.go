package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func main() {
	SpiderRun()
}

func SpiderRun() {
	fmt.Println("正在爬取豆瓣读书TOP250……")
	f, err := os.Create("豆瓣读书TOP250.txt")
	if err != nil {
		return
	}
	defer f.Close()

	for i := 1; i <= 10; i++ {
		SpiderPage(i, f)
	}
}

func SpiderPage(page int, f *os.File) {
	var txt string
	url := "https://book.douban.com/top250?start=" + strconv.Itoa((page-1)*25)
	pageHtml, err := SpiderGet(url)
	if err != nil {
		return
	}
	//fmt.Println(pageHtml)
	re := regexp.MustCompile(`<a href=".+?" onclick=&#34;moreurl.+?&#34; title="(.+?)"`)
	re2 := regexp.MustCompile(` <p class="pl">(.+?)</p>`)
	title := re.FindAllStringSubmatch(pageHtml, -1)
	detail := re2.FindAllStringSubmatch(pageHtml, -1)
	for i := 0; i < 25; i++ {
		no := (page-1)*25 + i + 1
		now := fmt.Sprintf("No.%03d %-30s%s \n", no, title[i][1], detail[i][1])
		txt += now
	}
	f.Write([]byte(txt))
	fmt.Printf("…………已完成 %d %% \n", page*10)
}

func SpiderGet(url string) (pageHtml string, err error) {
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36")

	client := new(http.Client)
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	buf := make([]byte, 8*1024)
	for {
		n, _ := resp.Body.Read(buf)
		if n == 0 {
			break
		}
		pageHtml += string(buf[:n])
	}
	return
}
