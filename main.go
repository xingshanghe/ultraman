package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
	"io/ioutil"
	"net/http"
	"path"
	"sync"
	"os"
)

var wg sync.WaitGroup


func main() {
	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	// 一共56页
	c.OnHTML("#mainCol-inner", func(e *colly.HTMLElement) {
		e.ForEach(" .cardCol", func(i int, ee *colly.HTMLElement) {
			name := ee.ChildText("h3.card-name")
			number := ee.ChildText("p.card-num")
			rarity := ee.ChildText("p.rarity")
			url := ee.ChildAttr("p.card-img img", "src")
			ext := path.Ext(url)
			rsp, err := http.Get(url)
			if err != nil {
				fmt.Println("err:", err)
			}
			defer rsp.Body.Close()
			body, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				fmt.Println("err:", err)
			}
			_, err = os.Stat("./cards/"+rarity)
			if err != nil {
				if os.IsNotExist(err) {
					err = os.Mkdir("./cards/"+rarity, os.ModePerm)
				}
			}
			if err != nil {
				fmt.Println("err:", err)
			}
			err = ioutil.WriteFile("./cards/"+rarity+"/"+name+"-"+number+ext, body, 0755)
			if err != nil {
				fmt.Println("err:", err)
			}
		})
		next := e.ChildAttr("div.paginator li:nth-last-child(2) a", "href")

		c.Visit(e.Request.AbsoluteURL(next))
	})
	c.Visit("http://www.dcd-ultraman.com.cn/index.php?m=list&a=index&id=2")
}
