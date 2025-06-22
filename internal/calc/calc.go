package calc

import (
	"bg-carinator/internal/offers"
	"bg-carinator/internal/searchcriteria"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
)

func calcUrl(brand searchcriteria.CarBrand, page int) string {
	var url string
	if brand != 0 {
		url = fmt.Sprintf("https://www.cars.bg/carslist.php?subm=1&add_search=1&typeoffer=1&brandId=%d&page=%d", brand, page)
	} else {
		url = fmt.Sprintf("https://www.cars.bg/carslist.php?page=%d", page)
	}
	return url
}

func GetAllCars(ctx context.Context, maxPages *int, brand *string) {
	var allOffers []offers.Offer

	carBrandId := searchcriteria.BrandNameToID(*brand)

	for page := 1; page <= *maxPages; page++ {
		url := calcUrl(carBrandId, page)
		fmt.Println("Scraping:", url)

		offersOnPage, err := offers.GetOffersByUrl(ctx, url)
		if err != nil {
			log.Printf("Error parsing offers on page %d: %v", page, err)
			continue
		}
		allOffers = append(allOffers, offersOnPage...)
	}

	// Print all offers
	for _, offer := range allOffers {
		fmt.Printf("data-item: %s\n", offer.DataItem)
		fmt.Printf("title: %s\n", offer.Title)
		fmt.Printf("image url: %s\n", offer.ImageURL)
		fmt.Printf("list-link: %s\n", offer.ListLink)
		fmt.Println("-----")
	}
}

func GetAllCarsFromBrand(ctx context.Context, maxPages *int) {
	var allOffers []offers.Offer

	for page := 1; page <= *maxPages; page++ {
		url := fmt.Sprintf("https://www.cars.bg/carslist.php?page=%d", page)
		fmt.Println("Scraping:", url)

		// Fetch HTML using net/http
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			log.Printf("Error creating request for page %d: %v", page, err)
			continue
		}
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error fetching page %d: %v", page, err)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Printf("Error reading body for page %d: %v", page, err)
			continue
		}

		offersOnPage, err := offers.ExtractAllOffers(string(body))
		if err != nil {
			log.Printf("Error parsing offers on page %d: %v", page, err)
			continue
		}
		allOffers = append(allOffers, offersOnPage...)
	}

	// Print all offers
	for _, offer := range allOffers {
		fmt.Printf("data-item: %s\n", offer.DataItem)
		fmt.Printf("title: %s\n", offer.Title)
		fmt.Printf("image url: %s\n", offer.ImageURL)
		fmt.Printf("list-link: %s\n", offer.ListLink)
		fmt.Println("-----")
	}
}
