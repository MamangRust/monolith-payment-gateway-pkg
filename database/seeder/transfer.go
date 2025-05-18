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

type transferSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewTransferSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *transferSeeder {
	return &transferSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *transferSeeder) Seed() error {
	cards, err := r.db.GetCards(r.ctx, db.GetCardsParams{
		Column1: "",
		Limit:   10,
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("failed to get card list", zap.Error(err))
		return fmt.Errorf("failed to get card list: %w", err)
	}

	if len(cards) < 2 {
		r.logger.Error("not enough cards available for transfer seeding")
		return fmt.Errorf("not enough cards available for transfer seeding")
	}

	statusOptions := []string{"pending", "success", "failed"}

	months := make([]time.Time, 12)
	currentYear := time.Now().Year()
	for i := 0; i < 12; i++ {
		months[i] = time.Date(currentYear, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
	}

	for i := 0; i < 40; i++ {
		fromIndex := rand.Intn(len(cards))
		toIndex := rand.Intn(len(cards))

		for fromIndex == toIndex {
			toIndex = rand.Intn(len(cards))
		}

		transferFrom := cards[fromIndex].CardNumber
		transferTo := cards[toIndex].CardNumber
		transferAmount := int32(rand.Intn(1000000) + 50000)

		status := statusOptions[rand.Intn(len(statusOptions))]

		monthIndex := i % 12
		transferTime := months[monthIndex].Add(time.Duration(rand.Intn(28)) * 24 * time.Hour)

		request := db.CreateTransferParams{
			TransferFrom:   transferFrom,
			TransferTo:     transferTo,
			TransferAmount: transferAmount,
			TransferTime:   transferTime,
			Status:         status,
		}

		transfer, err := r.db.CreateTransfer(r.ctx, request)
		if err != nil {
			r.logger.Error("failed to seed transfer", zap.Int("transfer", i+1), zap.Error(err))
			return fmt.Errorf("failed to seed transfer %d: %w", i+1, err)
		}

		if i < 20 {
			_, err = r.db.TrashTransfer(r.ctx, transfer.TransferID)
			if err != nil {
				r.logger.Error("failed to trash transfer", zap.Int("transfer", i+1), zap.Error(err))
				return fmt.Errorf("failed to trash transfer %d: %w", i+1, err)
			}
		}
	}

	r.logger.Info("transfer seeded successfully")

	return nil
}
