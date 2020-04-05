package traverser

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestTraverser(t *testing.T) {
	reader := strings.NewReader(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Document</title>
	</head>
	<body>
		
	</body>
	</html>`)
	node, err := html.Parse(reader)
	if err != nil {
		t.Error(err)
	}
	str := []string{}
	Traverse(node, func(eachNode *html.Node) {
		if eachNode.Type == html.ElementNode {
			str = append(str, eachNode.Data)
		}
	})
	if strings.Join(str, " ") != "html head meta meta title body" {
		t.Errorf("expected %s\ngot %s", "html head meta meta title body", strings.Join(str, " "))
	}
}
