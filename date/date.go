package date

import (
	"time"

	"golang.org/x/exp/rand"
)

func GenerateExpireDate() time.Time {
	now := time.Now()
	year := now.Year() + rand.Intn(5)
	month := time.Month(rand.Intn(12) + 1)
	return time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
}
