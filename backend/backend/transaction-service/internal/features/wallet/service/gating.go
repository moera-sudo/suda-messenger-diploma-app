package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
)

// CheckTokenGating проверяет, удовлетворяет ли юзер требованиям gating
// для конкретного чата.
//
// Правило (tx_gating_rules):
//   - min_suda_balance — минимальный on-chain баланс SUDA
//   - required_nft_collection_id — опционально, нужен NFT из этой коллекции
//
// Возвращает: (passed, required, min_balance_wei, user_balance_wei, reason, err)
//   - required: было ли вообще правило (false → чат open для всех)
//   - reason:
//       "ok"                       — passed
//       "no_rule"                  — нет правила, всем можно
//       "wallet_not_found"         — у юзера нет кошелька
//       "insufficient_balance"     — баланс ниже min_suda_balance
//       "nft_required"             — нет NFT из требуемой коллекции
func (s *Service) CheckTokenGating(
	ctx context.Context, chatIDStr, userIDStr string,
) (passed, required bool, minBalanceWei, userBalanceWei, priceWei, reason string, err error) {
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		return false, false, "0", "0", "0", "", fmt.Errorf("%w: chat_id must be uuid", ErrInvalidInput)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return false, false, "0", "0", "0", "", fmt.Errorf("%w: user_id must be uuid", ErrInvalidInput)
	}

	// 1. Правило для чата.
	rule, err := s.repo.GetGatingRule(ctx, chatID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			// Правила нет → чат не gated, всем можно.
			return true, false, "0", "0", "0", "no_rule", nil
		}
		return false, false, "0", "0", "0", "", fmt.Errorf("gating: read rule: %w", err)
	}

	required = true
	minBalanceWei = rule.MinSudaBalance.String()
	priceWei = "0"
	if rule.SubscriptionPrice != nil {
		priceWei = rule.SubscriptionPrice.String()
	}

	// 2. Кошелёк юзера.
	w, err := s.repo.GetUserWallet(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return false, true, minBalanceWei, "0", priceWei, "wallet_not_found", nil
		}
		return false, true, minBalanceWei, "0", priceWei, "", fmt.Errorf("gating: lookup wallet: %w", err)
	}

	// 3. SUDA balance check (legacy hold-порог; для платной модели min обычно = 0).
	balance, err := s.reader.SudaBalanceOf(ctx, common.HexToAddress(w.Address))
	if err != nil {
		return false, true, minBalanceWei, "0", priceWei, "", fmt.Errorf("gating: read balance: %w", err)
	}
	userBalanceWei = balance.String()

	if balance.Cmp(rule.MinSudaBalance) < 0 {
		return false, true, minBalanceWei, userBalanceWei, priceWei, "insufficient_balance", nil
	}

	// 4. NFT ownership check (опционально).
	if rule.RequiredNFTCollectionID != nil {
		owns, err := s.repo.UserOwnsNFTFromCollection(ctx, userID, *rule.RequiredNFTCollectionID)
		if err != nil {
			return false, true, minBalanceWei, userBalanceWei, priceWei, "", fmt.Errorf("gating: nft check: %w", err)
		}
		if !owns {
			return false, true, minBalanceWei, userBalanceWei, priceWei, "nft_required", nil
		}
	}

	return true, true, minBalanceWei, userBalanceWei, priceWei, "ok", nil
}