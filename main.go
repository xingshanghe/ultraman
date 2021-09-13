package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
	"github.com/gocolly/colly/v2/extensions"
	"github.com/ztrue/tracerr"
	"golang.org/x/sync/errgroup"
)

var g errgroup.Group

const PATH_CARDS = "./"

func main() {
	g.Go(func() error {
		return archive("http://archive.dcd-ultraman.com.cn/index.php?g=portal&m=list&a=index&id=2")
	})
	g.Go(func() error {
		return blazing("http://www.dcd-ultraman.com.cn/list.php?pid=4&ty=24")
	})
	g.Go(func() error {
		return blazing("http://www.dcd-ultraman.com.cn/list.php?pid=4&ty=10")
	})
	if err := g.Wait(); err != nil {
		tracerr.PrintSourceColor(err)
	}
}

// genColly
// @Description:
// @Date: 2021-09-13 16:40:38
// @return *colly.Collector
//
func genColly() *colly.Collector {
	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))
	extensions.RandomUserAgent(c)
	extensions.Referer(c)
	return c
}

// archive
// @Description: 往期
// @Date: 2021-09-13 16:39:51
// @param url
// @return error
//
func archive(url string) error {
	var err error
	c := genColly()
	// 一共56页
	c.OnHTML("#mainCol-inner", func(e *colly.HTMLElement) {
		e.ForEach(" .cardCol", func(i int, ee *colly.HTMLElement) {
			name := ee.ChildText("h3.card-name")
			number := ee.ChildText("p.card-num")
			rarity := ee.ChildText("p.rarity")
			src := ee.ChildAttr("p.card-img img", "src")
			ext := path.Ext(src)
			rsp, errInner := http.Get(src)
			if errInner != nil {
				err = errInner
				fmt.Println("err:", err)
			}
			defer rsp.Body.Close()
			body, errInner := ioutil.ReadAll(rsp.Body)
			if errInner != nil {
				err = errInner
				fmt.Println("err:", err)
			}
			_, errInner = os.Stat(PATH_CARDS + rarity)
			if errInner != nil {
				if os.IsNotExist(errInner) {
					errInner = os.Mkdir(PATH_CARDS+rarity, os.ModePerm)
				}
			}
			if errInner != nil {
				err = errInner
				fmt.Println("err:", err)
			}
			errInner = ioutil.WriteFile(PATH_CARDS+rarity+"/"+name+"-"+number+ext, body, 0755)
			if errInner != nil {
				err = errInner
				fmt.Println("err:", err)
			}
		})

		if next, found := e.DOM.Find("div.paginator li:nth-last-child(2) a").Attr("href"); found {
			e.Request.Visit(e.Request.AbsoluteURL(next))
		} else {
			return
		}

	})
	return c.Visit(url)
}

// blazing
// @Description: 炽热
// @Date: 2021-09-13 16:40:13
// @return error
//
func blazing(url string) error {
	var err error
	c := genColly()
	// 一共56页
	c.OnHTML("#mainCol", func(e *colly.HTMLElement) {
		e.ForEach("#list .cardCol", func(i int, ee *colly.HTMLElement) {
			name := ee.ChildText("h5.cardName")
			number := ee.ChildText("p.cardNum")
			rarity := ee.ChildText("p.rarity")
			src := ee.ChildAttr("div.cardImg img", "src")
			ext := path.Ext(src)
			rsp, errInner := http.Get(ee.Request.AbsoluteURL(src))
			if errInner != nil {
				err = errInner
				fmt.Println("err:", err)
			}
			defer rsp.Body.Close()
			body, errInner := ioutil.ReadAll(rsp.Body)
			if errInner != nil {
				err = errInner
				fmt.Println("err:", err)
			}
			_, errInner = os.Stat(PATH_CARDS + rarity)
			if errInner != nil {
				if os.IsNotExist(errInner) {
					errInner = os.Mkdir(PATH_CARDS+rarity, os.ModePerm)
				}
			}
			if errInner != nil {
				err = errInner
				fmt.Println("err:", err)
			}
			errInner = ioutil.WriteFile(PATH_CARDS+rarity+"/"+name+"-"+number+ext, body, 0755)
			if errInner != nil {
				err = errInner
				fmt.Println("err:", err)
			}
		})

		if next, found := e.DOM.Find("div.page > form > a:last-child").Attr("href"); found {
			e.Request.Visit(e.Request.AbsoluteURL(next))
		} else {
			return
		}

	})
	return c.Visit(url)
}
