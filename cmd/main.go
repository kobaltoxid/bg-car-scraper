package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/kobaltoxid/bg-car-scraper/internal/calc"
)

func main() {
	fmt.Println("hello")

	maxPages := flag.Int("maxPages", 10, "maximum number of pages to scrape")
	brand := flag.String("brand", "", "car brand ID (optional)")
	model := flag.String("model", "", "car model name (optional)")
	flag.Parse()

	ctx := context.Background()
	calc.GetAllCars(ctx, maxPages, brand, model)
}
