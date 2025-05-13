package main

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
)

var depth int

func main() {
	node, err := parseUrl()
	if err != nil {
		log.Fatal(err)
	}

	prettify(node)
}

func prettify(n *html.Node) {
	forEachNode(n, startElement, endElement)
}

func startElement(n *html.Node) {
	if n.Type == html.ElementNode {
		fmt.Printf("%*s<%s", depth*2, "", n.Data)
		for _, attr := range n.Attr {
			fmt.Printf(" %s='%s'", attr.Key, attr.Val)
		}
		if n.FirstChild != nil {
			fmt.Println(">")
		} else {
			fmt.Println(" />")
		}
		depth++
	} else if n.Type == html.TextNode {
		if n.Parent != nil && n.Parent.Data == "script" {
			printJS(n)
		} else {
			fmt.Printf("%*s%s\n", depth*2, "", n.Data)
		}
	}
}

func printJS(n *html.Node) {
	fmt.Printf("%*s", depth*2, "")
	runes := []rune(n.Data)
	disabledFmt := false
	for i, r := range runes {
		if r == '\'' {
			fmt.Print(string(r))
			disabledFmt = !disabledFmt
			continue
		}

		if disabledFmt {
			fmt.Print(string(r))
			continue
		}

		if r == '{' {
			fmt.Println(string(r))
			depth++
			fmt.Printf("%*s", depth*2, "")
		} else if r == '}' {
			fmt.Println()
			depth--
			fmt.Printf("%*s", depth*2, "")
			fmt.Print(string(r))
		} else if r == ';' && i < len(runes)-1 && runes[i+1] != '}' {
			fmt.Println(string(r))
			fmt.Printf("%*s", depth*2, "")
		} else {
			fmt.Print(string(r))
		}
	}
	fmt.Println()
}

func endElement(n *html.Node) {
	if n.Type == html.ElementNode {
		depth--
		if n.FirstChild != nil {
			fmt.Printf("%*s</%s>\n", depth*2, "", n.Data)
		}
	}
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}
