package extract

import (
	"github.com/iAziz786/wego/traverser"
	"golang.org/x/net/html"
)

// GetTextContent is like textContent property that you will find in your browsers
// DevTool on any DOM node. It extracts the text and concate all of them and return it
func GetTextContent(node *html.Node) (txt string) {
	traverser.Traverse(node, func(eachNode *html.Node) {
		if eachNode.Type == html.TextNode && eachNode.Parent != nil && eachNode.Parent.Data != "script" {
			txt += eachNode.Data
		}
	})
	return txt
}
