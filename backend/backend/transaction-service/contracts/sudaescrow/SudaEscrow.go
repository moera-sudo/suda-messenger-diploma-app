// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sudaescrow

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

// SudaEscrowMetaData contains all meta data concerning the SudaEscrow contract.
var SudaEscrowMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_sudaToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"assignee\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"}],\"name\":\"QuestApproved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"}],\"name\":\"QuestCancelled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"assignee\",\"type\":\"address\"}],\"name\":\"QuestClaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"QuestCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"}],\"name\":\"QuestExpired\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"assignee\",\"type\":\"address\"}],\"name\":\"QuestSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"}],\"name\":\"cancel\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"createQuest\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"}],\"name\":\"expire\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextQuestId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"quests\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"assignee\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"enumSudaEscrow.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"questId\",\"type\":\"uint256\"}],\"name\":\"submit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sudaToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// SudaEscrowABI is the input ABI used to generate the binding from.
// Deprecated: Use SudaEscrowMetaData.ABI instead.
var SudaEscrowABI = SudaEscrowMetaData.ABI

// SudaEscrow is an auto generated Go binding around an Ethereum contract.
type SudaEscrow struct {
	SudaEscrowCaller     // Read-only binding to the contract
	SudaEscrowTransactor // Write-only binding to the contract
	SudaEscrowFilterer   // Log filterer for contract events
}

// SudaEscrowCaller is an auto generated read-only Go binding around an Ethereum contract.
type SudaEscrowCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SudaEscrowTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SudaEscrowTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SudaEscrowFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SudaEscrowFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SudaEscrowSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SudaEscrowSession struct {
	Contract     *SudaEscrow       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SudaEscrowCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SudaEscrowCallerSession struct {
	Contract *SudaEscrowCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// SudaEscrowTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SudaEscrowTransactorSession struct {
	Contract     *SudaEscrowTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// SudaEscrowRaw is an auto generated low-level Go binding around an Ethereum contract.
type SudaEscrowRaw struct {
	Contract *SudaEscrow // Generic contract binding to access the raw methods on
}

// SudaEscrowCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SudaEscrowCallerRaw struct {
	Contract *SudaEscrowCaller // Generic read-only contract binding to access the raw methods on
}

// SudaEscrowTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SudaEscrowTransactorRaw struct {
	Contract *SudaEscrowTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSudaEscrow creates a new instance of SudaEscrow, bound to a specific deployed contract.
func NewSudaEscrow(address common.Address, backend bind.ContractBackend) (*SudaEscrow, error) {
	contract, err := bindSudaEscrow(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SudaEscrow{SudaEscrowCaller: SudaEscrowCaller{contract: contract}, SudaEscrowTransactor: SudaEscrowTransactor{contract: contract}, SudaEscrowFilterer: SudaEscrowFilterer{contract: contract}}, nil
}

// NewSudaEscrowCaller creates a new read-only instance of SudaEscrow, bound to a specific deployed contract.
func NewSudaEscrowCaller(address common.Address, caller bind.ContractCaller) (*SudaEscrowCaller, error) {
	contract, err := bindSudaEscrow(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SudaEscrowCaller{contract: contract}, nil
}

// NewSudaEscrowTransactor creates a new write-only instance of SudaEscrow, bound to a specific deployed contract.
func NewSudaEscrowTransactor(address common.Address, transactor bind.ContractTransactor) (*SudaEscrowTransactor, error) {
	contract, err := bindSudaEscrow(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SudaEscrowTransactor{contract: contract}, nil
}

// NewSudaEscrowFilterer creates a new log filterer instance of SudaEscrow, bound to a specific deployed contract.
func NewSudaEscrowFilterer(address common.Address, filterer bind.ContractFilterer) (*SudaEscrowFilterer, error) {
	contract, err := bindSudaEscrow(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SudaEscrowFilterer{contract: contract}, nil
}

// bindSudaEscrow binds a generic wrapper to an already deployed contract.
func bindSudaEscrow(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SudaEscrowMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SudaEscrow *SudaEscrowRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SudaEscrow.Contract.SudaEscrowCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SudaEscrow *SudaEscrowRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaEscrow.Contract.SudaEscrowTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SudaEscrow *SudaEscrowRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SudaEscrow.Contract.SudaEscrowTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SudaEscrow *SudaEscrowCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SudaEscrow.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SudaEscrow *SudaEscrowTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaEscrow.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SudaEscrow *SudaEscrowTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SudaEscrow.Contract.contract.Transact(opts, method, params...)
}

// NextQuestId is a free data retrieval call binding the contract method 0xd638c983.
//
// Solidity: function nextQuestId() view returns(uint256)
func (_SudaEscrow *SudaEscrowCaller) NextQuestId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SudaEscrow.contract.Call(opts, &out, "nextQuestId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextQuestId is a free data retrieval call binding the contract method 0xd638c983.
//
// Solidity: function nextQuestId() view returns(uint256)
func (_SudaEscrow *SudaEscrowSession) NextQuestId() (*big.Int, error) {
	return _SudaEscrow.Contract.NextQuestId(&_SudaEscrow.CallOpts)
}

// NextQuestId is a free data retrieval call binding the contract method 0xd638c983.
//
// Solidity: function nextQuestId() view returns(uint256)
func (_SudaEscrow *SudaEscrowCallerSession) NextQuestId() (*big.Int, error) {
	return _SudaEscrow.Contract.NextQuestId(&_SudaEscrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SudaEscrow *SudaEscrowCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SudaEscrow.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SudaEscrow *SudaEscrowSession) Owner() (common.Address, error) {
	return _SudaEscrow.Contract.Owner(&_SudaEscrow.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SudaEscrow *SudaEscrowCallerSession) Owner() (common.Address, error) {
	return _SudaEscrow.Contract.Owner(&_SudaEscrow.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SudaEscrow *SudaEscrowCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SudaEscrow.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SudaEscrow *SudaEscrowSession) Paused() (bool, error) {
	return _SudaEscrow.Contract.Paused(&_SudaEscrow.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SudaEscrow *SudaEscrowCallerSession) Paused() (bool, error) {
	return _SudaEscrow.Contract.Paused(&_SudaEscrow.CallOpts)
}

// Quests is a free data retrieval call binding the contract method 0xe085f980.
//
// Solidity: function quests(uint256 ) view returns(address creator, address assignee, uint256 reward, uint256 deadline, uint8 status)
func (_SudaEscrow *SudaEscrowCaller) Quests(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Creator  common.Address
	Assignee common.Address
	Reward   *big.Int
	Deadline *big.Int
	Status   uint8
}, error) {
	var out []interface{}
	err := _SudaEscrow.contract.Call(opts, &out, "quests", arg0)

	outstruct := new(struct {
		Creator  common.Address
		Assignee common.Address
		Reward   *big.Int
		Deadline *big.Int
		Status   uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Creator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Assignee = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Reward = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Deadline = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Status = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

// Quests is a free data retrieval call binding the contract method 0xe085f980.
//
// Solidity: function quests(uint256 ) view returns(address creator, address assignee, uint256 reward, uint256 deadline, uint8 status)
func (_SudaEscrow *SudaEscrowSession) Quests(arg0 *big.Int) (struct {
	Creator  common.Address
	Assignee common.Address
	Reward   *big.Int
	Deadline *big.Int
	Status   uint8
}, error) {
	return _SudaEscrow.Contract.Quests(&_SudaEscrow.CallOpts, arg0)
}

// Quests is a free data retrieval call binding the contract method 0xe085f980.
//
// Solidity: function quests(uint256 ) view returns(address creator, address assignee, uint256 reward, uint256 deadline, uint8 status)
func (_SudaEscrow *SudaEscrowCallerSession) Quests(arg0 *big.Int) (struct {
	Creator  common.Address
	Assignee common.Address
	Reward   *big.Int
	Deadline *big.Int
	Status   uint8
}, error) {
	return _SudaEscrow.Contract.Quests(&_SudaEscrow.CallOpts, arg0)
}

// SudaToken is a free data retrieval call binding the contract method 0xda71113f.
//
// Solidity: function sudaToken() view returns(address)
func (_SudaEscrow *SudaEscrowCaller) SudaToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SudaEscrow.contract.Call(opts, &out, "sudaToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SudaToken is a free data retrieval call binding the contract method 0xda71113f.
//
// Solidity: function sudaToken() view returns(address)
func (_SudaEscrow *SudaEscrowSession) SudaToken() (common.Address, error) {
	return _SudaEscrow.Contract.SudaToken(&_SudaEscrow.CallOpts)
}

// SudaToken is a free data retrieval call binding the contract method 0xda71113f.
//
// Solidity: function sudaToken() view returns(address)
func (_SudaEscrow *SudaEscrowCallerSession) SudaToken() (common.Address, error) {
	return _SudaEscrow.Contract.SudaToken(&_SudaEscrow.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0xb759f954.
//
// Solidity: function approve(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowTransactor) Approve(opts *bind.TransactOpts, questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.contract.Transact(opts, "approve", questId)
}

// Approve is a paid mutator transaction binding the contract method 0xb759f954.
//
// Solidity: function approve(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowSession) Approve(questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.Approve(&_SudaEscrow.TransactOpts, questId)
}

// Approve is a paid mutator transaction binding the contract method 0xb759f954.
//
// Solidity: function approve(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowTransactorSession) Approve(questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.Approve(&_SudaEscrow.TransactOpts, questId)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowTransactor) Cancel(opts *bind.TransactOpts, questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.contract.Transact(opts, "cancel", questId)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowSession) Cancel(questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.Cancel(&_SudaEscrow.TransactOpts, questId)
}

// Cancel is a paid mutator transaction binding the contract method 0x40e58ee5.
//
// Solidity: function cancel(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowTransactorSession) Cancel(questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.Cancel(&_SudaEscrow.TransactOpts, questId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowTransactor) Claim(opts *bind.TransactOpts, questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.contract.Transact(opts, "claim", questId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowSession) Claim(questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.Claim(&_SudaEscrow.TransactOpts, questId)
}

// Claim is a paid mutator transaction binding the contract method 0x379607f5.
//
// Solidity: function claim(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowTransactorSession) Claim(questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.Claim(&_SudaEscrow.TransactOpts, questId)
}

// CreateQuest is a paid mutator transaction binding the contract method 0x0cfd7cb1.
//
// Solidity: function createQuest(uint256 reward, uint256 deadline) returns(uint256 questId)
func (_SudaEscrow *SudaEscrowTransactor) CreateQuest(opts *bind.TransactOpts, reward *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.contract.Transact(opts, "createQuest", reward, deadline)
}

// CreateQuest is a paid mutator transaction binding the contract method 0x0cfd7cb1.
//
// Solidity: function createQuest(uint256 reward, uint256 deadline) returns(uint256 questId)
func (_SudaEscrow *SudaEscrowSession) CreateQuest(reward *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.CreateQuest(&_SudaEscrow.TransactOpts, reward, deadline)
}

// CreateQuest is a paid mutator transaction binding the contract method 0x0cfd7cb1.
//
// Solidity: function createQuest(uint256 reward, uint256 deadline) returns(uint256 questId)
func (_SudaEscrow *SudaEscrowTransactorSession) CreateQuest(reward *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.CreateQuest(&_SudaEscrow.TransactOpts, reward, deadline)
}

// Expire is a paid mutator transaction binding the contract method 0xbf81bf43.
//
// Solidity: function expire(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowTransactor) Expire(opts *bind.TransactOpts, questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.contract.Transact(opts, "expire", questId)
}

// Expire is a paid mutator transaction binding the contract method 0xbf81bf43.
//
// Solidity: function expire(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowSession) Expire(questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.Expire(&_SudaEscrow.TransactOpts, questId)
}

// Expire is a paid mutator transaction binding the contract method 0xbf81bf43.
//
// Solidity: function expire(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowTransactorSession) Expire(questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.Expire(&_SudaEscrow.TransactOpts, questId)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SudaEscrow *SudaEscrowTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaEscrow.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SudaEscrow *SudaEscrowSession) Pause() (*types.Transaction, error) {
	return _SudaEscrow.Contract.Pause(&_SudaEscrow.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SudaEscrow *SudaEscrowTransactorSession) Pause() (*types.Transaction, error) {
	return _SudaEscrow.Contract.Pause(&_SudaEscrow.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SudaEscrow *SudaEscrowTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaEscrow.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SudaEscrow *SudaEscrowSession) RenounceOwnership() (*types.Transaction, error) {
	return _SudaEscrow.Contract.RenounceOwnership(&_SudaEscrow.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SudaEscrow *SudaEscrowTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SudaEscrow.Contract.RenounceOwnership(&_SudaEscrow.TransactOpts)
}

// Submit is a paid mutator transaction binding the contract method 0xea99c2a6.
//
// Solidity: function submit(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowTransactor) Submit(opts *bind.TransactOpts, questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.contract.Transact(opts, "submit", questId)
}

// Submit is a paid mutator transaction binding the contract method 0xea99c2a6.
//
// Solidity: function submit(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowSession) Submit(questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.Submit(&_SudaEscrow.TransactOpts, questId)
}

// Submit is a paid mutator transaction binding the contract method 0xea99c2a6.
//
// Solidity: function submit(uint256 questId) returns()
func (_SudaEscrow *SudaEscrowTransactorSession) Submit(questId *big.Int) (*types.Transaction, error) {
	return _SudaEscrow.Contract.Submit(&_SudaEscrow.TransactOpts, questId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SudaEscrow *SudaEscrowTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SudaEscrow.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SudaEscrow *SudaEscrowSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SudaEscrow.Contract.TransferOwnership(&_SudaEscrow.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SudaEscrow *SudaEscrowTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SudaEscrow.Contract.TransferOwnership(&_SudaEscrow.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SudaEscrow *SudaEscrowTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaEscrow.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SudaEscrow *SudaEscrowSession) Unpause() (*types.Transaction, error) {
	return _SudaEscrow.Contract.Unpause(&_SudaEscrow.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SudaEscrow *SudaEscrowTransactorSession) Unpause() (*types.Transaction, error) {
	return _SudaEscrow.Contract.Unpause(&_SudaEscrow.TransactOpts)
}

// SudaEscrowOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SudaEscrow contract.
type SudaEscrowOwnershipTransferredIterator struct {
	Event *SudaEscrowOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SudaEscrowOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaEscrowOwnershipTransferred)
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
		it.Event = new(SudaEscrowOwnershipTransferred)
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
func (it *SudaEscrowOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaEscrowOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaEscrowOwnershipTransferred represents a OwnershipTransferred event raised by the SudaEscrow contract.
type SudaEscrowOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SudaEscrow *SudaEscrowFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SudaEscrowOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SudaEscrow.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SudaEscrowOwnershipTransferredIterator{contract: _SudaEscrow.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SudaEscrow *SudaEscrowFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SudaEscrowOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SudaEscrow.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaEscrowOwnershipTransferred)
				if err := _SudaEscrow.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SudaEscrow *SudaEscrowFilterer) ParseOwnershipTransferred(log types.Log) (*SudaEscrowOwnershipTransferred, error) {
	event := new(SudaEscrowOwnershipTransferred)
	if err := _SudaEscrow.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaEscrowPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the SudaEscrow contract.
type SudaEscrowPausedIterator struct {
	Event *SudaEscrowPaused // Event containing the contract specifics and raw log

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
func (it *SudaEscrowPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaEscrowPaused)
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
		it.Event = new(SudaEscrowPaused)
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
func (it *SudaEscrowPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaEscrowPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaEscrowPaused represents a Paused event raised by the SudaEscrow contract.
type SudaEscrowPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SudaEscrow *SudaEscrowFilterer) FilterPaused(opts *bind.FilterOpts) (*SudaEscrowPausedIterator, error) {

	logs, sub, err := _SudaEscrow.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &SudaEscrowPausedIterator{contract: _SudaEscrow.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SudaEscrow *SudaEscrowFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *SudaEscrowPaused) (event.Subscription, error) {

	logs, sub, err := _SudaEscrow.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaEscrowPaused)
				if err := _SudaEscrow.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_SudaEscrow *SudaEscrowFilterer) ParsePaused(log types.Log) (*SudaEscrowPaused, error) {
	event := new(SudaEscrowPaused)
	if err := _SudaEscrow.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaEscrowQuestApprovedIterator is returned from FilterQuestApproved and is used to iterate over the raw logs and unpacked data for QuestApproved events raised by the SudaEscrow contract.
type SudaEscrowQuestApprovedIterator struct {
	Event *SudaEscrowQuestApproved // Event containing the contract specifics and raw log

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
func (it *SudaEscrowQuestApprovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaEscrowQuestApproved)
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
		it.Event = new(SudaEscrowQuestApproved)
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
func (it *SudaEscrowQuestApprovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaEscrowQuestApprovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaEscrowQuestApproved represents a QuestApproved event raised by the SudaEscrow contract.
type SudaEscrowQuestApproved struct {
	QuestId  *big.Int
	Assignee common.Address
	Reward   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterQuestApproved is a free log retrieval operation binding the contract event 0x7c65ec455c543628d7fc06032bde9028ef702625279f55024de835a4c4361655.
//
// Solidity: event QuestApproved(uint256 indexed questId, address indexed assignee, uint256 reward)
func (_SudaEscrow *SudaEscrowFilterer) FilterQuestApproved(opts *bind.FilterOpts, questId []*big.Int, assignee []common.Address) (*SudaEscrowQuestApprovedIterator, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}
	var assigneeRule []interface{}
	for _, assigneeItem := range assignee {
		assigneeRule = append(assigneeRule, assigneeItem)
	}

	logs, sub, err := _SudaEscrow.contract.FilterLogs(opts, "QuestApproved", questIdRule, assigneeRule)
	if err != nil {
		return nil, err
	}
	return &SudaEscrowQuestApprovedIterator{contract: _SudaEscrow.contract, event: "QuestApproved", logs: logs, sub: sub}, nil
}

// WatchQuestApproved is a free log subscription operation binding the contract event 0x7c65ec455c543628d7fc06032bde9028ef702625279f55024de835a4c4361655.
//
// Solidity: event QuestApproved(uint256 indexed questId, address indexed assignee, uint256 reward)
func (_SudaEscrow *SudaEscrowFilterer) WatchQuestApproved(opts *bind.WatchOpts, sink chan<- *SudaEscrowQuestApproved, questId []*big.Int, assignee []common.Address) (event.Subscription, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}
	var assigneeRule []interface{}
	for _, assigneeItem := range assignee {
		assigneeRule = append(assigneeRule, assigneeItem)
	}

	logs, sub, err := _SudaEscrow.contract.WatchLogs(opts, "QuestApproved", questIdRule, assigneeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaEscrowQuestApproved)
				if err := _SudaEscrow.contract.UnpackLog(event, "QuestApproved", log); err != nil {
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

// ParseQuestApproved is a log parse operation binding the contract event 0x7c65ec455c543628d7fc06032bde9028ef702625279f55024de835a4c4361655.
//
// Solidity: event QuestApproved(uint256 indexed questId, address indexed assignee, uint256 reward)
func (_SudaEscrow *SudaEscrowFilterer) ParseQuestApproved(log types.Log) (*SudaEscrowQuestApproved, error) {
	event := new(SudaEscrowQuestApproved)
	if err := _SudaEscrow.contract.UnpackLog(event, "QuestApproved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaEscrowQuestCancelledIterator is returned from FilterQuestCancelled and is used to iterate over the raw logs and unpacked data for QuestCancelled events raised by the SudaEscrow contract.
type SudaEscrowQuestCancelledIterator struct {
	Event *SudaEscrowQuestCancelled // Event containing the contract specifics and raw log

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
func (it *SudaEscrowQuestCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaEscrowQuestCancelled)
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
		it.Event = new(SudaEscrowQuestCancelled)
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
func (it *SudaEscrowQuestCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaEscrowQuestCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaEscrowQuestCancelled represents a QuestCancelled event raised by the SudaEscrow contract.
type SudaEscrowQuestCancelled struct {
	QuestId *big.Int
	Creator common.Address
	Refund  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterQuestCancelled is a free log retrieval operation binding the contract event 0xd8ea0369deba99bb1177d79b8ced5ba0f6b4c36ee869ff2394b730233a21ffed.
//
// Solidity: event QuestCancelled(uint256 indexed questId, address indexed creator, uint256 refund)
func (_SudaEscrow *SudaEscrowFilterer) FilterQuestCancelled(opts *bind.FilterOpts, questId []*big.Int, creator []common.Address) (*SudaEscrowQuestCancelledIterator, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SudaEscrow.contract.FilterLogs(opts, "QuestCancelled", questIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &SudaEscrowQuestCancelledIterator{contract: _SudaEscrow.contract, event: "QuestCancelled", logs: logs, sub: sub}, nil
}

// WatchQuestCancelled is a free log subscription operation binding the contract event 0xd8ea0369deba99bb1177d79b8ced5ba0f6b4c36ee869ff2394b730233a21ffed.
//
// Solidity: event QuestCancelled(uint256 indexed questId, address indexed creator, uint256 refund)
func (_SudaEscrow *SudaEscrowFilterer) WatchQuestCancelled(opts *bind.WatchOpts, sink chan<- *SudaEscrowQuestCancelled, questId []*big.Int, creator []common.Address) (event.Subscription, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SudaEscrow.contract.WatchLogs(opts, "QuestCancelled", questIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaEscrowQuestCancelled)
				if err := _SudaEscrow.contract.UnpackLog(event, "QuestCancelled", log); err != nil {
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

// ParseQuestCancelled is a log parse operation binding the contract event 0xd8ea0369deba99bb1177d79b8ced5ba0f6b4c36ee869ff2394b730233a21ffed.
//
// Solidity: event QuestCancelled(uint256 indexed questId, address indexed creator, uint256 refund)
func (_SudaEscrow *SudaEscrowFilterer) ParseQuestCancelled(log types.Log) (*SudaEscrowQuestCancelled, error) {
	event := new(SudaEscrowQuestCancelled)
	if err := _SudaEscrow.contract.UnpackLog(event, "QuestCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaEscrowQuestClaimedIterator is returned from FilterQuestClaimed and is used to iterate over the raw logs and unpacked data for QuestClaimed events raised by the SudaEscrow contract.
type SudaEscrowQuestClaimedIterator struct {
	Event *SudaEscrowQuestClaimed // Event containing the contract specifics and raw log

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
func (it *SudaEscrowQuestClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaEscrowQuestClaimed)
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
		it.Event = new(SudaEscrowQuestClaimed)
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
func (it *SudaEscrowQuestClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaEscrowQuestClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaEscrowQuestClaimed represents a QuestClaimed event raised by the SudaEscrow contract.
type SudaEscrowQuestClaimed struct {
	QuestId  *big.Int
	Assignee common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterQuestClaimed is a free log retrieval operation binding the contract event 0x65fbb88a4254be2c0b66c696f552e3a8e933c782fed1b8e81a103b3e5f66b59c.
//
// Solidity: event QuestClaimed(uint256 indexed questId, address indexed assignee)
func (_SudaEscrow *SudaEscrowFilterer) FilterQuestClaimed(opts *bind.FilterOpts, questId []*big.Int, assignee []common.Address) (*SudaEscrowQuestClaimedIterator, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}
	var assigneeRule []interface{}
	for _, assigneeItem := range assignee {
		assigneeRule = append(assigneeRule, assigneeItem)
	}

	logs, sub, err := _SudaEscrow.contract.FilterLogs(opts, "QuestClaimed", questIdRule, assigneeRule)
	if err != nil {
		return nil, err
	}
	return &SudaEscrowQuestClaimedIterator{contract: _SudaEscrow.contract, event: "QuestClaimed", logs: logs, sub: sub}, nil
}

// WatchQuestClaimed is a free log subscription operation binding the contract event 0x65fbb88a4254be2c0b66c696f552e3a8e933c782fed1b8e81a103b3e5f66b59c.
//
// Solidity: event QuestClaimed(uint256 indexed questId, address indexed assignee)
func (_SudaEscrow *SudaEscrowFilterer) WatchQuestClaimed(opts *bind.WatchOpts, sink chan<- *SudaEscrowQuestClaimed, questId []*big.Int, assignee []common.Address) (event.Subscription, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}
	var assigneeRule []interface{}
	for _, assigneeItem := range assignee {
		assigneeRule = append(assigneeRule, assigneeItem)
	}

	logs, sub, err := _SudaEscrow.contract.WatchLogs(opts, "QuestClaimed", questIdRule, assigneeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaEscrowQuestClaimed)
				if err := _SudaEscrow.contract.UnpackLog(event, "QuestClaimed", log); err != nil {
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

// ParseQuestClaimed is a log parse operation binding the contract event 0x65fbb88a4254be2c0b66c696f552e3a8e933c782fed1b8e81a103b3e5f66b59c.
//
// Solidity: event QuestClaimed(uint256 indexed questId, address indexed assignee)
func (_SudaEscrow *SudaEscrowFilterer) ParseQuestClaimed(log types.Log) (*SudaEscrowQuestClaimed, error) {
	event := new(SudaEscrowQuestClaimed)
	if err := _SudaEscrow.contract.UnpackLog(event, "QuestClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaEscrowQuestCreatedIterator is returned from FilterQuestCreated and is used to iterate over the raw logs and unpacked data for QuestCreated events raised by the SudaEscrow contract.
type SudaEscrowQuestCreatedIterator struct {
	Event *SudaEscrowQuestCreated // Event containing the contract specifics and raw log

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
func (it *SudaEscrowQuestCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaEscrowQuestCreated)
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
		it.Event = new(SudaEscrowQuestCreated)
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
func (it *SudaEscrowQuestCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaEscrowQuestCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaEscrowQuestCreated represents a QuestCreated event raised by the SudaEscrow contract.
type SudaEscrowQuestCreated struct {
	QuestId  *big.Int
	Creator  common.Address
	Reward   *big.Int
	Deadline *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterQuestCreated is a free log retrieval operation binding the contract event 0x0867a04772979958deca46ac6224b77c7363fb1f765547324c56a23c1225b762.
//
// Solidity: event QuestCreated(uint256 indexed questId, address indexed creator, uint256 reward, uint256 deadline)
func (_SudaEscrow *SudaEscrowFilterer) FilterQuestCreated(opts *bind.FilterOpts, questId []*big.Int, creator []common.Address) (*SudaEscrowQuestCreatedIterator, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SudaEscrow.contract.FilterLogs(opts, "QuestCreated", questIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &SudaEscrowQuestCreatedIterator{contract: _SudaEscrow.contract, event: "QuestCreated", logs: logs, sub: sub}, nil
}

// WatchQuestCreated is a free log subscription operation binding the contract event 0x0867a04772979958deca46ac6224b77c7363fb1f765547324c56a23c1225b762.
//
// Solidity: event QuestCreated(uint256 indexed questId, address indexed creator, uint256 reward, uint256 deadline)
func (_SudaEscrow *SudaEscrowFilterer) WatchQuestCreated(opts *bind.WatchOpts, sink chan<- *SudaEscrowQuestCreated, questId []*big.Int, creator []common.Address) (event.Subscription, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SudaEscrow.contract.WatchLogs(opts, "QuestCreated", questIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaEscrowQuestCreated)
				if err := _SudaEscrow.contract.UnpackLog(event, "QuestCreated", log); err != nil {
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

// ParseQuestCreated is a log parse operation binding the contract event 0x0867a04772979958deca46ac6224b77c7363fb1f765547324c56a23c1225b762.
//
// Solidity: event QuestCreated(uint256 indexed questId, address indexed creator, uint256 reward, uint256 deadline)
func (_SudaEscrow *SudaEscrowFilterer) ParseQuestCreated(log types.Log) (*SudaEscrowQuestCreated, error) {
	event := new(SudaEscrowQuestCreated)
	if err := _SudaEscrow.contract.UnpackLog(event, "QuestCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaEscrowQuestExpiredIterator is returned from FilterQuestExpired and is used to iterate over the raw logs and unpacked data for QuestExpired events raised by the SudaEscrow contract.
type SudaEscrowQuestExpiredIterator struct {
	Event *SudaEscrowQuestExpired // Event containing the contract specifics and raw log

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
func (it *SudaEscrowQuestExpiredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaEscrowQuestExpired)
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
		it.Event = new(SudaEscrowQuestExpired)
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
func (it *SudaEscrowQuestExpiredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaEscrowQuestExpiredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaEscrowQuestExpired represents a QuestExpired event raised by the SudaEscrow contract.
type SudaEscrowQuestExpired struct {
	QuestId *big.Int
	Refund  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterQuestExpired is a free log retrieval operation binding the contract event 0x9efc8029c25ebf0053237270c247a73d6c08d42671be0ce5991cffa2a3fd8eef.
//
// Solidity: event QuestExpired(uint256 indexed questId, uint256 refund)
func (_SudaEscrow *SudaEscrowFilterer) FilterQuestExpired(opts *bind.FilterOpts, questId []*big.Int) (*SudaEscrowQuestExpiredIterator, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}

	logs, sub, err := _SudaEscrow.contract.FilterLogs(opts, "QuestExpired", questIdRule)
	if err != nil {
		return nil, err
	}
	return &SudaEscrowQuestExpiredIterator{contract: _SudaEscrow.contract, event: "QuestExpired", logs: logs, sub: sub}, nil
}

// WatchQuestExpired is a free log subscription operation binding the contract event 0x9efc8029c25ebf0053237270c247a73d6c08d42671be0ce5991cffa2a3fd8eef.
//
// Solidity: event QuestExpired(uint256 indexed questId, uint256 refund)
func (_SudaEscrow *SudaEscrowFilterer) WatchQuestExpired(opts *bind.WatchOpts, sink chan<- *SudaEscrowQuestExpired, questId []*big.Int) (event.Subscription, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}

	logs, sub, err := _SudaEscrow.contract.WatchLogs(opts, "QuestExpired", questIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaEscrowQuestExpired)
				if err := _SudaEscrow.contract.UnpackLog(event, "QuestExpired", log); err != nil {
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

// ParseQuestExpired is a log parse operation binding the contract event 0x9efc8029c25ebf0053237270c247a73d6c08d42671be0ce5991cffa2a3fd8eef.
//
// Solidity: event QuestExpired(uint256 indexed questId, uint256 refund)
func (_SudaEscrow *SudaEscrowFilterer) ParseQuestExpired(log types.Log) (*SudaEscrowQuestExpired, error) {
	event := new(SudaEscrowQuestExpired)
	if err := _SudaEscrow.contract.UnpackLog(event, "QuestExpired", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaEscrowQuestSubmittedIterator is returned from FilterQuestSubmitted and is used to iterate over the raw logs and unpacked data for QuestSubmitted events raised by the SudaEscrow contract.
type SudaEscrowQuestSubmittedIterator struct {
	Event *SudaEscrowQuestSubmitted // Event containing the contract specifics and raw log

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
func (it *SudaEscrowQuestSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaEscrowQuestSubmitted)
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
		it.Event = new(SudaEscrowQuestSubmitted)
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
func (it *SudaEscrowQuestSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaEscrowQuestSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaEscrowQuestSubmitted represents a QuestSubmitted event raised by the SudaEscrow contract.
type SudaEscrowQuestSubmitted struct {
	QuestId  *big.Int
	Assignee common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterQuestSubmitted is a free log retrieval operation binding the contract event 0x1089e840f3a5a41e7afdd0f2d871774341e15367cb89284fe68b2229ebb042cf.
//
// Solidity: event QuestSubmitted(uint256 indexed questId, address indexed assignee)
func (_SudaEscrow *SudaEscrowFilterer) FilterQuestSubmitted(opts *bind.FilterOpts, questId []*big.Int, assignee []common.Address) (*SudaEscrowQuestSubmittedIterator, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}
	var assigneeRule []interface{}
	for _, assigneeItem := range assignee {
		assigneeRule = append(assigneeRule, assigneeItem)
	}

	logs, sub, err := _SudaEscrow.contract.FilterLogs(opts, "QuestSubmitted", questIdRule, assigneeRule)
	if err != nil {
		return nil, err
	}
	return &SudaEscrowQuestSubmittedIterator{contract: _SudaEscrow.contract, event: "QuestSubmitted", logs: logs, sub: sub}, nil
}

// WatchQuestSubmitted is a free log subscription operation binding the contract event 0x1089e840f3a5a41e7afdd0f2d871774341e15367cb89284fe68b2229ebb042cf.
//
// Solidity: event QuestSubmitted(uint256 indexed questId, address indexed assignee)
func (_SudaEscrow *SudaEscrowFilterer) WatchQuestSubmitted(opts *bind.WatchOpts, sink chan<- *SudaEscrowQuestSubmitted, questId []*big.Int, assignee []common.Address) (event.Subscription, error) {

	var questIdRule []interface{}
	for _, questIdItem := range questId {
		questIdRule = append(questIdRule, questIdItem)
	}
	var assigneeRule []interface{}
	for _, assigneeItem := range assignee {
		assigneeRule = append(assigneeRule, assigneeItem)
	}

	logs, sub, err := _SudaEscrow.contract.WatchLogs(opts, "QuestSubmitted", questIdRule, assigneeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaEscrowQuestSubmitted)
				if err := _SudaEscrow.contract.UnpackLog(event, "QuestSubmitted", log); err != nil {
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

// ParseQuestSubmitted is a log parse operation binding the contract event 0x1089e840f3a5a41e7afdd0f2d871774341e15367cb89284fe68b2229ebb042cf.
//
// Solidity: event QuestSubmitted(uint256 indexed questId, address indexed assignee)
func (_SudaEscrow *SudaEscrowFilterer) ParseQuestSubmitted(log types.Log) (*SudaEscrowQuestSubmitted, error) {
	event := new(SudaEscrowQuestSubmitted)
	if err := _SudaEscrow.contract.UnpackLog(event, "QuestSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaEscrowUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the SudaEscrow contract.
type SudaEscrowUnpausedIterator struct {
	Event *SudaEscrowUnpaused // Event containing the contract specifics and raw log

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
func (it *SudaEscrowUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaEscrowUnpaused)
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
		it.Event = new(SudaEscrowUnpaused)
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
func (it *SudaEscrowUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaEscrowUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaEscrowUnpaused represents a Unpaused event raised by the SudaEscrow contract.
type SudaEscrowUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SudaEscrow *SudaEscrowFilterer) FilterUnpaused(opts *bind.FilterOpts) (*SudaEscrowUnpausedIterator, error) {

	logs, sub, err := _SudaEscrow.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &SudaEscrowUnpausedIterator{contract: _SudaEscrow.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SudaEscrow *SudaEscrowFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *SudaEscrowUnpaused) (event.Subscription, error) {

	logs, sub, err := _SudaEscrow.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaEscrowUnpaused)
				if err := _SudaEscrow.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_SudaEscrow *SudaEscrowFilterer) ParseUnpaused(log types.Log) (*SudaEscrowUnpaused, error) {
	event := new(SudaEscrowUnpaused)
	if err := _SudaEscrow.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
