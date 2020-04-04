package link

import (
	"net/url"

	"golang.org/x/net/html"
)

// GetLink parse the dom tree and send the link to anchorLink
func GetLink(n *html.Node, anchorLinks chan string) {
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
	} else {
		anchorLinks <- ""
	}
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
