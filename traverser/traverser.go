package traverser

import "golang.org/x/net/html"

type nodeHander func(*html.Node)

// Traverse loop through each node and passed it to callback function
func Traverse(startingNode *html.Node, fn nodeHander) {
	for n := startingNode.FirstChild; n != nil; {
		fn(n)
		if n.FirstChild != nil {
			n = n.FirstChild
			continue
		}
		if n.NextSibling != nil {
			n = n.NextSibling
			continue
		}
		n = n.Parent.NextSibling
	}
}
