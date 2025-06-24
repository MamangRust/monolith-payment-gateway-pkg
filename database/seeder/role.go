package seeder

import (
	"context"
	"fmt"

	"math/rand/v2"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
	"go.uber.org/zap"
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

	for roleIndex := 1; roleIndex <= totalRoles; roleIndex++ {
		prefix := prefixes[rand.IntN(len(prefixes))]
		suffix := suffixes[rand.IntN(len(suffixes))]
		roleName := fmt.Sprintf("%s_%s_%d", prefix, suffix, roleIndex)

		role, err := r.db.CreateRole(r.ctx, roleName)
		if err != nil {
			r.logger.Error("failed to seed role",
				zap.Int("role_index", roleIndex),
				zap.String("role_name", roleName),
				zap.Error(err),
			)
			return fmt.Errorf("failed to seed role %d (%s): %w", roleIndex, roleName, err)
		}

		if roleIndex > activeRoles {
			_, err = r.db.TrashRole(r.ctx, role.RoleID)
			if err != nil {
				r.logger.Error("failed to trash role",
					zap.Int("role_index", roleIndex),
					zap.String("role_name", roleName),
					zap.Error(err),
				)
				return fmt.Errorf("failed to trash role %d (%s): %w", roleIndex, roleName, err)
			}
		}
	}

	r.logger.Debug("role seeding completed",
		zap.Int("total_roles", totalRoles),
		zap.Int("active_roles", activeRoles),
		zap.Int("trashed_roles", trashedRoles),
	)

	return nil
}
