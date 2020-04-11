package extract

import (
	"github.com/iAziz786/wego/traverser"
	"golang.org/x/net/html"
)

func prepareTagExclusion(excludedTags []string) func(string) bool {
	tagMaps := make(map[string]bool)
	for _, tag := range excludedTags {
		tagMaps[tag] = true
	}
	return func(tag string) bool {
		return !tagMaps[tag]
	}
}

// GetTextContent is like textContent property that you will find in your browsers
// DevTool on any DOM node. It extracts the text and concate all of them and return it
func GetTextContent(node *html.Node) (txt string) {
	traverser.Traverse(node, func(eachNode *html.Node) {
		exclude := prepareTagExclusion([]string{"script", "style"})
		if eachNode.Type == html.TextNode && eachNode.Parent != nil && exclude(eachNode.Parent.Data) {
			txt += eachNode.Data
		}
	})
	return txt
}
