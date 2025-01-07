package scraper

import (
	"errors"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

func FetchTitleAndMeta(url string) (string, string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", "", err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", "", errors.New("failed to fetch: " + response.Status)
	}

	return extractMeta(response.Body)
}

func extractMeta(body io.Reader) (string, string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", "", err
	}

	var title, description string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "title" && n.FirstChild != nil {
				title = n.FirstChild.Data
			}
			if n.Data == "meta" {
				for _, attr := range n.Attr {
					if attr.Key == "name" && attr.Val == "description" {
						for _, attr := range n.Attr {
							if attr.Key == "content" {
								description = attr.Val
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	if title == "" {
		return "", "", errors.New("title not found")
	}

	return title, description, nil
}

func FetchLinks(url string) ([]string, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch: " + response.Status)
	}

	return extractLinks(response.Body)
}

func extractLinks(body io.Reader) ([]string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return nil, err
	}

	var links []string
	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					links = append(links, attr.Val)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}
	traverse(doc)

	return links, nil
}
