package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/repository"
)

// InitiateInput — параметры POST /purchase/initiate.
type InitiateInput struct {
	UserID        uuid.UUID
	PackageCode   string
	PaymentMethod string // CARD | APPLE_PAY | GOOGLE_PAY (декоративно)
}

// InitiateResult — то, что вернётся клиенту.
type InitiateResult struct {
	PurchaseID    uuid.UUID
	PackageCode   string
	SudaAmountWei string // десятичная строка
	FiatAmount    string // "9.99"
	FiatCurrency  string
	Status        string // всегда PENDING сразу после initiate
	PaymentMethod string
}

// Initiate создаёт row в статусе PENDING и возвращает purchase_id.
// Юзеру UI потом покажет форму "карты", после чего он позовёт /confirm.
func (s *Service) Initiate(ctx context.Context, in InitiateInput) (*InitiateResult, error) {
	pkg, ok := purchase.FindPackage(in.PackageCode)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrPackageUnknown, in.PackageCode)
	}

	method := in.PaymentMethod
	if method == "" {
		method = purchase.PaymentMethodCard
	}
	if !purchase.ValidPaymentMethod(method) {
		return nil, fmt.Errorf("%w: payment_method", ErrInvalidInput)
	}

	id, err := s.repo.Create(ctx, repository.CreateParams{
		UserID:        in.UserID,
		PackageCode:   pkg.Code,
		SudaAmountWei: pkg.SudaAmountWei(),
		FiatAmount:    pkg.FiatPrice,
		FiatCurrency:  pkg.FiatCurrency,
		PaymentMethod: method,
	})
	if err != nil {
		return nil, fmt.Errorf("initiate: %w", err)
	}

	log.Info().
		Str("user_id", in.UserID.String()).
		Str("purchase_id", id.String()).
		Str("package", pkg.Code).
		Str("fiat_amount", pkg.FiatPrice).
		Str("fiat_currency", pkg.FiatCurrency).
		Msg("purchase initiated")

	return &InitiateResult{
		PurchaseID:    id,
		PackageCode:   pkg.Code,
		SudaAmountWei: pkg.SudaAmountWei().String(),
		FiatAmount:    pkg.FiatPrice,
		FiatCurrency:  pkg.FiatCurrency,
		Status:        purchase.StatusPending,
		PaymentMethod: method,
	}, nil
}
