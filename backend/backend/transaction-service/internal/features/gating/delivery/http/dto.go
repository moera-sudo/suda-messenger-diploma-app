package http

// POST /gating/rule
//
// min_suda_balance_wei — строка десятичного числа в wei ("0" = без баланс-порога).
// required_nft_collection_id — опционально (NFT-фича пока stub; на практике null).
type CreateRuleRequest struct {
	ChatID string `json:"chat_id" validate:"required,uuid"`
	// subscription_price_wei — цена платной подписки на PUBLIC-канал в wei
	// ("0" или пусто = канал не платный). Главное поле для платной подписки.
	SubscriptionPriceWei string `json:"subscription_price_wei,omitempty"`
	// min_suda_balance_wei — legacy hold-порог; опционален ("" = "0").
	MinSudaBalanceWei       string `json:"min_suda_balance_wei,omitempty"`
	RequiredNFTCollectionID string `json:"required_nft_collection_id,omitempty" validate:"omitempty,uuid"`
}

// RuleResponse — текущее правило (GET /gating/rule/:chat_id и ответ на POST).
type RuleResponse struct {
	ChatID                  string `json:"chat_id"`
	SubscriptionPriceWei    string `json:"subscription_price_wei"`
	MinSudaBalanceWei       string `json:"min_suda_balance_wei"`
	RequiredNFTCollectionID string `json:"required_nft_collection_id,omitempty"`
	HasRule                 bool   `json:"has_rule"`
}
