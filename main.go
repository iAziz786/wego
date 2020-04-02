package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

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

func getLink(n *html.Node, anchorLinks chan string) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key != "href" {
				continue
			}
			link, err := url.Parse(a.Val)
			if err != nil {
				// ignore bad urls
				continue
			}
			if link.String() != "" {
				anchorLinks <- link.String()
			}
		}
	}
}

func traverse(startingNode *html.Node, wg *sync.WaitGroup, anchorLinks chan string) {
	getLink(startingNode, anchorLinks)
	for n := startingNode.FirstChild; n != nil; n = n.NextSibling {
		traverse(n, wg, anchorLinks)
	}
}

// GetCrawableURLs will crawl all the elements of the node and find all the
// all the anchor links and pass them down the channel
func GetCrawableURLs(nodeStream <-chan *html.Node) chan string {
	anchorLinks := make(chan string)
	go func() {
		var wg sync.WaitGroup
		defer close(anchorLinks)
		node := <-nodeStream
		traverse(node, &wg, anchorLinks)
	}()
	return anchorLinks
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
	crawlableLinks := GetCrawableURLs(nodeStream)
	for link := range crawlableLinks {
		parse, err := url.Parse(link)
		if err != nil {
			continue
		}
		base, err := url.Parse(crawlURL)
		nextURLToCrawl := base.ResolveReference(parse)
		fmt.Println(nextURLToCrawl.String())
	}
}
