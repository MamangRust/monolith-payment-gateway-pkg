package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
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
	totalSaldos := 40
	activeSaldos := 20
	trashedSaldos := 20

	cards, err := r.db.GetCards(r.ctx, db.GetCardsParams{
		Column1: "",
		Limit:   40,
		Offset:  0,
	})

	if err != nil {
		r.logger.Error("failed to get card list", zap.Error(err))
		return fmt.Errorf("failed to get card list: %w", err)
	}

	if len(cards) == 0 {
		r.logger.Error("no cards available to seed saldo")
		return fmt.Errorf("no cards available to seed saldo")
	}

	for i := 0; i < totalSaldos; i++ {
		card := cards[i%len(cards)]

		totalBalance := int32(rand.Intn(10000000) + 1000000)

		request := db.CreateSaldoParams{
			CardNumber:   card.CardNumber,
			TotalBalance: totalBalance,
		}

		saldo, err := r.db.CreateSaldo(r.ctx, request)
		if err != nil {
			r.logger.Error("failed to seed saldo for card", zap.String("card", card.CardNumber), zap.Error(err))
			return fmt.Errorf("failed to seed saldo for card %s: %w", card.CardNumber, err)
		}

		if i >= activeSaldos {
			_, err = r.db.TrashSaldo(r.ctx, saldo.SaldoID)
			if err != nil {
				r.logger.Error("failed to trash saldo", zap.Int("saldo", i+1), zap.String("card", card.CardNumber), zap.Error(err))
				return fmt.Errorf("failed to trash saldo %d for card %s: %w", i+1, card.CardNumber, err)
			}
		}
	}

	r.logger.Debug("saldo seeded successfully", zap.Int("totalSaldos", totalSaldos), zap.Int("activeSaldos", activeSaldos), zap.Int("trashedSaldos", trashedSaldos))

	return nil
}
