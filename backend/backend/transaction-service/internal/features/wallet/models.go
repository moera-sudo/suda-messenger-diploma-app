package wallet

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

// Wallet — кастодиальный кошелёк юзера.
type Wallet struct {
	UserID              uuid.UUID
	Address             string  // 0x... EVM address (42 chars)
	EncryptedPrivateKey string  // base64(nonce || ciphertext+tag), AES-256-GCM
	KeyVersion          uint8
	SudaBalanceCache    *big.Int  // wei, кеш — не источник правды
	BalanceUpdatedAt    *time.Time
	CreatedAt           time.Time
}

// ChannelWallet — treasury канала (отдельный кошелёк).
type ChannelWallet struct {
	ChannelID           uuid.UUID
	Address             string
	EncryptedPrivateKey string
	KeyVersion          uint8
	SudaBalanceCache    *big.Int
	BalanceUpdatedAt    *time.Time
	CreatedAt           time.Time
}

// Operation — типы операций, которые пишутся в tx_signing_audit.
// Должны совпадать со строками в SQL-комментариях миграции.
const (
	OpTransfer         = "TRANSFER"
	OpDonation         = "DONATION"
	OpNFTGift          = "NFT_GIFT"
	OpNFTList          = "NFT_LIST"
	OpNFTBuy           = "NFT_BUY"
	OpNFTCancelList    = "NFT_CANCEL_LIST"
	OpQuestCreate      = "QUEST_CREATE"
	OpQuestClaim       = "QUEST_CLAIM"
	OpQuestSubmit      = "QUEST_SUBMIT"
	OpQuestApprove     = "QUEST_APPROVE"
	OpQuestCancel      = "QUEST_CANCEL"
	OpFundraiseDonate  = "FUNDRAISE_DONATE"
	OpFundraiseCreate  = "FUNDRAISE_CREATE"
	OpChannelWithdraw  = "CHANNEL_WITHDRAW"
)

// Subject — тип владельца ключа, для tx_signing_audit.subject_type.
const (
	SubjectUser    = "USER"
	SubjectChannel = "CHANNEL"
)