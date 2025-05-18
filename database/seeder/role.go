package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/exp/rand"
)

type roleSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewRoleSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *roleSeeder {
	return &roleSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *roleSeeder) Seed() error {
	prefixes := []string{"Super", "Admin", "User", "Manager", "Editor", "Viewer", "Guest", "Support", "Developer", "Analyst"}
	suffixes := []string{"Role", "Access", "Level", "Permission", "Group", "Team", "Control", "Admin", "User", "Manager"}

	totalRoles := 20
	activeRoles := 10
	trashedRoles := 10

	for i := 0; i < totalRoles; i++ {
		prefix := prefixes[rand.Intn(len(prefixes))]
		suffix := suffixes[rand.Intn(len(suffixes))]
		roleName := fmt.Sprintf("%s %s %d", prefix, suffix, i+1)

		role, err := r.db.CreateRole(r.ctx, roleName)
		if err != nil {
			r.logger.Error("failed to seed role", zap.Int("role", i+1), zap.String("roleName", roleName), zap.Error(err))
			return fmt.Errorf("failed to seed role %d (%s): %w", i+1, roleName, err)
		}

		if i >= activeRoles {
			_, err = r.db.TrashRole(r.ctx, role.RoleID)
			if err != nil {
				r.logger.Error("failed to trash role", zap.Int("role", i+1), zap.String("roleName", roleName), zap.Error(err))
				return fmt.Errorf("failed to trash role %d (%s): %w", i+1, roleName, err)
			}
		}
	}

	r.logger.Debug("role seeded successfully", zap.Int("totalRoles", totalRoles), zap.Int("activeRoles", activeRoles), zap.Int("trashedRoles", trashedRoles))

	return nil
}
