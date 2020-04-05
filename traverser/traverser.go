package traverser

import "golang.org/x/net/html"

type nodeHander func(*html.Node)

func Traverse(startingNode *html.Node, fn nodeHander) {
	fn(startingNode)
	for n := startingNode.FirstChild; n != nil; n = n.NextSibling {
		Traverse(n, fn)
	}
}
