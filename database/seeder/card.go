package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/date"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-pkg/randomvcc"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type cardSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewCardSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *cardSeeder {
	return &cardSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *cardSeeder) Seed() error {
	cardTypes := []string{"credit", "debit"}
	cardProviders := []string{"mandiri", "bni", "bri"}

	generatedCards := make(map[string]struct{})

	totalCards := 40
	activeCards := 20
	trashedCards := 20

	for i := 1; i <= totalCards; i++ {
		var random string
		var err error

		for {
			random, err = randomvcc.RandomCardNumber()
			if err != nil {
				r.logger.Error("failed to generate card number for card", zap.Int("card", i), zap.Error(err))
				return fmt.Errorf("failed to generate card number for card %d: %w", i, err)
			}

			if _, exists := generatedCards[random]; !exists {
				generatedCards[random] = struct{}{}
				break
			}
		}

		request := db.CreateCardParams{
			UserID:       int32(i),
			CardNumber:   random,
			CardType:     cardTypes[i%len(cardTypes)],
			ExpireDate:   date.GenerateExpireDate(),
			Cvv:          fmt.Sprintf("%03d", rand.Intn(1000)),
			CardProvider: cardProviders[i%len(cardProviders)],
		}

		card, err := r.db.CreateCard(r.ctx, request)
		if err != nil {
			r.logger.Error("failed to seed card", zap.Int("card", i), zap.Error(err))
			return fmt.Errorf("failed to seed card %d: %w", i, err)
		}

		if i > activeCards {
			_, err = r.db.TrashCard(r.ctx, card.CardID)
			if err != nil {
				r.logger.Error("failed to trash card", zap.Int("card", i), zap.Error(err))
				return fmt.Errorf("failed to trash card %d: %w", i, err)
			}
		}
	}

	r.logger.Debug("card seeded successfully", zap.Int("totalCards", totalCards), zap.Int("activeCards", activeCards), zap.Int("trashedCards", trashedCards))

	return nil
}
