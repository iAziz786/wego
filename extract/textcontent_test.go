package extract

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestGetTextContent(t *testing.T) {
	reader := strings.NewReader(`
<p>Please add this to the main function and check the output.</p> <p>Let's conclude this section by writing one more program. This program will perform the same operations on each element of a slice and return the result. For example if we want to multiply all integers in a slice by 5 and return the output, it can be easily done using first class functions. These kind of functions which operate on every element of a collection are called <code>map</code> functions. I have provided the program below. It is self explanatory.</p>`)
	node, err := html.Parse(reader)
	if err != nil {
		t.Error(err)
	}
	text := GetTextContent(node)
	if text != "Please add this to the main function and check the output. Let's conclude this section by writing one more program. This program will perform the same operations on each element of a slice and return the result. For example if we want to multiply all integers in a slice by 5 and return the output, it can be easily done using first class functions. These kind of functions which operate on every element of a collection are called map functions. I have provided the program below. It is self explanatory." {
		t.Errorf("expected %s got %s", "Please add this to the main function and check the output. Let's conclude this section by writing one more program. This program will perform the same operations on each element of a slice and return the result. For example if we want to multiply all integers in a slice by 5 and return the output, it can be easily done using first class functions. These kind of functions which operate on every element of a collection are called map functions. I have provided the program below. It is self explanatory.", text)
	}
	reader = strings.NewReader(`
<div>many<div>level<div>nested<div>content</div></div></div></div>`)
	node, err = html.Parse(reader)
	if err != nil {
		t.Error(err)
	}
	text = GetTextContent(node)
	if text != "manylevelnestedcontent" {
		t.Errorf("expected %s got %s", "manylevelnestedcontent", text)
	}
}

func TestGetTextContentScriptTag(t *testing.T) {
	reader := strings.NewReader(`
<p>the script <strong>should not</strong> be includes</p><script>var go = "awesome";</script><script src="https://dns.com" />`)
	node, err := html.Parse(reader)
	if err != nil {
		t.Error(err)
	}
	text := GetTextContent(node)
	if text != "the script should not be includes" {
		t.Errorf("expected %s got %s", "the script should not be includes", text)
	}
}
