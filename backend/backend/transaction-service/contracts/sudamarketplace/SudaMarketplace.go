// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sudamarketplace

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// SudaMarketplaceMetaData contains all meta data concerning the SudaMarketplace contract.
var SudaMarketplaceMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_sudaToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"}],\"name\":\"Cancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"nftContract\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"Listed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldBps\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBps\",\"type\":\"uint256\"}],\"name\":\"MarketplaceFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"buyer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"royaltyAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"marketFee\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"sellerProceeds\",\"type\":\"uint256\"}],\"name\":\"Sold\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAX_FEE_BPS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"}],\"name\":\"buy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"}],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"nftContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"}],\"name\":\"list\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"listingId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"listings\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"seller\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"nftContract\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"price\",\"type\":\"uint256\"},{\"internalType\":\"enumSudaMarketplace.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"marketplaceFeeBps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextListingId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"onERC721Received\",\"outputs\":[{\"internalType\":\"bytes4\",\"name\":\"\",\"type\":\"bytes4\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newBps\",\"type\":\"uint256\"}],\"name\":\"setMarketplaceFeeBps\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sudaToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// SudaMarketplaceABI is the input ABI used to generate the binding from.
// Deprecated: Use SudaMarketplaceMetaData.ABI instead.
var SudaMarketplaceABI = SudaMarketplaceMetaData.ABI

// SudaMarketplace is an auto generated Go binding around an Ethereum contract.
type SudaMarketplace struct {
	SudaMarketplaceCaller     // Read-only binding to the contract
	SudaMarketplaceTransactor // Write-only binding to the contract
	SudaMarketplaceFilterer   // Log filterer for contract events
}

// SudaMarketplaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type SudaMarketplaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SudaMarketplaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SudaMarketplaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SudaMarketplaceFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SudaMarketplaceFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SudaMarketplaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SudaMarketplaceSession struct {
	Contract     *SudaMarketplace  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SudaMarketplaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SudaMarketplaceCallerSession struct {
	Contract *SudaMarketplaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// SudaMarketplaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SudaMarketplaceTransactorSession struct {
	Contract     *SudaMarketplaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// SudaMarketplaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type SudaMarketplaceRaw struct {
	Contract *SudaMarketplace // Generic contract binding to access the raw methods on
}

// SudaMarketplaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SudaMarketplaceCallerRaw struct {
	Contract *SudaMarketplaceCaller // Generic read-only contract binding to access the raw methods on
}

// SudaMarketplaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SudaMarketplaceTransactorRaw struct {
	Contract *SudaMarketplaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSudaMarketplace creates a new instance of SudaMarketplace, bound to a specific deployed contract.
func NewSudaMarketplace(address common.Address, backend bind.ContractBackend) (*SudaMarketplace, error) {
	contract, err := bindSudaMarketplace(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SudaMarketplace{SudaMarketplaceCaller: SudaMarketplaceCaller{contract: contract}, SudaMarketplaceTransactor: SudaMarketplaceTransactor{contract: contract}, SudaMarketplaceFilterer: SudaMarketplaceFilterer{contract: contract}}, nil
}

// NewSudaMarketplaceCaller creates a new read-only instance of SudaMarketplace, bound to a specific deployed contract.
func NewSudaMarketplaceCaller(address common.Address, caller bind.ContractCaller) (*SudaMarketplaceCaller, error) {
	contract, err := bindSudaMarketplace(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SudaMarketplaceCaller{contract: contract}, nil
}

// NewSudaMarketplaceTransactor creates a new write-only instance of SudaMarketplace, bound to a specific deployed contract.
func NewSudaMarketplaceTransactor(address common.Address, transactor bind.ContractTransactor) (*SudaMarketplaceTransactor, error) {
	contract, err := bindSudaMarketplace(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SudaMarketplaceTransactor{contract: contract}, nil
}

// NewSudaMarketplaceFilterer creates a new log filterer instance of SudaMarketplace, bound to a specific deployed contract.
func NewSudaMarketplaceFilterer(address common.Address, filterer bind.ContractFilterer) (*SudaMarketplaceFilterer, error) {
	contract, err := bindSudaMarketplace(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SudaMarketplaceFilterer{contract: contract}, nil
}

// bindSudaMarketplace binds a generic wrapper to an already deployed contract.
func bindSudaMarketplace(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SudaMarketplaceMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SudaMarketplace *SudaMarketplaceRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SudaMarketplace.Contract.SudaMarketplaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SudaMarketplace *SudaMarketplaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.SudaMarketplaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SudaMarketplace *SudaMarketplaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.SudaMarketplaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SudaMarketplace *SudaMarketplaceCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SudaMarketplace.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SudaMarketplace *SudaMarketplaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SudaMarketplace *SudaMarketplaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.contract.Transact(opts, method, params...)
}

// MAXFEEBPS is a free data retrieval call binding the contract method 0xd55be8c6.
//
// Solidity: function MAX_FEE_BPS() view returns(uint256)
func (_SudaMarketplace *SudaMarketplaceCaller) MAXFEEBPS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SudaMarketplace.contract.Call(opts, &out, "MAX_FEE_BPS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXFEEBPS is a free data retrieval call binding the contract method 0xd55be8c6.
//
// Solidity: function MAX_FEE_BPS() view returns(uint256)
func (_SudaMarketplace *SudaMarketplaceSession) MAXFEEBPS() (*big.Int, error) {
	return _SudaMarketplace.Contract.MAXFEEBPS(&_SudaMarketplace.CallOpts)
}

// MAXFEEBPS is a free data retrieval call binding the contract method 0xd55be8c6.
//
// Solidity: function MAX_FEE_BPS() view returns(uint256)
func (_SudaMarketplace *SudaMarketplaceCallerSession) MAXFEEBPS() (*big.Int, error) {
	return _SudaMarketplace.Contract.MAXFEEBPS(&_SudaMarketplace.CallOpts)
}

// Listings is a free data retrieval call binding the contract method 0xde74e57b.
//
// Solidity: function listings(uint256 ) view returns(address seller, address nftContract, uint256 tokenId, uint256 price, uint8 status)
func (_SudaMarketplace *SudaMarketplaceCaller) Listings(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Seller      common.Address
	NftContract common.Address
	TokenId     *big.Int
	Price       *big.Int
	Status      uint8
}, error) {
	var out []interface{}
	err := _SudaMarketplace.contract.Call(opts, &out, "listings", arg0)

	outstruct := new(struct {
		Seller      common.Address
		NftContract common.Address
		TokenId     *big.Int
		Price       *big.Int
		Status      uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Seller = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.NftContract = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.TokenId = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Price = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Status = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

// Listings is a free data retrieval call binding the contract method 0xde74e57b.
//
// Solidity: function listings(uint256 ) view returns(address seller, address nftContract, uint256 tokenId, uint256 price, uint8 status)
func (_SudaMarketplace *SudaMarketplaceSession) Listings(arg0 *big.Int) (struct {
	Seller      common.Address
	NftContract common.Address
	TokenId     *big.Int
	Price       *big.Int
	Status      uint8
}, error) {
	return _SudaMarketplace.Contract.Listings(&_SudaMarketplace.CallOpts, arg0)
}

// Listings is a free data retrieval call binding the contract method 0xde74e57b.
//
// Solidity: function listings(uint256 ) view returns(address seller, address nftContract, uint256 tokenId, uint256 price, uint8 status)
func (_SudaMarketplace *SudaMarketplaceCallerSession) Listings(arg0 *big.Int) (struct {
	Seller      common.Address
	NftContract common.Address
	TokenId     *big.Int
	Price       *big.Int
	Status      uint8
}, error) {
	return _SudaMarketplace.Contract.Listings(&_SudaMarketplace.CallOpts, arg0)
}

// MarketplaceFeeBps is a free data retrieval call binding the contract method 0xd3929c8a.
//
// Solidity: function marketplaceFeeBps() view returns(uint256)
func (_SudaMarketplace *SudaMarketplaceCaller) MarketplaceFeeBps(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SudaMarketplace.contract.Call(opts, &out, "marketplaceFeeBps")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MarketplaceFeeBps is a free data retrieval call binding the contract method 0xd3929c8a.
//
// Solidity: function marketplaceFeeBps() view returns(uint256)
func (_SudaMarketplace *SudaMarketplaceSession) MarketplaceFeeBps() (*big.Int, error) {
	return _SudaMarketplace.Contract.MarketplaceFeeBps(&_SudaMarketplace.CallOpts)
}

// MarketplaceFeeBps is a free data retrieval call binding the contract method 0xd3929c8a.
//
// Solidity: function marketplaceFeeBps() view returns(uint256)
func (_SudaMarketplace *SudaMarketplaceCallerSession) MarketplaceFeeBps() (*big.Int, error) {
	return _SudaMarketplace.Contract.MarketplaceFeeBps(&_SudaMarketplace.CallOpts)
}

// NextListingId is a free data retrieval call binding the contract method 0xaaccf1ec.
//
// Solidity: function nextListingId() view returns(uint256)
func (_SudaMarketplace *SudaMarketplaceCaller) NextListingId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SudaMarketplace.contract.Call(opts, &out, "nextListingId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextListingId is a free data retrieval call binding the contract method 0xaaccf1ec.
//
// Solidity: function nextListingId() view returns(uint256)
func (_SudaMarketplace *SudaMarketplaceSession) NextListingId() (*big.Int, error) {
	return _SudaMarketplace.Contract.NextListingId(&_SudaMarketplace.CallOpts)
}

// NextListingId is a free data retrieval call binding the contract method 0xaaccf1ec.
//
// Solidity: function nextListingId() view returns(uint256)
func (_SudaMarketplace *SudaMarketplaceCallerSession) NextListingId() (*big.Int, error) {
	return _SudaMarketplace.Contract.NextListingId(&_SudaMarketplace.CallOpts)
}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_SudaMarketplace *SudaMarketplaceCaller) OnERC721Received(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	var out []interface{}
	err := _SudaMarketplace.contract.Call(opts, &out, "onERC721Received", arg0, arg1, arg2, arg3)

	if err != nil {
		return *new([4]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([4]byte)).(*[4]byte)

	return out0, err

}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_SudaMarketplace *SudaMarketplaceSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	return _SudaMarketplace.Contract.OnERC721Received(&_SudaMarketplace.CallOpts, arg0, arg1, arg2, arg3)
}

// OnERC721Received is a free data retrieval call binding the contract method 0x150b7a02.
//
// Solidity: function onERC721Received(address , address , uint256 , bytes ) pure returns(bytes4)
func (_SudaMarketplace *SudaMarketplaceCallerSession) OnERC721Received(arg0 common.Address, arg1 common.Address, arg2 *big.Int, arg3 []byte) ([4]byte, error) {
	return _SudaMarketplace.Contract.OnERC721Received(&_SudaMarketplace.CallOpts, arg0, arg1, arg2, arg3)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SudaMarketplace *SudaMarketplaceCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SudaMarketplace.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SudaMarketplace *SudaMarketplaceSession) Owner() (common.Address, error) {
	return _SudaMarketplace.Contract.Owner(&_SudaMarketplace.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SudaMarketplace *SudaMarketplaceCallerSession) Owner() (common.Address, error) {
	return _SudaMarketplace.Contract.Owner(&_SudaMarketplace.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SudaMarketplace *SudaMarketplaceCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SudaMarketplace.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SudaMarketplace *SudaMarketplaceSession) Paused() (bool, error) {
	return _SudaMarketplace.Contract.Paused(&_SudaMarketplace.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SudaMarketplace *SudaMarketplaceCallerSession) Paused() (bool, error) {
	return _SudaMarketplace.Contract.Paused(&_SudaMarketplace.CallOpts)
}

// SudaToken is a free data retrieval call binding the contract method 0xda71113f.
//
// Solidity: function sudaToken() view returns(address)
func (_SudaMarketplace *SudaMarketplaceCaller) SudaToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SudaMarketplace.contract.Call(opts, &out, "sudaToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SudaToken is a free data retrieval call binding the contract method 0xda71113f.
//
// Solidity: function sudaToken() view returns(address)
func (_SudaMarketplace *SudaMarketplaceSession) SudaToken() (common.Address, error) {
	return _SudaMarketplace.Contract.SudaToken(&_SudaMarketplace.CallOpts)
}

// SudaToken is a free data retrieval call binding the contract method 0xda71113f.
//
// Solidity: function sudaToken() view returns(address)
func (_SudaMarketplace *SudaMarketplaceCallerSession) SudaToken() (common.Address, error) {
	return _SudaMarketplace.Contract.SudaToken(&_SudaMarketplace.CallOpts)
}

// Buy is a paid mutator transaction binding the contract method 0xd96a094a.
//
// Solidity: function buy(uint256 listingId) returns()
func (_SudaMarketplace *SudaMarketplaceTransactor) Buy(opts *bind.TransactOpts, listingId *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.contract.Transact(opts, "buy", listingId)
}

// Buy is a paid mutator transaction binding the contract method 0xd96a094a.
//
// Solidity: function buy(uint256 listingId) returns()
func (_SudaMarketplace *SudaMarketplaceSession) Buy(listingId *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.Buy(&_SudaMarketplace.TransactOpts, listingId)
}

// Buy is a paid mutator transaction binding the contract method 0xd96a094a.
//
// Solidity: function buy(uint256 listingId) returns()
func (_SudaMarketplace *SudaMarketplaceTransactorSession) Buy(listingId *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.Buy(&_SudaMarketplace.TransactOpts, listingId)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 listingId) returns()
func (_SudaMarketplace *SudaMarketplaceTransactor) Cancel(opts *bind.TransactOpts, listingId *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.contract.Transact(opts, "cancel", listingId)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 listingId) returns()
func (_SudaMarketplace *SudaMarketplaceSession) Cancel(listingId *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.Cancel(&_SudaMarketplace.TransactOpts, listingId)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 listingId) returns()
func (_SudaMarketplace *SudaMarketplaceTransactorSession) Cancel(listingId *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.Cancel(&_SudaMarketplace.TransactOpts, listingId)
}

// List is a paid mutator transaction binding the contract method 0xdda342bb.
//
// Solidity: function list(address nftContract, uint256 tokenId, uint256 price) returns(uint256 listingId)
func (_SudaMarketplace *SudaMarketplaceTransactor) List(opts *bind.TransactOpts, nftContract common.Address, tokenId *big.Int, price *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.contract.Transact(opts, "list", nftContract, tokenId, price)
}

// List is a paid mutator transaction binding the contract method 0xdda342bb.
//
// Solidity: function list(address nftContract, uint256 tokenId, uint256 price) returns(uint256 listingId)
func (_SudaMarketplace *SudaMarketplaceSession) List(nftContract common.Address, tokenId *big.Int, price *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.List(&_SudaMarketplace.TransactOpts, nftContract, tokenId, price)
}

// List is a paid mutator transaction binding the contract method 0xdda342bb.
//
// Solidity: function list(address nftContract, uint256 tokenId, uint256 price) returns(uint256 listingId)
func (_SudaMarketplace *SudaMarketplaceTransactorSession) List(nftContract common.Address, tokenId *big.Int, price *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.List(&_SudaMarketplace.TransactOpts, nftContract, tokenId, price)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SudaMarketplace *SudaMarketplaceTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaMarketplace.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SudaMarketplace *SudaMarketplaceSession) Pause() (*types.Transaction, error) {
	return _SudaMarketplace.Contract.Pause(&_SudaMarketplace.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SudaMarketplace *SudaMarketplaceTransactorSession) Pause() (*types.Transaction, error) {
	return _SudaMarketplace.Contract.Pause(&_SudaMarketplace.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SudaMarketplace *SudaMarketplaceTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaMarketplace.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SudaMarketplace *SudaMarketplaceSession) RenounceOwnership() (*types.Transaction, error) {
	return _SudaMarketplace.Contract.RenounceOwnership(&_SudaMarketplace.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SudaMarketplace *SudaMarketplaceTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SudaMarketplace.Contract.RenounceOwnership(&_SudaMarketplace.TransactOpts)
}

// SetMarketplaceFeeBps is a paid mutator transaction binding the contract method 0x0b43e840.
//
// Solidity: function setMarketplaceFeeBps(uint256 newBps) returns()
func (_SudaMarketplace *SudaMarketplaceTransactor) SetMarketplaceFeeBps(opts *bind.TransactOpts, newBps *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.contract.Transact(opts, "setMarketplaceFeeBps", newBps)
}

// SetMarketplaceFeeBps is a paid mutator transaction binding the contract method 0x0b43e840.
//
// Solidity: function setMarketplaceFeeBps(uint256 newBps) returns()
func (_SudaMarketplace *SudaMarketplaceSession) SetMarketplaceFeeBps(newBps *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.SetMarketplaceFeeBps(&_SudaMarketplace.TransactOpts, newBps)
}

// SetMarketplaceFeeBps is a paid mutator transaction binding the contract method 0x0b43e840.
//
// Solidity: function setMarketplaceFeeBps(uint256 newBps) returns()
func (_SudaMarketplace *SudaMarketplaceTransactorSession) SetMarketplaceFeeBps(newBps *big.Int) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.SetMarketplaceFeeBps(&_SudaMarketplace.TransactOpts, newBps)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SudaMarketplace *SudaMarketplaceTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SudaMarketplace.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SudaMarketplace *SudaMarketplaceSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.TransferOwnership(&_SudaMarketplace.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SudaMarketplace *SudaMarketplaceTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SudaMarketplace.Contract.TransferOwnership(&_SudaMarketplace.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SudaMarketplace *SudaMarketplaceTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaMarketplace.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SudaMarketplace *SudaMarketplaceSession) Unpause() (*types.Transaction, error) {
	return _SudaMarketplace.Contract.Unpause(&_SudaMarketplace.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SudaMarketplace *SudaMarketplaceTransactorSession) Unpause() (*types.Transaction, error) {
	return _SudaMarketplace.Contract.Unpause(&_SudaMarketplace.TransactOpts)
}

// SudaMarketplaceCancelledIterator is returned from FilterCancelled and is used to iterate over the raw logs and unpacked data for Cancelled events raised by the SudaMarketplace contract.
type SudaMarketplaceCancelledIterator struct {
	Event *SudaMarketplaceCancelled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SudaMarketplaceCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaMarketplaceCancelled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SudaMarketplaceCancelled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SudaMarketplaceCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaMarketplaceCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaMarketplaceCancelled represents a Cancelled event raised by the SudaMarketplace contract.
type SudaMarketplaceCancelled struct {
	ListingId *big.Int
	Seller    common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCancelled is a free log retrieval operation binding the contract event 0x26deca31ff8139a06c52453ce8985d34f7648a6d9af1d283c4063d052c355a0f.
//
// Solidity: event Cancelled(uint256 indexed listingId, address indexed seller)
func (_SudaMarketplace *SudaMarketplaceFilterer) FilterCancelled(opts *bind.FilterOpts, listingId []*big.Int, seller []common.Address) (*SudaMarketplaceCancelledIterator, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _SudaMarketplace.contract.FilterLogs(opts, "Cancelled", listingIdRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return &SudaMarketplaceCancelledIterator{contract: _SudaMarketplace.contract, event: "Cancelled", logs: logs, sub: sub}, nil
}

// WatchCancelled is a free log subscription operation binding the contract event 0x26deca31ff8139a06c52453ce8985d34f7648a6d9af1d283c4063d052c355a0f.
//
// Solidity: event Cancelled(uint256 indexed listingId, address indexed seller)
func (_SudaMarketplace *SudaMarketplaceFilterer) WatchCancelled(opts *bind.WatchOpts, sink chan<- *SudaMarketplaceCancelled, listingId []*big.Int, seller []common.Address) (event.Subscription, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _SudaMarketplace.contract.WatchLogs(opts, "Cancelled", listingIdRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaMarketplaceCancelled)
				if err := _SudaMarketplace.contract.UnpackLog(event, "Cancelled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCancelled is a log parse operation binding the contract event 0x26deca31ff8139a06c52453ce8985d34f7648a6d9af1d283c4063d052c355a0f.
//
// Solidity: event Cancelled(uint256 indexed listingId, address indexed seller)
func (_SudaMarketplace *SudaMarketplaceFilterer) ParseCancelled(log types.Log) (*SudaMarketplaceCancelled, error) {
	event := new(SudaMarketplaceCancelled)
	if err := _SudaMarketplace.contract.UnpackLog(event, "Cancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaMarketplaceListedIterator is returned from FilterListed and is used to iterate over the raw logs and unpacked data for Listed events raised by the SudaMarketplace contract.
type SudaMarketplaceListedIterator struct {
	Event *SudaMarketplaceListed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SudaMarketplaceListedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaMarketplaceListed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SudaMarketplaceListed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SudaMarketplaceListedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaMarketplaceListedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaMarketplaceListed represents a Listed event raised by the SudaMarketplace contract.
type SudaMarketplaceListed struct {
	ListingId   *big.Int
	Seller      common.Address
	NftContract common.Address
	TokenId     *big.Int
	Price       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterListed is a free log retrieval operation binding the contract event 0x723f73331eaee88eec7fc68ef60ab6ed15e4b90d0472b55eb92fa43910bab6dd.
//
// Solidity: event Listed(uint256 indexed listingId, address indexed seller, address indexed nftContract, uint256 tokenId, uint256 price)
func (_SudaMarketplace *SudaMarketplaceFilterer) FilterListed(opts *bind.FilterOpts, listingId []*big.Int, seller []common.Address, nftContract []common.Address) (*SudaMarketplaceListedIterator, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}
	var nftContractRule []interface{}
	for _, nftContractItem := range nftContract {
		nftContractRule = append(nftContractRule, nftContractItem)
	}

	logs, sub, err := _SudaMarketplace.contract.FilterLogs(opts, "Listed", listingIdRule, sellerRule, nftContractRule)
	if err != nil {
		return nil, err
	}
	return &SudaMarketplaceListedIterator{contract: _SudaMarketplace.contract, event: "Listed", logs: logs, sub: sub}, nil
}

// WatchListed is a free log subscription operation binding the contract event 0x723f73331eaee88eec7fc68ef60ab6ed15e4b90d0472b55eb92fa43910bab6dd.
//
// Solidity: event Listed(uint256 indexed listingId, address indexed seller, address indexed nftContract, uint256 tokenId, uint256 price)
func (_SudaMarketplace *SudaMarketplaceFilterer) WatchListed(opts *bind.WatchOpts, sink chan<- *SudaMarketplaceListed, listingId []*big.Int, seller []common.Address, nftContract []common.Address) (event.Subscription, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}
	var nftContractRule []interface{}
	for _, nftContractItem := range nftContract {
		nftContractRule = append(nftContractRule, nftContractItem)
	}

	logs, sub, err := _SudaMarketplace.contract.WatchLogs(opts, "Listed", listingIdRule, sellerRule, nftContractRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaMarketplaceListed)
				if err := _SudaMarketplace.contract.UnpackLog(event, "Listed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseListed is a log parse operation binding the contract event 0x723f73331eaee88eec7fc68ef60ab6ed15e4b90d0472b55eb92fa43910bab6dd.
//
// Solidity: event Listed(uint256 indexed listingId, address indexed seller, address indexed nftContract, uint256 tokenId, uint256 price)
func (_SudaMarketplace *SudaMarketplaceFilterer) ParseListed(log types.Log) (*SudaMarketplaceListed, error) {
	event := new(SudaMarketplaceListed)
	if err := _SudaMarketplace.contract.UnpackLog(event, "Listed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaMarketplaceMarketplaceFeeUpdatedIterator is returned from FilterMarketplaceFeeUpdated and is used to iterate over the raw logs and unpacked data for MarketplaceFeeUpdated events raised by the SudaMarketplace contract.
type SudaMarketplaceMarketplaceFeeUpdatedIterator struct {
	Event *SudaMarketplaceMarketplaceFeeUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SudaMarketplaceMarketplaceFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaMarketplaceMarketplaceFeeUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SudaMarketplaceMarketplaceFeeUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SudaMarketplaceMarketplaceFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaMarketplaceMarketplaceFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaMarketplaceMarketplaceFeeUpdated represents a MarketplaceFeeUpdated event raised by the SudaMarketplace contract.
type SudaMarketplaceMarketplaceFeeUpdated struct {
	OldBps *big.Int
	NewBps *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterMarketplaceFeeUpdated is a free log retrieval operation binding the contract event 0x936008f1e981d6f48d4a529714c27d74d644750a1d215451212ce0f2476f2fbc.
//
// Solidity: event MarketplaceFeeUpdated(uint256 oldBps, uint256 newBps)
func (_SudaMarketplace *SudaMarketplaceFilterer) FilterMarketplaceFeeUpdated(opts *bind.FilterOpts) (*SudaMarketplaceMarketplaceFeeUpdatedIterator, error) {

	logs, sub, err := _SudaMarketplace.contract.FilterLogs(opts, "MarketplaceFeeUpdated")
	if err != nil {
		return nil, err
	}
	return &SudaMarketplaceMarketplaceFeeUpdatedIterator{contract: _SudaMarketplace.contract, event: "MarketplaceFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchMarketplaceFeeUpdated is a free log subscription operation binding the contract event 0x936008f1e981d6f48d4a529714c27d74d644750a1d215451212ce0f2476f2fbc.
//
// Solidity: event MarketplaceFeeUpdated(uint256 oldBps, uint256 newBps)
func (_SudaMarketplace *SudaMarketplaceFilterer) WatchMarketplaceFeeUpdated(opts *bind.WatchOpts, sink chan<- *SudaMarketplaceMarketplaceFeeUpdated) (event.Subscription, error) {

	logs, sub, err := _SudaMarketplace.contract.WatchLogs(opts, "MarketplaceFeeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaMarketplaceMarketplaceFeeUpdated)
				if err := _SudaMarketplace.contract.UnpackLog(event, "MarketplaceFeeUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMarketplaceFeeUpdated is a log parse operation binding the contract event 0x936008f1e981d6f48d4a529714c27d74d644750a1d215451212ce0f2476f2fbc.
//
// Solidity: event MarketplaceFeeUpdated(uint256 oldBps, uint256 newBps)
func (_SudaMarketplace *SudaMarketplaceFilterer) ParseMarketplaceFeeUpdated(log types.Log) (*SudaMarketplaceMarketplaceFeeUpdated, error) {
	event := new(SudaMarketplaceMarketplaceFeeUpdated)
	if err := _SudaMarketplace.contract.UnpackLog(event, "MarketplaceFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaMarketplaceOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SudaMarketplace contract.
type SudaMarketplaceOwnershipTransferredIterator struct {
	Event *SudaMarketplaceOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SudaMarketplaceOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaMarketplaceOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SudaMarketplaceOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SudaMarketplaceOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaMarketplaceOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaMarketplaceOwnershipTransferred represents a OwnershipTransferred event raised by the SudaMarketplace contract.
type SudaMarketplaceOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SudaMarketplace *SudaMarketplaceFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SudaMarketplaceOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SudaMarketplace.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SudaMarketplaceOwnershipTransferredIterator{contract: _SudaMarketplace.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SudaMarketplace *SudaMarketplaceFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SudaMarketplaceOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SudaMarketplace.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaMarketplaceOwnershipTransferred)
				if err := _SudaMarketplace.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SudaMarketplace *SudaMarketplaceFilterer) ParseOwnershipTransferred(log types.Log) (*SudaMarketplaceOwnershipTransferred, error) {
	event := new(SudaMarketplaceOwnershipTransferred)
	if err := _SudaMarketplace.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaMarketplacePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the SudaMarketplace contract.
type SudaMarketplacePausedIterator struct {
	Event *SudaMarketplacePaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SudaMarketplacePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaMarketplacePaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SudaMarketplacePaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SudaMarketplacePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaMarketplacePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaMarketplacePaused represents a Paused event raised by the SudaMarketplace contract.
type SudaMarketplacePaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SudaMarketplace *SudaMarketplaceFilterer) FilterPaused(opts *bind.FilterOpts) (*SudaMarketplacePausedIterator, error) {

	logs, sub, err := _SudaMarketplace.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &SudaMarketplacePausedIterator{contract: _SudaMarketplace.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SudaMarketplace *SudaMarketplaceFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *SudaMarketplacePaused) (event.Subscription, error) {

	logs, sub, err := _SudaMarketplace.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaMarketplacePaused)
				if err := _SudaMarketplace.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SudaMarketplace *SudaMarketplaceFilterer) ParsePaused(log types.Log) (*SudaMarketplacePaused, error) {
	event := new(SudaMarketplacePaused)
	if err := _SudaMarketplace.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaMarketplaceSoldIterator is returned from FilterSold and is used to iterate over the raw logs and unpacked data for Sold events raised by the SudaMarketplace contract.
type SudaMarketplaceSoldIterator struct {
	Event *SudaMarketplaceSold // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SudaMarketplaceSoldIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaMarketplaceSold)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SudaMarketplaceSold)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SudaMarketplaceSoldIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaMarketplaceSoldIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaMarketplaceSold represents a Sold event raised by the SudaMarketplace contract.
type SudaMarketplaceSold struct {
	ListingId      *big.Int
	Buyer          common.Address
	Seller         common.Address
	Price          *big.Int
	RoyaltyAmount  *big.Int
	MarketFee      *big.Int
	SellerProceeds *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSold is a free log retrieval operation binding the contract event 0x6d8b7be78706a9b4a0413c29b82265e6de593cf330075475acbef8c50fed4fb2.
//
// Solidity: event Sold(uint256 indexed listingId, address indexed buyer, address indexed seller, uint256 price, uint256 royaltyAmount, uint256 marketFee, uint256 sellerProceeds)
func (_SudaMarketplace *SudaMarketplaceFilterer) FilterSold(opts *bind.FilterOpts, listingId []*big.Int, buyer []common.Address, seller []common.Address) (*SudaMarketplaceSoldIterator, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _SudaMarketplace.contract.FilterLogs(opts, "Sold", listingIdRule, buyerRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return &SudaMarketplaceSoldIterator{contract: _SudaMarketplace.contract, event: "Sold", logs: logs, sub: sub}, nil
}

// WatchSold is a free log subscription operation binding the contract event 0x6d8b7be78706a9b4a0413c29b82265e6de593cf330075475acbef8c50fed4fb2.
//
// Solidity: event Sold(uint256 indexed listingId, address indexed buyer, address indexed seller, uint256 price, uint256 royaltyAmount, uint256 marketFee, uint256 sellerProceeds)
func (_SudaMarketplace *SudaMarketplaceFilterer) WatchSold(opts *bind.WatchOpts, sink chan<- *SudaMarketplaceSold, listingId []*big.Int, buyer []common.Address, seller []common.Address) (event.Subscription, error) {

	var listingIdRule []interface{}
	for _, listingIdItem := range listingId {
		listingIdRule = append(listingIdRule, listingIdItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}
	var sellerRule []interface{}
	for _, sellerItem := range seller {
		sellerRule = append(sellerRule, sellerItem)
	}

	logs, sub, err := _SudaMarketplace.contract.WatchLogs(opts, "Sold", listingIdRule, buyerRule, sellerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaMarketplaceSold)
				if err := _SudaMarketplace.contract.UnpackLog(event, "Sold", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSold is a log parse operation binding the contract event 0x6d8b7be78706a9b4a0413c29b82265e6de593cf330075475acbef8c50fed4fb2.
//
// Solidity: event Sold(uint256 indexed listingId, address indexed buyer, address indexed seller, uint256 price, uint256 royaltyAmount, uint256 marketFee, uint256 sellerProceeds)
func (_SudaMarketplace *SudaMarketplaceFilterer) ParseSold(log types.Log) (*SudaMarketplaceSold, error) {
	event := new(SudaMarketplaceSold)
	if err := _SudaMarketplace.contract.UnpackLog(event, "Sold", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaMarketplaceUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the SudaMarketplace contract.
type SudaMarketplaceUnpausedIterator struct {
	Event *SudaMarketplaceUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *SudaMarketplaceUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaMarketplaceUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(SudaMarketplaceUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *SudaMarketplaceUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaMarketplaceUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaMarketplaceUnpaused represents a Unpaused event raised by the SudaMarketplace contract.
type SudaMarketplaceUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SudaMarketplace *SudaMarketplaceFilterer) FilterUnpaused(opts *bind.FilterOpts) (*SudaMarketplaceUnpausedIterator, error) {

	logs, sub, err := _SudaMarketplace.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &SudaMarketplaceUnpausedIterator{contract: _SudaMarketplace.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SudaMarketplace *SudaMarketplaceFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *SudaMarketplaceUnpaused) (event.Subscription, error) {

	logs, sub, err := _SudaMarketplace.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaMarketplaceUnpaused)
				if err := _SudaMarketplace.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SudaMarketplace *SudaMarketplaceFilterer) ParseUnpaused(log types.Log) (*SudaMarketplaceUnpaused, error) {
	event := new(SudaMarketplaceUnpaused)
	if err := _SudaMarketplace.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
