package methodtopup

import "strings"

func PaymentMethodValidator(paymentMethod string) bool {
	paymentRules := []string{
		"alfamart",
		"indomart",
		"lawson",
		"dana",
		"ovo",
		"gopay",
		"linkaja",
		"jenius",
		"fastpay",
		"kudo",
		"bri",
		"mandiri",
		"bca",
		"bni",
		"bukopin",
		"e-banking",
		"visa",
		"mastercard",
		"discover",
		"american express",
		"paypal",
	}

	paymentMethodLower := strings.ToLower(paymentMethod)
	for _, rule := range paymentRules {
		if paymentMethodLower == rule {
			return true
		}
	}

	return false
}
