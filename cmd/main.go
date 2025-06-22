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

	maxPages := flag.Int("maxPages", 10, "maximum number of pages to scrape")
	brand := flag.String("brand", "", "car brand ID (optional)")
	flag.Parse()

	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	calc.GetAllCars(ctx, maxPages, brand)
}
