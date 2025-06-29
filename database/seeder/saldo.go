package seeder

import (
	"context"
	"fmt"

	"math/rand"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.uber.org/zap"
)

type saldoSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewSaldoSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *saldoSeeder {
	return &saldoSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *saldoSeeder) Seed() error {
	totalSaldos := 10
	activeSaldos := 5
	trashedSaldos := 5

	var cards []db.Card
	for i := 1; i <= totalSaldos; i++ {
		card, err := r.db.GetCardByUserID(r.ctx, int32(i))
		if err != nil {
			r.logger.Error("failed to get card for user", zap.Int("userID", i), zap.Error(err))
			return fmt.Errorf("failed to get card for user %d: %w", i, err)
		}
		if card == nil {
			r.logger.Error("no card found for user", zap.Int("userID", i))
			continue
		}
		cards = append(cards, *card)
	}

	if len(cards) < totalSaldos {
		r.logger.Error("not enough cards to seed saldo", zap.Int("required", totalSaldos), zap.Int("available", len(cards)))
		return fmt.Errorf("not enough cards to seed saldo: required %d, got %d", totalSaldos, len(cards))
	}

	for i, card := range cards {
		request := db.CreateSaldoParams{
			CardNumber:   card.CardNumber,
			TotalBalance: int32(rand.Intn(9_000_000) + 1_000_000),
		}

		saldo, err := r.db.CreateSaldo(r.ctx, request)
		if err != nil {
			r.logger.Error("failed to seed saldo", zap.Int("index", i), zap.String("card", card.CardNumber), zap.Error(err))
			return fmt.Errorf("failed to seed saldo for card %s: %w", card.CardNumber, err)
		}

		if i >= activeSaldos {
			_, err = r.db.TrashSaldo(r.ctx, saldo.SaldoID)
			if err != nil {
				r.logger.Error("failed to trash saldo", zap.Int("index", i), zap.String("card", card.CardNumber), zap.Error(err))
				return fmt.Errorf("failed to trash saldo %d for card %s: %w", i+1, card.CardNumber, err)
			}
		}
	}

	r.logger.Info("saldo seeded successfully",
		zap.Int("totalSaldos", totalSaldos),
		zap.Int("activeSaldos", activeSaldos),
		zap.Int("trashedSaldos", trashedSaldos))

	return nil
}
