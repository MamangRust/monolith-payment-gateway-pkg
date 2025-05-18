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

type transactionSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewTransactionSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *transactionSeeder {
	return &transactionSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *transactionSeeder) Seed() error {
	paymentMethods := []string{"Bank Alpha", "Bank Beta", "Bank Gamma"}
	statusOptions := []string{"pending", "success", "failed"}

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
		r.logger.Error("no cards available for transaction seeding")
		return fmt.Errorf("no cards available for transaction seeding")
	}

	merchants, err := r.db.GetMerchants(r.ctx, db.GetMerchantsParams{
		Column1: "",
		Limit:   10,
		Offset:  0,
	})
	if err != nil {
		r.logger.Error("failed to get merchant list", zap.Error(err))
		return fmt.Errorf("failed to get merchant list: %w", err)
	}

	if len(merchants) == 0 {
		r.logger.Error("no merchants available for transaction seeding")
		return fmt.Errorf("no merchants available for transaction seeding")
	}

	months := make([]time.Time, 12)
	currentYear := time.Now().Year()
	for i := 0; i < 12; i++ {
		months[i] = time.Date(currentYear, time.Month(i+1), 1, 0, 0, 0, 0, time.UTC)
	}

	for i := 0; i < 40; i++ {
		selectedCard := cards[rand.Intn(len(cards))]
		selectedMerchant := merchants[rand.Intn(len(merchants))]
		selectedPaymentMethod := paymentMethods[rand.Intn(len(paymentMethods))]
		transactionAmount := int32(rand.Intn(1000000-50000) + 50000)

		status := statusOptions[rand.Intn(len(statusOptions))]

		monthIndex := i % 12
		transactionTime := months[monthIndex].Add(time.Duration(rand.Intn(28)) * 24 * time.Hour)

		request := db.CreateTransactionParams{
			CardNumber:      selectedCard.CardNumber,
			Amount:          transactionAmount,
			PaymentMethod:   selectedPaymentMethod,
			MerchantID:      selectedMerchant.MerchantID,
			TransactionTime: transactionTime,
		}

		transaction, err := r.db.CreateTransaction(r.ctx, request)
		if err != nil {
			r.logger.Error("failed to seed transaction", zap.Int("transaction", i+1), zap.Error(err))
			return fmt.Errorf("failed to seed transaction %d: %w", i+1, err)
		}

		_, err = r.db.UpdateTransactionStatus(r.ctx, db.UpdateTransactionStatusParams{
			TransactionID: transaction.TransactionID,
			Status:        status,
		})

		if err != nil {
			r.logger.Error("failed to update transaction status", zap.Int("transactionID", int(transaction.TransactionID)), zap.String("status", status), zap.Error(err))
			return fmt.Errorf("failed to update transaction status for transaction ID %d: %w", transaction.TransactionID, err)
		}

		if i < 20 {
			_, err = r.db.TrashTransaction(r.ctx, transaction.TransactionID)
			if err != nil {
				r.logger.Error("failed to trash transaction", zap.Int("transaction", i+1), zap.Error(err))
				return fmt.Errorf("failed to trash transaction %d: %w", i+1, err)
			}
		}
	}

	r.logger.Info("transaction seeded successfully")

	return nil
}
