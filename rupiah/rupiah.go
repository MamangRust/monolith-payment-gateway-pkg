package rupiah

import (
	"fmt"
	"strconv"
)

func RupiahFormat(digit string) string {
	digitNumber, err := strconv.ParseFloat(digit, 64)

	if err != nil {
		return "Rp 0"
	}

	formatter := fmt.Sprintf("Rp.%.0f", digitNumber)
	return formatter
}
