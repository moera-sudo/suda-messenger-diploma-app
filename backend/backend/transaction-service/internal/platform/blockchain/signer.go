package blockchain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Signer — абстракция «кто и как подписывает транзакцию».
//
// Реализации:
//
//	TreasurySigner — приватник казначейства в памяти (mint, expire, refund admin ops)
//	UserSigner     — приватник юзера/канала, дешифрованный из БД на лету
//
// SignTx ставит подпись на готовую *types.Transaction (chainID должен
// совпадать с client.ChainID() — это проверяется внутри).
type Signer interface {
	Address() common.Address
	SignTx(ctx context.Context, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error)

	// TransactOpts собирает *bind.TransactOpts для типизированных вызовов
	// контрактов через abigen-биндинги. Сервис выставляет nonce/gas
	// извне (через NonceManager), здесь только подпись.
	TransactOpts(ctx context.Context, chainID *big.Int) (*bind.TransactOpts, error)
}

// ────────────────────────────────────────────────────────────
//  TreasurySigner — единственный экземпляр на весь процесс.
//  Используется для onlyOwner операций смарт-контрактов (mint новых SUDA,
//  изначальный welcome-bonus новому юзеру, expire/cancel просроченных
//  квестов от лица сервиса).
// ────────────────────────────────────────────────────────────

type TreasurySigner struct {
	privKey *ecdsa.PrivateKey
	addr    common.Address
}

// NewTreasurySigner парсит hex-приватник (с 0x или без) и валидирует его.
func NewTreasurySigner(hexPrivKey string) (*TreasurySigner, error) {
	if hexPrivKey == "" {
		return nil, fmt.Errorf("blockchain: ETH_PRIVATE_KEY is empty")
	}
	clean := strings.TrimPrefix(hexPrivKey, "0x")

	priv, err := crypto.HexToECDSA(clean)
	if err != nil {
		return nil, fmt.Errorf("blockchain: parse treasury private key: %w", err)
	}

	pubKey, ok := priv.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("blockchain: cannot cast treasury public key")
	}
	addr := crypto.PubkeyToAddress(*pubKey)

	return &TreasurySigner{privKey: priv, addr: addr}, nil
}

func (s *TreasurySigner) Address() common.Address { return s.addr }

func (s *TreasurySigner) SignTx(_ context.Context, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return types.SignTx(tx, types.LatestSignerForChainID(chainID), s.privKey)
}

func (s *TreasurySigner) TransactOpts(_ context.Context, chainID *big.Int) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(s.privKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("blockchain: treasury transact opts: %w", err)
	}
	opts.GasPrice = big.NewInt(0) // Besu zeroBaseFee
	return opts, nil
}

// ────────────────────────────────────────────────────────────
//  UserSigner — конструируется на каждый запрос юзера, живёт ровно
//  одну операцию. Приватник дешифруется из tx_wallets через encryptor.
//  После использования рекомендуется НЕ держать ссылки на UserSigner,
//  чтобы GC побыстрее обнулил приватник в памяти.
// ────────────────────────────────────────────────────────────

type UserSigner struct {
	privKey *ecdsa.PrivateKey
	addr    common.Address
}

// NewUserSigner принимает уже расшифрованный hex-приватник.
// Дешифровка и audit-logging выполняются в feature-слое (wallet service),
// чтобы знать operation/user_id для записи в tx_signing_audit.
func NewUserSigner(hexPrivKey string) (*UserSigner, error) {
	clean := strings.TrimPrefix(hexPrivKey, "0x")
	priv, err := crypto.HexToECDSA(clean)
	if err != nil {
		return nil, fmt.Errorf("blockchain: parse user private key: %w", err)
	}
	pubKey, ok := priv.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("blockchain: cannot cast user public key")
	}
	addr := crypto.PubkeyToAddress(*pubKey)
	return &UserSigner{privKey: priv, addr: addr}, nil
}

func (s *UserSigner) Address() common.Address { return s.addr }

func (s *UserSigner) SignTx(_ context.Context, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return types.SignTx(tx, types.LatestSignerForChainID(chainID), s.privKey)
}

func (s *UserSigner) TransactOpts(_ context.Context, chainID *big.Int) (*bind.TransactOpts, error) {
	opts, err := bind.NewKeyedTransactorWithChainID(s.privKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("blockchain: user transact opts: %w", err)
	}
	opts.GasPrice = big.NewInt(0)
	return opts, nil
}

// GenerateNewWallet создаёт новую пару ключей secp256k1 и возвращает
// hex-приватник (без 0x) и EVM-адрес. Используется при регистрации
// нового юзера / канала.
func GenerateNewWallet() (hexPrivKey string, addr common.Address, err error) {
	priv, err := crypto.GenerateKey()
	if err != nil {
		return "", common.Address{}, fmt.Errorf("blockchain: generate key: %w", err)
	}
	pubKey, ok := priv.Public().(*ecdsa.PublicKey)
	if !ok {
		return "", common.Address{}, fmt.Errorf("blockchain: cast public key")
	}
	addr = crypto.PubkeyToAddress(*pubKey)

	// crypto.FromECDSA возвращает 32-байтовый приватник.
	hexPrivKey = common.Bytes2Hex(crypto.FromECDSA(priv))
	return hexPrivKey, addr, nil
}