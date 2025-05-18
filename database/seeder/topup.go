package seeder

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.uber.org/zap"
)

type topupSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewTopupSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *topupSeeder {
	return &topupSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *topupSeeder) Seed() error {
	totalTopups := 40
	activeTopups := 20
	trashedTopups := 20

	cards, err := r.db.GetCards(r.ctx, db.GetCardsParams{
		Column1: "",
		Limit:   int32(totalTopups),
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("failed to get card list", zap.Error(err))
		return fmt.Errorf("failed to get card list: %w", err)
	}

	if len(cards) == 0 {
		r.logger.Error("no cards available for topup seeding")
		return fmt.Errorf("no cards available for topup seeding")
	}

	topupMethods := []string{"Bank Alpha", "Bank Beta", "Bank Gamma"}
	statusOptions := []string{"pending", "success", "failed"}

	months := make([]time.Time, 12)
	currentYear := time.Now().Year()
	for i := 0; i < 12; i++ {
		months[i] = time.Date(currentYear, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
	}

	for i := 0; i < totalTopups; i++ {
		card := cards[i%len(cards)]
		cardNumber := card.CardNumber

		if len(cardNumber) < 5 {
			r.logger.Error("card number is too short", zap.String("card", cardNumber))
			return fmt.Errorf("card number %s is too short", cardNumber)
		}

		monthIndex := i % 12
		topupTime := months[monthIndex].Add(time.Duration(rand.Intn(28)) * 24 * time.Hour)

		request := db.CreateTopupParams{
			CardNumber:  cardNumber,
			TopupAmount: int32(rand.Intn(10000000) + 1000000),
			TopupMethod: topupMethods[rand.Intn(len(topupMethods))],
			TopupTime:   topupTime,
		}

		topup, err := r.db.CreateTopup(r.ctx, request)
		if err != nil {
			r.logger.Error("failed to seed topup for card", zap.String("card", cardNumber), zap.Error(err))
			return fmt.Errorf("failed to seed topup for card %s: %w", cardNumber, err)
		}

		status := statusOptions[rand.Intn(len(statusOptions))]
		_, err = r.db.UpdateTopupStatus(r.ctx, db.UpdateTopupStatusParams{
			TopupID: topup.TopupID,
			Status:  status,
		})
		if err != nil {
			r.logger.Error("failed to update topup status", zap.Int("topupID", int(topup.TopupID)), zap.String("status", status), zap.Error(err))
			return fmt.Errorf("failed to update topup status for topup ID %d: %w", topup.TopupID, err)
		}

		if i >= activeTopups {
			_, err = r.db.TrashTopup(r.ctx, topup.TopupID)
			if err != nil {
				r.logger.Error("failed to trash topup", zap.Int("topup", i+1), zap.String("card", cardNumber), zap.Error(err))
				return fmt.Errorf("failed to trash topup %d for card %s: %w", i+1, cardNumber, err)
			}
		}
	}

	r.logger.Debug("topup seeded successfully", zap.Int("totalTopups", totalTopups), zap.Int("activeTopups", activeTopups), zap.Int("trashedTopups", trashedTopups))

	return nil
}
