package offers

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

type Offer struct {
	DataItem string
	Title    string
	ImageURL string
	ListLink string
}

func ExtractOffers(htmlStr string) ([]Offer, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return nil, err
	}

	var offers []Offer
	var f func(*html.Node)
	f = func(n *html.Node) {
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
				}
				var findDetails func(*html.Node)
				findDetails = func(nn *html.Node) {
					if nn.Type == html.ElementNode && nn.Data == "div" {
						for _, attr := range nn.Attr {
							if attr.Key == "style" && strings.Contains(attr.Val, "background-image") {
								re := regexp.MustCompile(`url\(['"]?([^'")]+)['"]?\)`)
								matches := re.FindStringSubmatch(attr.Val)
								if len(matches) > 1 {
									offer.ImageURL = matches[1]
								}
							}
						}
					}
					if nn.Type == html.ElementNode && nn.Data == "a" {
						for _, attr := range nn.Attr {
							if attr.Key == "list-link" {
								offer.ListLink = attr.Val
							}
						}
					}
					for c := nn.FirstChild; c != nil; c = c.NextSibling {
						findDetails(c)
					}
				}
				findDetails(n)
				offers = append(offers, offer)
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return offers, nil
}
