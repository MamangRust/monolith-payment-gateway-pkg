package seeder

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userSeeder struct {
	db     *db.Queries
	hash   hash.HashPassword
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewUserSeeder(db *db.Queries, ctx context.Context, hash hash.HashPassword, logger logger.LoggerInterface) *userSeeder {
	return &userSeeder{
		db:     db,
		hash:   hash,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *userSeeder) Seed() error {
	totalUsers := 15
	activeUsers := 10
	trashedUsers := totalUsers - activeUsers

	randomRoles := []string{
		"Super Admin",
		"Admin",
		"Merchant Admin",
		"Merchant Operator",
		"Finance",
		"Compliance",
		"Auditor",
		"Support",
		"Viewer",
		"User",
	}

	for i := 1; i <= totalUsers; i++ {
		email := fmt.Sprintf("user_%s@example.com", uuid.NewString())

		hash, err := r.hash.HashPassword("password")
		if err != nil {
			r.logger.Error("failed to hash password", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to hash password for user %d: %w", i, err)
		}

		user := db.CreateUserParams{
			Firstname:        fmt.Sprintf("User%d", i),
			Lastname:         fmt.Sprintf("Last%d", i),
			Email:            email,
			Password:         hash,
			VerificationCode: uuid.NewString(),
			IsVerified:       sql.NullBool{Bool: true, Valid: true},
		}

		createdUser, err := r.db.CreateUser(r.ctx, user)
		if err != nil {
			r.logger.Error("failed to create user", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to create user %d: %w", i, err)
		}

		randomRole := randomRoles[rand.Intn(len(randomRoles))]
		role, err := r.db.GetRoleByName(r.ctx, randomRole)
		if err != nil {
			r.logger.Error("failed to get role", zap.String("role", randomRole), zap.Error(err))
			return fmt.Errorf("failed to get role %s: %w", randomRole, err)
		}

		_, err = r.db.AssignRoleToUser(r.ctx, db.AssignRoleToUserParams{
			RoleID: role.RoleID,
			UserID: createdUser.UserID,
		})
		if err != nil {
			r.logger.Error("failed to assign role to user", zap.Int("userID", int(createdUser.UserID)), zap.String("role", randomRole), zap.Error(err))
			return fmt.Errorf("failed to assign role %s to user %d: %w", randomRole, createdUser.UserID, err)
		}

		if i > activeUsers {
			_, err := r.db.TrashUser(r.ctx, createdUser.UserID)
			if err != nil {
				r.logger.Error("failed to trash user", zap.Int("user", i), zap.Error(err))
				return fmt.Errorf("failed to trash user %d: %w", i, err)
			}
		}
	}

	r.logger.Info("user seeded successfully", zap.Int("totalUsers", totalUsers), zap.Int("activeUsers", activeUsers), zap.Int("trashedUsers", trashedUsers))

	return nil
}
