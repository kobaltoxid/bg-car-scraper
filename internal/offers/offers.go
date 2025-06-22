package offers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Offer struct {
	DataItem string
	Title    string
	ImageURL string
	ListLink string
}

func ExtractAllOffers(htmlStr string) ([]Offer, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	var offers []Offer

	var findOffers func(*html.Node)
	findOffers = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "div" {
			var dataItem, title string
			for _, attr := range n.Attr {
				if attr.Key == "data-item" {
					dataItem = attr.Val
				}
				if attr.Key == "title" {
					title = attr.Val
				}
			}
			if dataItem != "" {
				offer := Offer{
					DataItem: dataItem,
					Title:    title,
					ImageURL: findImageURL(n),
					ListLink: findListLink(n),
				}
				offers = append(offers, offer)
				// Do not recurse further into this offer node
				return
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			findOffers(c)
		}
	}
	findOffers(doc)
	return offers, nil
}

// Helper to find image URL in the subtree
func findImageURL(n *html.Node) string {
	var imgURL string
	var f func(*html.Node)
	f = func(nn *html.Node) {
		if nn.Type == html.ElementNode && nn.Data == "div" {
			for _, attr := range nn.Attr {
				if attr.Key == "style" && strings.Contains(attr.Val, "background-image") {
					re := regexp.MustCompile(`url\(['"]?([^'")]+)['"]?\)`)
					matches := re.FindStringSubmatch(attr.Val)
					if len(matches) > 1 {
						imgURL = matches[1]
						return
					}
				}
			}
		}
		for c := nn.FirstChild; c != nil; c = c.NextSibling {
			if imgURL == "" {
				f(c)
			}
		}
	}
	f(n)
	return imgURL
}

// Helper to find list-link in the subtree
func findListLink(n *html.Node) string {
	var link string
	var f func(*html.Node)
	f = func(nn *html.Node) {
		if nn.Type == html.ElementNode && nn.Data == "a" {
			for _, attr := range nn.Attr {
				if attr.Key == "list-link" {
					link = attr.Val
					return
				}
			}
		}
		for c := nn.FirstChild; c != nil; c = c.NextSibling {
			if link == "" {
				f(c)
			}
		}
	}
	f(n)
	return link
}

func GetOffersByUrl(ctx context.Context, url string) ([]Offer, error) {
	// Create a request with timeout using the provided context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Save HTML to a file for debugging
	err = os.WriteFile("debug.html", body, 0644)
	if err != nil {
		fmt.Println("Failed to write debug.html:", err)
	}
	return ExtractAllOffers(string(body))
}
