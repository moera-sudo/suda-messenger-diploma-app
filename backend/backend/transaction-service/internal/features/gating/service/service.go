// Package service — бизнес-логика управления token-gating правилами.
//
// Сама ПРОВЕРКА gating (проходит ли юзер) живёт в wallet/service/gating.go
// и вызывается через gRPC. Здесь — только CRUD правил для owner'а канала.
package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/gating"
	walletrepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ────────────────────────────────────────────────────────────
//  Маленькие интерфейсы для зависимостей (мокаются в тестах).
// ────────────────────────────────────────────────────────────

type gatingRepo interface {
	GetGatingRule(ctx context.Context, chatID uuid.UUID) (*walletrepo.GatingRule, error)
	UpsertGatingRule(ctx context.Context, chatID uuid.UUID, minSudaBalance *big.Int, nftCollectionID *uuid.UUID, subscriptionPrice *big.Int, createdBy uuid.UUID) error
	DeleteGatingRule(ctx context.Context, chatID uuid.UUID) error
}

type messengerClient interface {
	CheckChannelPermission(ctx context.Context, userID, channelID, permission string) (*platgrpc.ChannelPermissionResult, error)
}

// Service — CRUD token-gating правил.
type Service struct {
	repo      gatingRepo
	messenger messengerClient
}

// Deps — зависимости конструктора.
type Deps struct {
	Postgres        *pgxpool.Pool
	MessengerClient *platgrpc.MessengerClient
}

func New(d Deps) *Service {
	return &Service{
		repo:      walletrepo.New(d.Postgres),
		messenger: d.MessengerClient,
	}
}

// NewWithDeps — для тестов: инжект моков напрямую.
func NewWithDeps(repo gatingRepo, mc messengerClient) *Service {
	return &Service{repo: repo, messenger: mc}
}

// Rule — token-gating правило для отдачи наружу (HTTP / observer).
type Rule struct {
	ChatID                  uuid.UUID
	MinSudaBalanceWei       string
	RequiredNFTCollectionID *uuid.UUID
	SubscriptionPriceWei    string // цена платной подписки на PUBLIC-канал ("0" = не платный)
	HasRule                 bool   // false → правила нет, канал open для всех
}

// RuleInput — параметры POST /gating/rule.
type RuleInput struct {
	OwnerID                 uuid.UUID
	ChatID                  uuid.UUID
	MinSudaBalanceWei       string // legacy hold-порог; "" трактуется как "0"
	RequiredNFTCollectionID *uuid.UUID
	SubscriptionPriceWei    string // цена платной подписки; "" → "0"
}

// CreateRule создаёт или обновляет правило gating. Доступно только OWNER'у канала.
func (s *Service) CreateRule(ctx context.Context, in RuleInput) error {
	minBalance, err := parseNonNegWei(in.MinSudaBalanceWei)
	if err != nil {
		return fmt.Errorf("%w: min_suda_balance_wei must be a non-negative integer string", gating.ErrInvalidInput)
	}
	price, err := parseNonNegWei(in.SubscriptionPriceWei)
	if err != nil {
		return fmt.Errorf("%w: subscription_price_wei must be a non-negative integer string", gating.ErrInvalidInput)
	}

	if err := s.requireOwner(ctx, in.OwnerID, in.ChatID); err != nil {
		return err
	}

	if err := s.repo.UpsertGatingRule(ctx, in.ChatID, minBalance, in.RequiredNFTCollectionID, price, in.OwnerID); err != nil {
		return fmt.Errorf("create rule: %w", err)
	}

	log.Info().
		Str("chat_id", in.ChatID.String()).
		Str("owner_id", in.OwnerID.String()).
		Str("min_suda_balance_wei", minBalance.String()).
		Str("subscription_price_wei", price.String()).
		Msg("gating rule upserted")
	return nil
}

// parseNonNegWei парсит десятичную wei-строку; "" трактуется как 0.
// Возвращает ошибку для не-числа или отрицательного значения.
func parseNonNegWei(s string) (*big.Int, error) {
	if s == "" {
		return big.NewInt(0), nil
	}
	n, ok := new(big.Int).SetString(s, 10)
	if !ok || n.Sign() < 0 {
		return nil, fmt.Errorf("invalid wei")
	}
	return n, nil
}

// GetRule возвращает текущее правило канала. Permission НЕ требуется —
// порог доступа может видеть любой (UI показывает "вход стоит N SUDA").
func (s *Service) GetRule(ctx context.Context, chatID uuid.UUID) (*Rule, error) {
	rule, err := s.repo.GetGatingRule(ctx, chatID)
	if err != nil {
		if errors.Is(err, walletrepo.ErrNotFound) {
			return &Rule{ChatID: chatID, MinSudaBalanceWei: "0", SubscriptionPriceWei: "0", HasRule: false}, nil
		}
		return nil, fmt.Errorf("get rule: %w", err)
	}
	return &Rule{
		ChatID:                  rule.ChatID,
		MinSudaBalanceWei:       rule.MinSudaBalance.String(),
		RequiredNFTCollectionID: rule.RequiredNFTCollectionID,
		SubscriptionPriceWei:    bigStrOrZero(rule.SubscriptionPrice),
		HasRule:                 true,
	}, nil
}

// bigStrOrZero — *big.Int в строку, nil → "0".
func bigStrOrZero(n *big.Int) string {
	if n == nil {
		return "0"
	}
	return n.String()
}

// DeleteRule снимает правило gating. Доступно только OWNER'у канала.
func (s *Service) DeleteRule(ctx context.Context, ownerID, chatID uuid.UUID) error {
	if err := s.requireOwner(ctx, ownerID, chatID); err != nil {
		return err
	}
	if err := s.repo.DeleteGatingRule(ctx, chatID); err != nil {
		return fmt.Errorf("delete rule: %w", err)
	}
	log.Info().Str("chat_id", chatID.String()).Str("owner_id", ownerID.String()).Msg("gating rule deleted")
	return nil
}

// requireOwner проверяет через messenger, что userID — OWNER канала chatID.
func (s *Service) requireOwner(ctx context.Context, userID, chatID uuid.UUID) error {
	perm, err := s.messenger.CheckChannelPermission(ctx, userID.String(), chatID.String(), platgrpc.PermissionOwner)
	if err != nil {
		return fmt.Errorf("gating: permission check: %w", err)
	}
	if !perm.Granted {
		return gating.ErrForbidden
	}
	return nil
}
