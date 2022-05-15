package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/*
	获取页面的json结构体
	在线转换处理 https://oktools.net/json2go
*/
type AutoJson struct {
	Status  int         `json:"status"`
	Info    string      `json:"info"`
	AllData CompanyData `json:"data"`
}
type List struct {
	Avatar      string `json:"avatar"`
	CompanyID   string `json:"company_id"`
	CompanyName string `json:"company_name"`
}
type CompanyData struct {
	Count   int    `json:"count"`
	Current int    `json:"current"`
	List    []List `json:"list"`
}

// OpenFile 判断文件是否存在  存在则OpenFile 不存在则Create
func OpenFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Println("文件不存在")
		return os.Create(filename) //创建文件
	}
	fmt.Println("文件存在")
	return os.OpenFile(filename, os.O_APPEND, 0666) //打开文件
}

func getrandomtime() int {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(500) + 500
	fmt.Println(r)
	return r
}

func getfunc(url, cookiessend string) string {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Cookie", cookiessend)

	resp, err := client.Do(req)
	if err != nil {
		panic(nil)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return string(body)
}

func postfunc(url, posttype, allcontent, cookiessend string) string {

	resp, err := http.Post(url, posttype, strings.NewReader(allcontent))
	//设置Header
	resp.Header.Set("Host", "www.butian.net")
	resp.Header.Set("Connection", "close")
	resp.Header.Set("Content-Length", "14")
	resp.Header.Set("sec-ch-ua", "\"Not A;Brand\";v=\"99\", \"Chromium\";v=\"100\", \"Google Chrome\";v=\"100\"")
	resp.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	resp.Header.Set("Content-Type", " application/x-www-form-urlencoded; charset=UTF-8")
	resp.Header.Set("X-Requested-With", "XMLHttpRequest")
	resp.Header.Set("sec-ch-ua-mobile", "?0")
	resp.Header.Set("User-Agent", " Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
	resp.Header.Set("sec-ch-ua-platform", "\"Windows\"")
	resp.Header.Set("Origin", "https://www.butian.net")
	resp.Header.Set("Sec-Fetch-Site", "same-origin")
	resp.Header.Set("Sec-Fetch-Mode", "cors")
	resp.Header.Set("Sec-Fetch-Dest", "empty")
	resp.Header.Set("Referer", "https://www.butian.net/Reward/pub/Message/send")
	resp.Header.Set("Accept-Encoding", "gzip, deflate")
	resp.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	//设置cookie
	resp.Header.Set("Cookie", cookiessend)

	if err != nil {
		fmt.Println("http.Get err=", err)
	}

	bytess, err := ioutil.ReadAll(resp.Body) // 获取网页返回的内容

	return string(bytess)
}

func webdeal(url, cookiessend string) string{ //爬取厂商的URL
	time.Sleep(time.Duration(getrandomtime()) * time.Millisecond) // 每获取底部页面依次，即睡眠500~1000毫秒

	company_url := ""

	urlcontent1 := getfunc(url, cookiessend)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(urlcontent1))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find(".tabs-con.tabs-con-loo ul li input").Each(func(i int, selection *goquery.Selection){
		//fmt.Println(urls)
		temp_values,_ := selection.Attr("value") //匹配所有字段内容
		temp_placeholder,_ := selection.Attr("placeholder") // 仅有URL内容被写入文件

		if(temp_placeholder=="输入所属域名或ip，若存在多个以,分隔") {
			if(temp_values == ""){
				println("未获取到厂商URL",url)
			}else{
				println("厂商URL内容为",temp_values)
				company_url = temp_values
			}
		}
	})
	return company_url
}

func main() {
	start := time.Now() // 计算代码运行时间

	//这里输入登录的cookie值
	cookiessend := ""

	fist_html := postfunc("https://www.butian.net/Reward/pub","application/x-www-form-urlencoded","s=1&p=1&token=",cookiessend)

	xin := AutoJson{}

	errxin := json.Unmarshal([]byte(fist_html), &xin) // 新的结构体用于处理获取的json数据
	if errxin != nil {
		log.Print(errxin)
	} else {
		log.Print(xin)
	}

	allpage := xin.AllData.Count // 获取所有的页面数量

	var all_info []string // 创建保存所有公司id信息的列表

	/*
		****************** 使用 bufio.NewWriter 写入文件 ****************
	*/
	filename := string(time.Now().Format("2006-01-02")) + "-" + strconv.Itoa(int(time.Now().Unix()))
	files, err4 := OpenFile(filename)
	if err4 != nil {
		log.Fatal(err4.Error())
	}
	w := bufio.NewWriter(files) //创建新的 Writer 对象

	for i := 1; i <= allpage; i++ {
		time.Sleep(time.Duration(getrandomtime()) * time.Millisecond) // 每获取底部页面依次，即睡眠500~1000毫秒

		url_append := "s=1&p=" + strconv.Itoa(i) + "&token="
		fmt.Println(url_append)

		second_html := postfunc("https://www.butian.net/Reward/pub", "application/x-www-form-urlencoded", url_append,cookiessend)

		xin := AutoJson{} // 创建新的结构体，用于处理获取的json数据

		errxin := json.Unmarshal([]byte(second_html), &xin) // 新的结构体用于处理获取的json数据
		if errxin != nil {
			log.Print(errxin)
		}

		for x := 0; x < len(xin.AllData.List); x++ {
			n_string, errn_string := w.WriteString("https://www.butian.net/Loo/submit?cid=" + xin.AllData.List[x].CompanyID + "\n")
			if errn_string != nil {
				log.Fatal(errn_string)
				log.Fatal(n_string)
			}
			fmt.Println("获得厂商id对应url为：%s", xin.AllData.List[x].CompanyID)

			all_info = append(all_info, "https://www.butian.net/Loo/submit?cid="+xin.AllData.List[x].CompanyID)
		}
	}

	url_filename := "URL" + string(time.Now().Format("2006-01-02")) + "-" + strconv.Itoa(int(time.Now().Unix()))
	url_files, url_err := OpenFile(url_filename)
	if url_err != nil {
		log.Fatal(url_err.Error())
	}
	url_w := bufio.NewWriter(url_files) //创建新的 Writer 对象

	for in := 0; in < len(all_info); in++ {
		temp_values := webdeal(all_info[in],cookiessend)

		url_n_string, url_errn_string := url_w.WriteString(temp_values + "\n")
		if url_errn_string != nil {
			log.Fatal(url_errn_string)
			log.Fatal(url_n_string)
		}
	}
	url_w.Flush()
	url_files.Close()

	w.Flush()
	files.Close()

	//计算代码运行时间
	cost := time.Since(start)
	fmt.Printf("cost=[%s]",cost)
}