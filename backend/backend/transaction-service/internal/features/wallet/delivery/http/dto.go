package http

// GET /wallet/me

// MeResponse — мой адрес + on-chain баланс.
type MeResponse struct {
	Address        string `json:"address"`
	SudaBalanceWei string `json:"suda_balance_wei"`
	Decimals       int32  `json:"decimals"`
}

// POST /wallet/resolve
//
// Используется UI: «я хочу перевести @john, покажи мне его адрес и
// display_name перед подтверждением».

type ResolveRequest struct {
	Username string `json:"username" validate:"required,min=1,max=50"`
}

type ResolveResponse struct {
	Found         bool   `json:"found"`
	UserID        string `json:"user_id,omitempty"`
	DisplayName   string `json:"display_name,omitempty"`
	WalletAddress string `json:"wallet_address,omitempty"`
}

// POST /wallet/transfer

// TransferRequest — перевод SUDA по @username.
//
// amount_wei — строка десятичного числа в wei (uint256 не лезет в JSON number).
// 100 SUDA = "100000000000000000000".
type TransferRequest struct {
	ToUsername string `json:"to_username" validate:"required,min=1,max=50"`
	AmountWei  string `json:"amount_wei"  validate:"required"`
	Note       string `json:"note,omitempty" validate:"max=200"`
	ChatID     string `json:"chat_id,omitempty" validate:"omitempty,uuid"`
}

// TransferResponse — async ответ: tx_hash сразу, статус через WS-event.
type TransferResponse struct {
	TxHash      string `json:"tx_hash"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	ToUserID    string `json:"to_user_id"`
	AmountWei   string `json:"amount_wei"`
	SubmittedAt string `json:"submitted_at"`
	ChatID      string `json:"chat_id,omitempty"`
}

// WithdrawRequest — вывод средств из казны канала на кошелёк владельца.
type WithdrawRequest struct {
	AmountWei string `json:"amount_wei" validate:"required"`
}

// WithdrawResponse — async ответ на вывод из казны.
type WithdrawResponse struct {
	TxHash      string `json:"tx_hash"`
	FromAddress string `json:"from_address"` // адрес казны
	ToAddress   string `json:"to_address"`   // кошелёк владельца
	AmountWei   string `json:"amount_wei"`
	SubmittedAt string `json:"submitted_at"`
}

// GET /wallet/history?limit=&offset=

type HistoryItem struct {
	ID                string `json:"id"`
	TxHash            string `json:"tx_hash"`
	LogIndex          int32  `json:"log_index"`
	FromUserID        string `json:"from_user_id,omitempty"`
	ToUserID          string `json:"to_user_id,omitempty"`
	FromAddress       string `json:"from_address"`
	ToAddress         string `json:"to_address"`
	AmountWei         string `json:"amount_wei"`
	Type              string `json:"type"`
	Status            string `json:"status"`
	RelatedEntityType string `json:"related_entity_type,omitempty"`
	RelatedEntityID   string `json:"related_entity_id,omitempty"`
	Note              string `json:"note,omitempty"`
	BlockNumber       int64  `json:"block_number"`
	ConfirmedAt       string `json:"confirmed_at"`
}

type HistoryResponse struct {
	Items []HistoryItem `json:"items"`
}

// GET /wallet/channel/:channel_id
//
// Только для OWNER/ADMIN канала (проверка через gRPC к messenger).

type ChannelWalletResponse struct {
	ChannelID      string `json:"channel_id"`
	Address        string `json:"address"`
	SudaBalanceWei string `json:"suda_balance_wei"`
	Decimals       int32  `json:"decimals"`
}

// GET /wallet/channel/:channel_id/treasury
//
// Статистика treasury канала. Только для OWNER/ADMIN.

type TopDonorDTO struct {
	UserID        string `json:"user_id"`
	Username      string `json:"username,omitempty"`
	DisplayName   string `json:"display_name,omitempty"`
	AmountWei     string `json:"amount_wei"`
	DonationCount int64  `json:"donation_count"`
}

type TreasuryResponse struct {
	ChannelID           string        `json:"channel_id"`
	Address             string        `json:"address"`
	SudaBalanceWei      string        `json:"suda_balance_wei"`
	Decimals            int32         `json:"decimals"`
	TotalDonationsCount int64         `json:"total_donations_count"`
	TotalDonationsWei   string        `json:"total_donations_wei"`
	TopDonors           []TopDonorDTO `json:"top_donors"`
}

// GET /wallet/channel/:channel_id/donations?limit=&offset=
//
// Пагинированный список донатов канала. Только для OWNER/ADMIN.

type DonationDTO struct {
	ID              string `json:"id"`
	FromUserID      string `json:"from_user_id"`
	FromUsername    string `json:"from_username,omitempty"`
	FromDisplayName string `json:"from_display_name,omitempty"`
	FromAddress     string `json:"from_address"`
	AmountWei       string `json:"amount_wei"`
	Message         string `json:"message,omitempty"`
	TxHash          string `json:"tx_hash"`
	CreatedAt       string `json:"created_at"`
}

type DonationsListResponse struct {
	Items []DonationDTO `json:"items"`
	Total int64         `json:"total"`
}
