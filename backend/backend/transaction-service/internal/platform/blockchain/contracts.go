package blockchain

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/moera-sudo/backend/backend/transaction-service/contracts/sudaescrow"
	"github.com/moera-sudo/backend/backend/transaction-service/contracts/sudafundraising"
	"github.com/moera-sudo/backend/backend/transaction-service/contracts/sudamarketplace"
	"github.com/moera-sudo/backend/backend/transaction-service/contracts/sudanft"
	"github.com/moera-sudo/backend/backend/transaction-service/contracts/sudatoken"
)

// Contracts — централизованный контейнер инстансов всех 5 задеплоенных
// смарт-контрактов. Создаётся один раз при старте сервиса. Каждое поле —
// типизированная обёртка от abigen, которую можно дальше использовать
// для view-вызовов и формирования транзакций.
//
// Адреса самих контрактов передаются через config из .env и должны
// соответствовать deployments/<chainId>.json после `npx hardhat run deploy`.
type Contracts struct {
	Token       *sudatoken.SudaToken
	NFT         *sudanft.SudaNFT
	Marketplace *sudamarketplace.SudaMarketplace
	Escrow      *sudaescrow.SudaEscrow
	Fundraising *sudafundraising.SudaFundraising

	// Адреса контрактов в строковой форме — пригодятся для логов и
	// сравнений (например, в observer'е чтобы понять "к чему относится event").
	TokenAddr       common.Address
	NFTAddr         common.Address
	MarketplaceAddr common.Address
	EscrowAddr      common.Address
	FundraisingAddr common.Address
}

// ContractAddresses — адреса для NewContracts. Все 5 обязательны;
// если хоть один пустой — конструктор вернёт ошибку.
type ContractAddresses struct {
	Token       string
	NFT         string
	Marketplace string
	Escrow      string
	Fundraising string
}

// NewContracts создаёт типизированные обёртки над всеми 5 контрактами.
// Если какой-то адрес невалиден или не задан — возвращает ошибку,
// чтобы сервис не стартанул с битой конфигурацией.
func NewContracts(client *Client, addrs ContractAddresses) (*Contracts, error) {
	if err := requireAddr("SUDA_TOKEN_ADDRESS", addrs.Token); err != nil {
		return nil, err
	}
	if err := requireAddr("SUDA_NFT_ADDRESS", addrs.NFT); err != nil {
		return nil, err
	}
	if err := requireAddr("SUDA_MARKETPLACE_ADDRESS", addrs.Marketplace); err != nil {
		return nil, err
	}
	if err := requireAddr("SUDA_ESCROW_ADDRESS", addrs.Escrow); err != nil {
		return nil, err
	}
	if err := requireAddr("SUDA_FUNDRAISING_ADDRESS", addrs.Fundraising); err != nil {
		return nil, err
	}

	tokenAddr := common.HexToAddress(addrs.Token)
	token, err := sudatoken.NewSudaToken(tokenAddr, client.Eth())
	if err != nil {
		return nil, fmt.Errorf("blockchain: bind SudaToken: %w", err)
	}

	nftAddr := common.HexToAddress(addrs.NFT)
	nft, err := sudanft.NewSudaNFT(nftAddr, client.Eth())
	if err != nil {
		return nil, fmt.Errorf("blockchain: bind SudaNFT: %w", err)
	}

	marketAddr := common.HexToAddress(addrs.Marketplace)
	market, err := sudamarketplace.NewSudaMarketplace(marketAddr, client.Eth())
	if err != nil {
		return nil, fmt.Errorf("blockchain: bind SudaMarketplace: %w", err)
	}

	escrowAddr := common.HexToAddress(addrs.Escrow)
	escrow, err := sudaescrow.NewSudaEscrow(escrowAddr, client.Eth())
	if err != nil {
		return nil, fmt.Errorf("blockchain: bind SudaEscrow: %w", err)
	}

	fundAddr := common.HexToAddress(addrs.Fundraising)
	fund, err := sudafundraising.NewSudaFundraising(fundAddr, client.Eth())
	if err != nil {
		return nil, fmt.Errorf("blockchain: bind SudaFundraising: %w", err)
	}

	return &Contracts{
		Token:           token,
		NFT:             nft,
		Marketplace:     market,
		Escrow:          escrow,
		Fundraising:     fund,
		TokenAddr:       tokenAddr,
		NFTAddr:         nftAddr,
		MarketplaceAddr: marketAddr,
		EscrowAddr:      escrowAddr,
		FundraisingAddr: fundAddr,
	}, nil
}

func requireAddr(envName, value string) error {
	if value == "" {
		return fmt.Errorf("blockchain: %s is empty (deploy contracts and set in .env)", envName)
	}
	if !common.IsHexAddress(value) {
		return fmt.Errorf("blockchain: %s is not a valid EVM address: %s", envName, value)
	}
	return nil
}