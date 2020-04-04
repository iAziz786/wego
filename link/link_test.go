package link

import (
	"testing"
)

// type Response struct {
// 	Body string
// }

// func (r *Response) Read(p []byte) (n int, err error) {
// 	fmt.Print("reader", p)
// 	return len(r.Body), nil
// }

// func TestGetLink(t *testing.T) {
// 	anchorLinks := make(chan string)
// 	aLink := "https://dev.to"
// 	resp := strings.NewReader("<a href='" + aLink + "'>Hello</a>")

// 	node, err := html.Parse(resp)
// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	go GetLink(node, anchorLinks)

// 	select {
// 	case link := <-anchorLinks:
// 		fmt.Println("received anchor link")
// 		if link != aLink {
// 			t.Error("link didn't match")
// 		}
// 	case <-time.After(5 * time.Second):
// 		t.Error("timeout when running test")
// 	}
// }

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
