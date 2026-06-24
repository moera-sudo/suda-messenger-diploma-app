package http

// POST /donate
//
// Ровно одно из to_username / to_channel_id должно быть задано.
// amount_wei — строка десятичного числа в wei (uint256 не лезет в JSON number).
type DonateRequest struct {
	ToUsername  string `json:"to_username,omitempty"   validate:"omitempty,min=1,max=50"`
	ToChannelID string `json:"to_channel_id,omitempty" validate:"omitempty,uuid"`
	AmountWei   string `json:"amount_wei"              validate:"required"`
	Message     string `json:"message,omitempty"       validate:"max=200"`
	ChatID      string `json:"chat_id,omitempty"       validate:"omitempty,uuid"`
}

// DonateResponse — async-ответ: tx_hash сразу, статус через WS-event.
type DonateResponse struct {
	TxHash      string `json:"tx_hash"`
	FromAddress string `json:"from_address"`
	ToAddress   string `json:"to_address"`
	AmountWei   string `json:"amount_wei"`
	IsChannel   bool   `json:"is_channel"`
	ToUserID    string `json:"to_user_id,omitempty"`
	ToChannelID string `json:"to_channel_id,omitempty"`
	SubmittedAt string `json:"submitted_at"`
}
