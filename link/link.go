package link

import (
	"net/url"

	"github.com/iAziz786/wego/traverser"
	"golang.org/x/net/html"
)

// FoundURL is the struct of the base URL of the value of href
type FoundURL struct {
	BaseURL      string
	RelativeHref string
}

// GetLink parse the dom node and return the href if it's a link
func GetLink(n *html.Node) string {
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
				return link.String()
			}
		}
	}
	return ""
}

// JoinURLs with accept two URLs and append them so that it can be feedback
// to the crawler so it can crawl continuously
// Example: https://example.com, /next-page
// Both will be joined and create https://example.com/next-page
func JoinURLs(baseURL, hyperlink string) (string, error) {
	parse, err := url.Parse(hyperlink)
	if err != nil {
		return "", err
	}
	base, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	nextURLToCrawl := base.ResolveReference(parse)
	return nextURLToCrawl.String(), nil
}

// GetCrawableURLs will crawl all the elements of the node and find
// all the anchor links and pass them down the channel
func GetCrawableURLs(nodeStream <-chan *html.Node, locationHref string) chan *FoundURL {
	anchorLinks := make(chan *FoundURL)
	go func() {
		defer close(anchorLinks)
		node := <-nodeStream
		var baseURL string
		traverser.Traverse(node, func(node *html.Node) {
			if node.Data == "base" {
				for _, attr := range node.Attr {
					if attr.Key != "href" {
						continue
					}
					baseURL = attr.Val
				}
			}
			href := GetLink(node)
			if baseURL == "" {
				baseURL = locationHref
			}
			anchorLinks <- &FoundURL{
				BaseURL:      baseURL,
				RelativeHref: href,
			}
		})
	}()
	return anchorLinks
}
