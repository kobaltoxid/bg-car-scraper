package calc

import (
	"bg-carinator/internal/offers"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func GetAllCarsUpToGivenPage(ctx context.Context, maxPages *int) {

	var html string
	var allOffers []offers.Offer

	offers, err := offers.ExtractOffers(html)
	if err != nil {
		log.Fatal("Error parsing offers:", err)
	}

	allOffers = append(allOffers, offers...)

	for page := 1; page <= *maxPages; page++ {
		var pageHTML string
		url := fmt.Sprintf("https://www.cars.bg/carslist.php?page=%d", page)
		fmt.Println("Scraping:", url)

		selector := fmt.Sprintf("div.mdc-layout-grid__inner.white-background.pageContainer.page-%d", page)

		err := chromedp.Run(ctx,
			chromedp.Navigate(url),
			chromedp.Sleep(1000*time.Millisecond),
			chromedp.OuterHTML(
				selector,
				&pageHTML,
			),
		)
		if err != nil {
			log.Printf("Error on page %d: %v", page, err)
			continue
		}

		// 5. Print all offers
		for _, offer := range allOffers {
			fmt.Printf("data-item: %s\n", offer.DataItem)
			fmt.Printf("title: %s\n", offer.Title)
			fmt.Printf("image url: %s\n", offer.ImageURL)
			fmt.Printf("list-link: %s\n", offer.ListLink)
			fmt.Println("-----")
		}
	}
}
