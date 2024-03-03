package tests

import (
	"testing"
)

func TestSendMail(t *testing.T, fun error) {

	price := 100
	discountPercent := 10
	discountedPrice := CalculateDiscount(price, discountPercent)
	expectedPrice := 90
	if discountedPrice != expectedPrice {
		t.Errorf("Incorrect discounted price. Expected: %d, Got: %d", expectedPrice, discountedPrice)
	}
}
