package seeder

import (
	"context"
	"fmt"

	apikey "github.com/MamangRust/monolith-payment-gateway-pkg/api-key"
	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
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
	adjectives := []string{"Blue", "Green", "Red", "Yellow", "Fast", "Smart", "Global", "Local", "Happy", "Bright"}
	nouns := []string{"Shop", "Store", "Mart", "Market", "Hub", "Center", "Place", "Corner", "Zone", "Point"}

	statusOptions := []string{"pending", "active", "deactive"}

	totalMerchants := 40
	activeMerchants := 20
	trashedMerchants := 20

	for i := 0; i < totalMerchants; i++ {
		adjective := adjectives[rand.Intn(len(adjectives))]
		noun := nouns[rand.Intn(len(nouns))]
		merchantName := fmt.Sprintf("%s %s", adjective, noun)

		apiKey := apikey.GenerateApiKey()

		req := db.CreateMerchantParams{
			Name:   merchantName,
			UserID: int32(rand.Intn(40) + 1),
			ApiKey: apiKey,
		}

		merchant, err := r.db.CreateMerchant(r.ctx, req)
		if err != nil {
			r.logger.Error("failed to seed merchant", zap.Int("merchant", i+1), zap.Error(err))
			return fmt.Errorf("failed to seed merchant %d: %w", i+1, err)
		}

		status := statusOptions[rand.Intn(len(statusOptions))]

		_, err = r.db.UpdateMerchantStatus(r.ctx, db.UpdateMerchantStatusParams{
			MerchantID: merchant.MerchantID,
			Status:     status,
		})

		if err != nil {
			r.logger.Error("failed to update merchant status", zap.Int("merchantID", int(merchant.MerchantID)), zap.String("status", status), zap.Error(err))
			return fmt.Errorf("failed to update merchant status for merchant ID %d: %w", merchant.MerchantID, err)
		}

		if i >= activeMerchants {
			_, err = r.db.TrashMerchant(r.ctx, merchant.MerchantID)
			if err != nil {
				r.logger.Error("failed to trash merchant", zap.Int("merchant", i+1), zap.Error(err))
				return fmt.Errorf("failed to trash merchant %d: %w", i+1, err)
			}
		}
	}

	r.logger.Debug("merchant seeded successfully", zap.Int("totalMerchants", totalMerchants), zap.Int("activeMerchants", activeMerchants), zap.Int("trashedMerchants", trashedMerchants))

	return nil
}
