package purchase

import (
	"math/big"
)

// Package — фиксированный "пакет SUDA", доступный для покупки.
//
// Захардкожен в коде (не в БД): набор небольшой, меняется редко,
// и любая попытка эксплойта через DB-injection сразу будет видна
// при ревью кода.
type Package struct {
	Code         string   // машинный идентификатор: "SMALL"
	Title        string   // отображаемый заголовок для UI
	SudaAmount   int64    // количество SUDA (не wei)
	FiatPrice    string   // десятичная строка цены: "1.99"
	FiatCurrency string   // ISO 4217: "USD"
}

// SudaAmountWei возвращает количество SUDA, переведённое в wei (1 SUDA = 1e18 wei).
func (p Package) SudaAmountWei() *big.Int {
	weiPerSuda := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	return new(big.Int).Mul(big.NewInt(p.SudaAmount), weiPerSuda)
}

// Packages — список доступных пакетов. Порядок отражает то, что увидит UI.
var Packages = []Package{
	{Code: "SMALL", Title: "Small bag", SudaAmount: 100, FiatPrice: "1.99", FiatCurrency: "USD"},
	{Code: "MEDIUM", Title: "Medium bag", SudaAmount: 500, FiatPrice: "8.99", FiatCurrency: "USD"},
	{Code: "LARGE", Title: "Large bag", SudaAmount: 1000, FiatPrice: "14.99", FiatCurrency: "USD"},
	{Code: "MEGA", Title: "Whale bag", SudaAmount: 5000, FiatPrice: "49.99", FiatCurrency: "USD"},
}

// FindPackage возвращает пакет по его коду или (Package{}, false).
func FindPackage(code string) (Package, bool) {
	for _, p := range Packages {
		if p.Code == code {
			return p, true
		}
	}
	return Package{}, false
}

// ValidPaymentMethod — true для известных способов оплаты.
// Поле декоративное, поэтому валидируем мягко.
func ValidPaymentMethod(m string) bool {
	switch m {
	case PaymentMethodCard, PaymentMethodApplePay, PaymentMethodGooglePay:
		return true
	default:
		return false
	}
}
