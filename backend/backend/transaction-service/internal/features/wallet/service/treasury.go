package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
)

// topDonorsLimit — сколько донатеров показываем в treasury-статистике.
const topDonorsLimit = 10

// ────────────────────────────────────────────────────────────
//  Малые интерфейсы для treasury-методов (мокируются в тестах).
// ────────────────────────────────────────────────────────────

// TreasuryRepo — подмножество wallet repo, нужное только для treasury.
type TreasuryRepo interface {
	GetChannelWallet(ctx context.Context, channelID uuid.UUID) (*wallet.ChannelWallet, error)
	GetChannelDonationStats(ctx context.Context, channelID uuid.UUID) (int64, *big.Int, error)
	GetChannelTopDonors(ctx context.Context, channelID uuid.UUID, limit int) ([]repository.TopDonor, error)
	GetChannelDonations(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]repository.DonationRecord, error)
}

// TreasuryReader — on-chain reader для treasury.
type TreasuryReader interface {
	SudaBalanceOf(ctx context.Context, addr common.Address) (*big.Int, error)
}

// TreasuryMessenger — подмножество messenger client для treasury.
type TreasuryMessenger interface {
	CheckChannelPermission(ctx context.Context, userID, channelID, permission string) (*platgrpc.ChannelPermissionResult, error)
	GetUsersByIDs(ctx context.Context, ids []string) (map[string]platgrpc.UserBrief, error)
}

// ────────────────────────────────────────────────────────────
//  TreasuryService
// ────────────────────────────────────────────────────────────

// TreasuryService реализует treasury-статистику через маленькие интерфейсы.
// Встраивается в wallet.Service через поле ts.
type TreasuryService struct {
	repo      TreasuryRepo
	reader    TreasuryReader
	messenger TreasuryMessenger
}

// NewTreasuryService создаёт TreasuryService с конкретными зависимостями.
func NewTreasuryService(r TreasuryRepo, rd TreasuryReader, m TreasuryMessenger) *TreasuryService {
	return &TreasuryService{repo: r, reader: rd, messenger: m}
}

// NewTreasuryForTest создаёт TreasuryService для unit-тестов (можно передавать моки).
func NewTreasuryForTest(r TreasuryRepo, rd TreasuryReader, m TreasuryMessenger) *TreasuryService {
	return &TreasuryService{repo: r, reader: rd, messenger: m}
}

// ────────────────────────────────────────────────────────────
//  Доменные типы
// ────────────────────────────────────────────────────────────

// TopDonorInfo — донатер с резолвнутым username для UI.
type TopDonorInfo struct {
	UserID        string
	Username      string
	DisplayName   string
	TotalWei      string
	DonationCount int64
}

// TreasuryStats — агрегированная статистика treasury канала.
type TreasuryStats struct {
	ChannelID           string
	Address             string
	SudaBalanceWei      string
	TotalDonationsCount int64
	TotalDonationsWei   string
	TopDonors           []TopDonorInfo
}

// DonationInfo — одна строка списка донатов с резолвнутым username.
type DonationInfo struct {
	ID              string
	FromUserID      string
	FromUsername    string
	FromDisplayName string
	FromAddress     string
	AmountWei       string
	Message         string
	TxHash          string
	CreatedAt       time.Time
}

// DonationsList — пагинированный список донатов канала.
type DonationsList struct {
	Items []DonationInfo
	Total int64
}

// ────────────────────────────────────────────────────────────
//  Методы TreasuryService
// ────────────────────────────────────────────────────────────

// GetChannelTreasury возвращает статистику treasury канала.
// Доступно только OWNER/ADMIN канала (проверка через gRPC к messenger).
func (ts *TreasuryService) GetChannelTreasury(
	ctx context.Context, userID, channelID uuid.UUID,
) (*TreasuryStats, error) {
	if err := ts.requireChannelOwnerOrAdmin(ctx, userID, channelID); err != nil {
		return nil, err
	}

	// Кошелёк канала + on-chain баланс.
	w, err := ts.repo.GetChannelWallet(ctx, channelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, fmt.Errorf("treasury: lookup channel wallet: %w", err)
	}
	balance, err := ts.reader.SudaBalanceOf(ctx, common.HexToAddress(w.Address))
	if err != nil {
		return nil, fmt.Errorf("treasury: on-chain read: %w", err)
	}

	// Агрегаты донатов.
	count, totalWei, err := ts.repo.GetChannelDonationStats(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("treasury: donation stats: %w", err)
	}

	// Топ-донатеры + резолв username'ов.
	top, err := ts.repo.GetChannelTopDonors(ctx, channelID, topDonorsLimit)
	if err != nil {
		return nil, fmt.Errorf("treasury: top donors: %w", err)
	}
	donorIDs := make([]string, 0, len(top))
	for _, d := range top {
		donorIDs = append(donorIDs, d.FromUserID.String())
	}
	users := ts.resolveUsers(ctx, donorIDs)

	topDonors := make([]TopDonorInfo, 0, len(top))
	for _, d := range top {
		brief := users[d.FromUserID.String()]
		topDonors = append(topDonors, TopDonorInfo{
			UserID:        d.FromUserID.String(),
			Username:      brief.Username,
			DisplayName:   brief.DisplayName,
			TotalWei:      bigStr(d.TotalWei),
			DonationCount: d.DonationCount,
		})
	}

	return &TreasuryStats{
		ChannelID:           channelID.String(),
		Address:             w.Address,
		SudaBalanceWei:      balance.String(),
		TotalDonationsCount: count,
		TotalDonationsWei:   bigStr(totalWei),
		TopDonors:           topDonors,
	}, nil
}

// GetChannelDonations возвращает пагинированный список донатов канала.
// Доступно только OWNER/ADMIN канала.
func (ts *TreasuryService) GetChannelDonations(
	ctx context.Context, userID, channelID uuid.UUID, limit, offset int,
) (*DonationsList, error) {
	if err := ts.requireChannelOwnerOrAdmin(ctx, userID, channelID); err != nil {
		return nil, err
	}

	records, err := ts.repo.GetChannelDonations(ctx, channelID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("treasury: list donations: %w", err)
	}

	count, _, err := ts.repo.GetChannelDonationStats(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("treasury: donation count: %w", err)
	}

	// Резолв username'ов уникальных отправителей.
	idSet := make(map[string]struct{}, len(records))
	for _, r := range records {
		idSet[r.FromUserID.String()] = struct{}{}
	}
	ids := make([]string, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	users := ts.resolveUsers(ctx, ids)

	items := make([]DonationInfo, 0, len(records))
	for _, r := range records {
		brief := users[r.FromUserID.String()]
		items = append(items, DonationInfo{
			ID:              r.ID.String(),
			FromUserID:      r.FromUserID.String(),
			FromUsername:    brief.Username,
			FromDisplayName: brief.DisplayName,
			FromAddress:     r.FromAddress,
			AmountWei:       bigStr(r.Amount),
			Message:         r.Message,
			TxHash:          r.TxHash,
			CreatedAt:       r.CreatedAt,
		})
	}

	return &DonationsList{Items: items, Total: count}, nil
}

// requireChannelOwnerOrAdmin — проверка прав через messenger gRPC.
func (ts *TreasuryService) requireChannelOwnerOrAdmin(ctx context.Context, userID, channelID uuid.UUID) error {
	perm, err := ts.messenger.CheckChannelPermission(
		ctx, userID.String(), channelID.String(), platgrpc.PermissionOwnerOrAdmin,
	)
	if err != nil {
		log.Error().Err(err).
			Str("user_id", userID.String()).
			Str("channel_id", channelID.String()).
			Msg("treasury: permission check failed")
		return fmt.Errorf("treasury: permission check: %w", err)
	}
	if !perm.Granted {
		return ErrForbidden
	}
	return nil
}

// resolveUsers — best-effort batch-резолв username'ов. При ошибке gRPC
// возвращает пустую map: treasury-статистика отдаётся без username'ов.
func (ts *TreasuryService) resolveUsers(ctx context.Context, ids []string) map[string]platgrpc.UserBrief {
	if len(ids) == 0 {
		return map[string]platgrpc.UserBrief{}
	}
	users, err := ts.messenger.GetUsersByIDs(ctx, ids)
	if err != nil {
		log.Warn().Err(err).Int("ids", len(ids)).Msg("treasury: resolve usernames failed")
		return map[string]platgrpc.UserBrief{}
	}
	return users
}

// bigStr — *big.Int в строку, nil → "0".
func bigStr(n *big.Int) string {
	if n == nil {
		return "0"
	}
	return n.String()
}

// ────────────────────────────────────────────────────────────
//  Делегация с wallet.Service
// ────────────────────────────────────────────────────────────

// GetChannelTreasury — делегирует в TreasuryService.
func (s *Service) GetChannelTreasury(ctx context.Context, userID, channelID uuid.UUID) (*TreasuryStats, error) {
	return s.ts.GetChannelTreasury(ctx, userID, channelID)
}

// GetChannelDonations — делегирует в TreasuryService.
func (s *Service) GetChannelDonations(ctx context.Context, userID, channelID uuid.UUID, limit, offset int) (*DonationsList, error) {
	return s.ts.GetChannelDonations(ctx, userID, channelID, limit, offset)
}
