package utils

import (
	"github.com/codevault-llc/minerva/pkg/types"
	"golang.org/x/net/html"
)

// ExtractTitle retrieves the title from the parsed HTML.
func ExtractTitle(doc *html.Node) string {
	var title string
	TraverseHTML(doc, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "title" && node.FirstChild != nil {
			title = node.FirstChild.Data
		}
	})
	return title
}

// ProcessScriptNode extracts data from a script element.
func ProcessScriptNode(node *html.Node) types.FileRequest {
	// Handle inline script.
	if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		return CreateFileRequest("inline-script", node.FirstChild.Data, "application/javascript")
	}
	return types.FileRequest{}
}

// ProcessStyleNode extracts data from a style element.
func ProcessStyleNode(node *html.Node) types.FileRequest {
	if node.FirstChild != nil && node.FirstChild.Type == html.TextNode {
		return CreateFileRequest("inline-style", node.FirstChild.Data, "text/css")
	}
	return types.FileRequest{}
}

// ProcessLinkNode extracts data from a link element if it's a stylesheet.
func ProcessLinkNode(node *html.Node) types.FileRequest {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return types.FileRequest{
				Src:      attr.Val,
				FileType: "text/css",
			}
		}
	}
	return types.FileRequest{}
}

// ProcessFontNode extracts data from a link element if it's a font.
func ProcessFontNode(node *html.Node) types.FileRequest {
	for _, attr := range node.Attr {
		if attr.Key == "href" {
			return types.FileRequest{
				Src:      attr.Val,
				FileType: "font",
			}
		}
	}
	return types.FileRequest{}
}

// CreateFileRequest constructs a FileRequest with content details.
func CreateFileRequest(src, content, fileType string) types.FileRequest {
	return types.FileRequest{
		Src:        src,
		Content:    content,
		HashedBody: SHA256(content),
		FileSize:   uint(len(content)),
		FileType:   fileType,
	}
}

// TraverseHTML recursively walks through the HTML nodes and applies the given function.
func TraverseHTML(node *html.Node, fn func(*html.Node)) {
	fn(node)
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		TraverseHTML(child, fn)
	}
}
