package link

import (
	"strings"
	"testing"

	"github.com/iAziz786/wego/traverser"
	"golang.org/x/net/html"
)

// type Response struct {
// 	Body string
// }

// func (r *Response) Read(p []byte) (n int, err error) {
// 	fmt.Print("reader", p)
// 	return len(r.Body), nil
// }

func TestGetLink(t *testing.T) {
	aLink := "https://dev.to"
	resp := strings.NewReader("<a href='" + aLink + "'>Hello</a>")

	node, err := html.Parse(resp)
	if err != nil {
		t.Error(err.Error())
	}

	var anchorNode *html.Node
	traverser.Traverse(node, func(eachNode *html.Node) {
		if eachNode.Data == "a" {
			anchorNode = eachNode
		}
	})
	linkFromHref := GetLink(anchorNode)

	if linkFromHref != aLink {
		t.Errorf("expected %s got %s", aLink, linkFromHref)
	}
}

func TestJoinURLs(t *testing.T) {
	withTrailingSlash, _ := JoinURLs("https://example.com", "/next-page")
	expected := "https://example.com/next-page"
	if withTrailingSlash != expected {
		t.Errorf("expected %s got %s", expected, withTrailingSlash)
	}
	withoutTrailingSlash, _ := JoinURLs("https://example.com/first", "first/second")
	expected = "https://example.com/first/second"
	if withoutTrailingSlash != expected {
		t.Errorf("expected %s got %s", expected, withoutTrailingSlash)
	}
}

func TestGetCrawableURLsWithBaseURL(t *testing.T) {
	reader := strings.NewReader(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<base href="https://different-url.com" />
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>
		<a href="/next">Go to next of different url</a>
	</body>
	</html>
	`)
	locationHref := "https://example.com"
	nodeStream := make(chan *html.Node)
	foundLinkStream := GetCrawableURLs(nodeStream, locationHref)
	node, err := html.Parse(reader)
	if err != nil {
		t.Error(err)
	}

	select {
	case foundLink := <-foundLinkStream:
		if foundLink.BaseURL == locationHref {
			t.Error(locationHref, "should not be base URL")
		}
		if !(foundLink.BaseURL == "https://different-url.com" && foundLink.RelativeHref == "/next") {
			t.Errorf("expected %s got %s", "https://different-url.com", foundLink.BaseURL)
			t.Errorf("expected %s got %s", "/next", foundLink.RelativeHref)
		}
	case nodeStream <- node:
	}
}

func TestGetCrawableURLsWithAbsolutePath(t *testing.T) {
	reader := strings.NewReader(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<base href="https://different-url.com" />
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>
		<a href="https://example.com/next">Go to next of different url</a>
	</body>
	</html>
	`)
	locationHref := "https://example.com"
	nodeStream := make(chan *html.Node)
	foundLinkStream := GetCrawableURLs(nodeStream, locationHref)
	node, err := html.Parse(reader)
	if err != nil {
		t.Error(err)
	}

	select {
	case foundLink := <-foundLinkStream:
		joinedURL, err := JoinURLs(foundLink.BaseURL, foundLink.RelativeHref)
		if err != nil {
			t.Error(err)
		}
		if joinedURL != "https://example.com/next" {
			t.Errorf("expected %s got %s", "https://example.com/next", joinedURL)
		}
	case nodeStream <- node:
	}
}
