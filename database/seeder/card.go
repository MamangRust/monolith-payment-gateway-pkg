package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/date"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/MamangRust/monolith-payment-gateway-pkg/randomvcc"
	"go.uber.org/zap"
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

	totalCards := 10
	activeCards := 5
	trashedCards := 5

	cardNumbers := make([]string, totalCards)
	for i := 0; i < totalCards; i++ {
		cardNumber, err := randomvcc.RandomCardNumber()
		if err != nil {
			r.logger.Error("failed to generate card number", zap.Int("index", i), zap.Error(err))
			return fmt.Errorf("failed to generate card number: %w", err)
		}
		cardNumbers[i] = cardNumber
	}

	for i := 0; i < totalCards; i++ {
		request := db.CreateCardParams{
			UserID:       int32(i + 1),
			CardNumber:   cardNumbers[i],
			CardType:     cardTypes[i%len(cardTypes)],
			ExpireDate:   date.GenerateExpireDate(),
			Cvv:          fmt.Sprintf("%03d", i%1000),
			CardProvider: cardProviders[i%len(cardProviders)],
		}

		card, err := r.db.CreateCard(r.ctx, request)
		if err != nil {
			r.logger.Error("failed to seed card", zap.Int("card", i+1), zap.Error(err))
			return fmt.Errorf("failed to seed card %d: %w", i+1, err)
		}

		if i >= activeCards {
			_, err = r.db.TrashCard(r.ctx, card.CardID)
			if err != nil {
				r.logger.Error("failed to trash card", zap.Int("card", i+1), zap.Error(err))
				return fmt.Errorf("failed to trash card %d: %w", i+1, err)
			}
		}
	}

	r.logger.Info("card seeded successfully", zap.Int("totalCards", totalCards), zap.Int("activeCards", activeCards), zap.Int("trashedCards", trashedCards))
	return nil
}
