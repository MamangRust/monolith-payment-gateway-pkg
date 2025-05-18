package seeder

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	db "github.com/MamangRust/payment-gateway-monolith-grpc/pkg/database/schema"
	"github.com/MamangRust/payment-gateway-monolith-grpc/pkg/logger"
	"go.uber.org/zap"
)

type withdrawSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewWithdrawSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *withdrawSeeder {
	return &withdrawSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *withdrawSeeder) Seed() error {
	cards, err := r.db.GetCards(r.ctx, db.GetCardsParams{
		Column1: "",
		Limit:   10,
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("failed to get card list", zap.Error(err))
		return fmt.Errorf("failed to get card list: %w", err)
	}

	if len(cards) == 0 {
		r.logger.Error("no cards available for withdraw seeding")
		return fmt.Errorf("no cards available for withdraw seeding")
	}

	statusOptions := []string{"pending", "success", "failed"}

	months := make([]time.Time, 12)
	currentYear := time.Now().Year()
	for i := 0; i < 12; i++ {
		months[i] = time.Date(currentYear, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
	}

	for i := 0; i < 40; i++ {
		selectedCard := cards[rand.Intn(len(cards))]

		status := statusOptions[rand.Intn(len(statusOptions))]

		monthIndex := i % 12
		withdrawTime := months[monthIndex].Add(time.Duration(rand.Intn(28)) * 24 * time.Hour)

		request := db.CreateWithdrawParams{
			CardNumber:     selectedCard.CardNumber,
			WithdrawAmount: int32(rand.Intn(1000000) + 50000),
			WithdrawTime:   withdrawTime,
		}

		withdraw, err := r.db.CreateWithdraw(r.ctx, request)
		if err != nil {
			r.logger.Error("failed to seed withdraw", zap.Int("withdraw", i+1), zap.Error(err))
			return fmt.Errorf("failed to seed withdraw %d: %w", i+1, err)
		}

		_, err = r.db.UpdateWithdrawStatus(r.ctx, db.UpdateWithdrawStatusParams{
			WithdrawID: withdraw.WithdrawID,
			Status:     status,
		})

		if err != nil {
			r.logger.Error("failed to update withdraw status", zap.Int("withdrawID", int(withdraw.WithdrawID)), zap.String("status", status), zap.Error(err))
			return fmt.Errorf("failed to update withdraw status for withdraw ID %d: %w", withdraw.WithdrawID, err)
		}

		if i < 20 {
			_, err = r.db.TrashWithdraw(r.ctx, withdraw.WithdrawID)
			if err != nil {
				r.logger.Error("failed to trash withdraw", zap.Int("withdraw", i+1), zap.Error(err))
				return fmt.Errorf("failed to trash withdraw %d: %w", i+1, err)
			}
		}
	}

	r.logger.Info("withdraw seeded successfully")

	return nil
}
