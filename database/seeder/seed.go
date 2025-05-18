package seeder

import (
	"context"
	"fmt"
	"time"

	db "github.com/MamangRust/monolith-payment-gateway-pkg/database/schema"
	"github.com/MamangRust/monolith-payment-gateway-pkg/hash"
	"github.com/MamangRust/monolith-payment-gateway-pkg/logger"
)

type Deps struct {
	DB     *db.Queries
	Hash   hash.HashPassword
	Ctx    context.Context
	Logger logger.LoggerInterface
}

type Seeder struct {
	User        *userSeeder
	Role        *roleSeeder
	Saldo       *saldoSeeder
	Topup       *topupSeeder
	Withdraw    *withdrawSeeder
	Transfer    *transferSeeder
	Merchant    *merchantSeeder
	Card        *cardSeeder
	Transaction *transactionSeeder
}

func NewSeeder(deps Deps) *Seeder {
	return &Seeder{
		User:        NewUserSeeder(deps.DB, deps.Ctx, deps.Hash, deps.Logger),
		Role:        NewRoleSeeder(deps.DB, deps.Ctx, deps.Logger),
		Saldo:       NewSaldoSeeder(deps.DB, deps.Ctx, deps.Logger),
		Topup:       NewTopupSeeder(deps.DB, deps.Ctx, deps.Logger),
		Withdraw:    NewWithdrawSeeder(deps.DB, deps.Ctx, deps.Logger),
		Transfer:    NewTransferSeeder(deps.DB, deps.Ctx, deps.Logger),
		Merchant:    NewMerchantSeeder(deps.DB, deps.Ctx, deps.Logger),
		Card:        NewCardSeeder(deps.DB, deps.Ctx, deps.Logger),
		Transaction: NewTransactionSeeder(deps.DB, deps.Ctx, deps.Logger),
	}
}

func (s *Seeder) Run() error {
	if err := s.seedWithDelay("users", s.User.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("roles", s.Role.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("cards", s.Card.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("saldo", s.Saldo.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("topups", s.Topup.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("withdrawals", s.Withdraw.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("transfers", s.Transfer.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("merchants", s.Merchant.Seed); err != nil {
		return err
	}

	if err := s.seedWithDelay("transactions", s.Transaction.Seed); err != nil {
		return err
	}

	return nil
}

func (s *Seeder) seedWithDelay(entityName string, seedFunc func() error) error {
	if err := seedFunc(); err != nil {
		return fmt.Errorf("failed to seed %s: %w", entityName, err)
	}

	time.Sleep(30 * time.Second)
	return nil
}
