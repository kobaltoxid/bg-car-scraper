package main

import (
	"context"
	"flag"
	"fmt"

	"bg-carinator/internal/calc"

	"github.com/chromedp/chromedp"
)

func main() {
	fmt.Println("hello")

	maxPages := flag.Int("maxPages", 100, "maximum number of pages to scrape")
	flag.Parse()

	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	calc.GetAllCarsUpToGivenPage(ctx, maxPages)
}
