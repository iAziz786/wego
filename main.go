package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/iAziz786/wego/extract"
	"github.com/iAziz786/wego/link"
	"golang.org/x/net/html"
)

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

type RobotsTxt string

func (r RobotsTxt) Get() RobotsTxt {
	return r
}

func isValidURL(crawlURL string) bool {
	_, err := url.ParseRequestURI(crawlURL)
	if err != nil {
		return false
	}

	u, err := url.Parse(crawlURL)

	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func getBody(crawlURL string, nodeStream chan<- *html.Node) {
	resp, err := http.Get(crawlURL)
	panicIfError(err)
	defer resp.Body.Close()
	node, err := html.Parse(resp.Body)
	panicIfError(err)
	nodeStream <- node
}

func main() {
	crawlURL := os.Args[1]
	fmt.Println(crawlURL)
	if !isValidURL(crawlURL) {
		log.Panic("Invalid URL passed to crawl")
	}
	nodeStream := make(chan *html.Node)
	go getBody(crawlURL, nodeStream)
	crawlableLinks := link.GetCrawableURLs(nodeStream, crawlURL)
	var wg sync.WaitGroup
	for foundLink := range crawlableLinks {
		if foundLink.RelativeHref != "" {
			wg.Add(1)
			joinedLink, _ := link.JoinURLs(foundLink.BaseURL, foundLink.RelativeHref)
			time.Sleep(100 * time.Millisecond)
			go func() {
				defer wg.Done()
				resp, err := http.Get(joinedLink)
				panicIfError(err)
				defer resp.Body.Close()
				node, err := html.Parse(resp.Body)
				panicIfError(err)
				text := extract.GetTextContent(node)
				fmt.Println(text)
			}()
		}
	}
	wg.Wait()
}
