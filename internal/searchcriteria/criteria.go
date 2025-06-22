package searchcriteria

import "strings"

type CarBrand int

const (
	BrandUnknown CarBrand = 0
	BrandBMW     CarBrand = 10
	BrandAudi    CarBrand = 20
	BrandVW      CarBrand = 30
	// Add more brands as needed
)

// BrandNameToID maps a string (case-insensitive) to a CarBrand type.
func BrandNameToID(brand string) CarBrand {
	switch strings.ToLower(brand) {
	case "bmw":
		return BrandBMW
	case "audi":
		return BrandAudi
	case "vw", "volkswagen":
		return BrandVW
	// Add more cases as needed
	default:
		return BrandUnknown
	}
}
