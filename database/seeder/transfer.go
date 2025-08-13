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

// transferSeeder is a struct that represents a seeder for the transfers table.
type transferSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

// NewTransferSeeder creates a new instance of the transferSeeder, which is
// responsible for populating the transfers table with fake data.
//
// Args:
// db: a pointer to the database queries
// ctx: a context.Context object
// logger: a logger.LoggerInterface object
//
// Returns:
// a pointer to the transferSeeder struct
func NewTransferSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *transferSeeder {
	return &transferSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

// Seed populates the transfers table with fake data.
//
// It creates a total of 10 transfers, with 5 of them being active and 5 of them being
// trashed. The active transfers have a status of "pending", "success", or "failed", while
// the trashed transfers have a status of "trashed". The transfers are seeded with random
// transfer from and to, transfer time, and transfer amount. The transfer time is randomly
// selected from the first day of each month of the current year.
//
// If any of the transfers fail to be created, the function returns an error.
//
// Returns:
// an error if any of the transfers fail to be created, otherwise nil
func (r *transferSeeder) Seed() error {
	total := 5
	active := 3
	trashed := 2

	var cards []db.Card
	for i := 1; i <= total; i++ {
		card, err := r.db.GetCardByUserID(r.ctx, int32(i))
		if err != nil {
			r.logger.Debug("failed to get card for user", zap.Int("userID", i), zap.Error(err))
			continue
		}
		if card != nil {
			cards = append(cards, *card)
		}
	}

	if len(cards) < 2 {
		r.logger.Error("not enough cards available for transfer seeding", zap.Int("available", len(cards)))
		return fmt.Errorf("need at least 2 cards, got %d", len(cards))
	}

	statusOptions := []string{"pending", "success", "failed"}

	months := make([]time.Time, 12)
	currentYear := time.Now().Year()
	for i := 0; i < 12; i++ {
		months[i] = time.Date(currentYear, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
	}

	for i := 0; i < total; i++ {
		fromIndex := rand.Intn(len(cards))
		toIndex := rand.Intn(len(cards))
		for fromIndex == toIndex {
			toIndex = rand.Intn(len(cards))
		}

		transferFrom := cards[fromIndex].CardNumber
		transferTo := cards[toIndex].CardNumber
		amount := int32(rand.Intn(1000000) + 50000)
		status := statusOptions[rand.Intn(len(statusOptions))]

		monthIndex := i % 12
		transferTime := months[monthIndex].Add(time.Duration(rand.Intn(28)) * 24 * time.Hour)

		req := db.CreateTransferParams{
			TransferFrom:   transferFrom,
			TransferTo:     transferTo,
			TransferAmount: amount,
			TransferTime:   transferTime,
			Status:         status,
		}

		transfer, err := r.db.CreateTransfer(r.ctx, req)
		if err != nil {
			r.logger.Error("failed to seed transfer", zap.Int("transfer", i+1), zap.Error(err))
			return fmt.Errorf("failed to seed transfer %d: %w", i+1, err)
		}

		if i >= active {
			_, err = r.db.TrashTransfer(r.ctx, transfer.TransferID)
			if err != nil {
				r.logger.Error("failed to trash transfer", zap.Int("transfer", i+1), zap.Error(err))
				return fmt.Errorf("failed to trash transfer %d: %w", i+1, err)
			}
		}
	}

	r.logger.Info("transfer seeded successfully",
		zap.Int("total", total),
		zap.Int("active", active),
		zap.Int("trashed", trashed))

	return nil
}
