package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	res, err := http.Get("https://dev.to")
	panicIfError(err)
	defer res.Body.Close()
	z := html.NewTokenizer(res.Body)
	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			return
		case html.TextToken:
			fmt.Println(string(z.Text()))
		}
	}
}
