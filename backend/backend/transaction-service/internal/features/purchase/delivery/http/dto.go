package http

// GET /purchase/packages

// PackageDTO — то, что покажем юзеру в магазине.
type PackageDTO struct {
	Code          string `json:"code"`
	Title         string `json:"title"`
	SudaAmount    int64  `json:"suda_amount"`
	SudaAmountWei string `json:"suda_amount_wei"`
	FiatPrice     string `json:"fiat_price"`
	FiatCurrency  string `json:"fiat_currency"`
}

type PackagesResponse struct {
	Packages []PackageDTO `json:"packages"`
}

// POST /purchase/initiate

type InitiateRequest struct {
	PackageCode   string `json:"package_code"   validate:"required,min=1,max=50"`
	PaymentMethod string `json:"payment_method" validate:"omitempty,oneof=CARD APPLE_PAY GOOGLE_PAY"`
}

type InitiateResponse struct {
	PurchaseID    string `json:"purchase_id"`
	PackageCode   string `json:"package_code"`
	SudaAmountWei string `json:"suda_amount_wei"`
	FiatAmount    string `json:"fiat_amount"`
	FiatCurrency  string `json:"fiat_currency"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
}

// POST /purchase/:id/confirm
//
// card_number/cvv намеренно НЕ валидируются (это имитация). Поля есть
// в DTO, чтобы клиент мог отправить их и UI выглядел как настоящий чекаут.
type ConfirmRequest struct {
	CardNumber string `json:"card_number,omitempty"`
	CVV        string `json:"cvv,omitempty"`
	CardholderName string `json:"cardholder_name,omitempty"`
	ExpiryMonth string `json:"expiry_month,omitempty"`
	ExpiryYear  string `json:"expiry_year,omitempty"`
}

type ConfirmResponse struct {
	PurchaseID    string `json:"purchase_id"`
	Status        string `json:"status"`
	TxHash        string `json:"tx_hash,omitempty"`
	SudaAmountWei string `json:"suda_amount_wei"`
	CompletedAt   string `json:"completed_at"`
}

// GET /purchase/history

type HistoryItem struct {
	ID            string `json:"id"`
	PackageCode   string `json:"package_code"`
	SudaAmountWei string `json:"suda_amount_wei"`
	FiatAmount    string `json:"fiat_amount"`
	FiatCurrency  string `json:"fiat_currency"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method,omitempty"`
	TxHash        string `json:"tx_hash,omitempty"`
	FailureReason string `json:"failure_reason,omitempty"`
	CreatedAt     string `json:"created_at"`
	CompletedAt   string `json:"completed_at,omitempty"`
}

type HistoryResponse struct {
	Items []HistoryItem `json:"items"`
}
