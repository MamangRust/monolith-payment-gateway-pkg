package seeder

import (
	"context"
	"fmt"

	apikey "github.com/MamangRust/monolith-payment-gateway-pkg/api-key"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.uber.org/zap"
)

type merchantSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewMerchantSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *merchantSeeder {
	return &merchantSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *merchantSeeder) Seed() error {
	adjectives := []string{"Blue", "Green", "Red", "Yellow", "Fast"}
	nouns := []string{"Shop", "Store", "Mart", "Market", "Hub"}

	totalMerchants := 10
	activeMerchants := 5
	trashedMerchants := 5

	for i := 0; i < totalMerchants; i++ {
		adjective := adjectives[i%len(adjectives)]
		noun := nouns[i%len(nouns)]
		merchantName := fmt.Sprintf("%s %s", adjective, noun)

		apiKey := apikey.GenerateApiKey()

		req := db.CreateMerchantParams{
			Name:   merchantName,
			UserID: int32((i % 5) + 1),
			ApiKey: apiKey,
		}

		merchant, err := r.db.CreateMerchant(r.ctx, req)
		if err != nil {
			r.logger.Error("failed to seed merchant", zap.Int("merchant", i+1), zap.Error(err))
			return fmt.Errorf("failed to seed merchant %d: %w", i+1, err)
		}

		var status string
		if i < activeMerchants {
			status = "active"
		} else {
			status = "deactive"
		}

		_, err = r.db.UpdateMerchantStatus(r.ctx, db.UpdateMerchantStatusParams{
			MerchantID: merchant.MerchantID,
			Status:     status,
		})
		if err != nil {
			r.logger.Error("failed to update merchant status", zap.Int("merchantID", int(merchant.MerchantID)), zap.String("status", status), zap.Error(err))
			return fmt.Errorf("failed to update status for merchant ID %d: %w", merchant.MerchantID, err)
		}

		if i >= activeMerchants {
			_, err = r.db.TrashMerchant(r.ctx, merchant.MerchantID)
			if err != nil {
				r.logger.Error("failed to trash merchant", zap.Int("merchant", i+1), zap.Error(err))
				return fmt.Errorf("failed to trash merchant %d: %w", i+1, err)
			}
		}
	}

	r.logger.Info("merchant seeded successfully",
		zap.Int("totalMerchants", totalMerchants),
		zap.Int("activeMerchants", activeMerchants),
		zap.Int("trashedMerchants", trashedMerchants))

	return nil
}
