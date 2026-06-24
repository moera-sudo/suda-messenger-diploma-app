// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sudafundraising

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

// SudaFundraisingMetaData contains all meta data concerning the SudaFundraising contract.
var SudaFundraisingMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"},{\"internalType\":\"contractIERC20\",\"name\":\"_sudaToken\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"EnforcedPause\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExpectedPause\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fundraiserId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"donor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalRaised\",\"type\":\"uint256\"}],\"name\":\"Donated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fundraiserId\",\"type\":\"uint256\"}],\"name\":\"Expired\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fundraiserId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"goal\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"FundraiserCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fundraiserId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalRaised\",\"type\":\"uint256\"}],\"name\":\"GoalReached\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fundraiserId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"donor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Refunded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fundraiserId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"contributions\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"goal\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"create\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"fundraiserId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fundraiserId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"donate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"fundraisers\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"creator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"goal\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"raised\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"enumSudaFundraising.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"nextFundraiserId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fundraiserId\",\"type\":\"uint256\"}],\"name\":\"refund\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"sudaToken\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"fundraiserId\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// SudaFundraisingABI is the input ABI used to generate the binding from.
// Deprecated: Use SudaFundraisingMetaData.ABI instead.
var SudaFundraisingABI = SudaFundraisingMetaData.ABI

// SudaFundraising is an auto generated Go binding around an Ethereum contract.
type SudaFundraising struct {
	SudaFundraisingCaller     // Read-only binding to the contract
	SudaFundraisingTransactor // Write-only binding to the contract
	SudaFundraisingFilterer   // Log filterer for contract events
}

// SudaFundraisingCaller is an auto generated read-only Go binding around an Ethereum contract.
type SudaFundraisingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SudaFundraisingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SudaFundraisingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SudaFundraisingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SudaFundraisingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SudaFundraisingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SudaFundraisingSession struct {
	Contract     *SudaFundraising  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SudaFundraisingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SudaFundraisingCallerSession struct {
	Contract *SudaFundraisingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// SudaFundraisingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SudaFundraisingTransactorSession struct {
	Contract     *SudaFundraisingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// SudaFundraisingRaw is an auto generated low-level Go binding around an Ethereum contract.
type SudaFundraisingRaw struct {
	Contract *SudaFundraising // Generic contract binding to access the raw methods on
}

// SudaFundraisingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SudaFundraisingCallerRaw struct {
	Contract *SudaFundraisingCaller // Generic read-only contract binding to access the raw methods on
}

// SudaFundraisingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SudaFundraisingTransactorRaw struct {
	Contract *SudaFundraisingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSudaFundraising creates a new instance of SudaFundraising, bound to a specific deployed contract.
func NewSudaFundraising(address common.Address, backend bind.ContractBackend) (*SudaFundraising, error) {
	contract, err := bindSudaFundraising(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SudaFundraising{SudaFundraisingCaller: SudaFundraisingCaller{contract: contract}, SudaFundraisingTransactor: SudaFundraisingTransactor{contract: contract}, SudaFundraisingFilterer: SudaFundraisingFilterer{contract: contract}}, nil
}

// NewSudaFundraisingCaller creates a new read-only instance of SudaFundraising, bound to a specific deployed contract.
func NewSudaFundraisingCaller(address common.Address, caller bind.ContractCaller) (*SudaFundraisingCaller, error) {
	contract, err := bindSudaFundraising(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingCaller{contract: contract}, nil
}

// NewSudaFundraisingTransactor creates a new write-only instance of SudaFundraising, bound to a specific deployed contract.
func NewSudaFundraisingTransactor(address common.Address, transactor bind.ContractTransactor) (*SudaFundraisingTransactor, error) {
	contract, err := bindSudaFundraising(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingTransactor{contract: contract}, nil
}

// NewSudaFundraisingFilterer creates a new log filterer instance of SudaFundraising, bound to a specific deployed contract.
func NewSudaFundraisingFilterer(address common.Address, filterer bind.ContractFilterer) (*SudaFundraisingFilterer, error) {
	contract, err := bindSudaFundraising(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingFilterer{contract: contract}, nil
}

// bindSudaFundraising binds a generic wrapper to an already deployed contract.
func bindSudaFundraising(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SudaFundraisingMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SudaFundraising *SudaFundraisingRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SudaFundraising.Contract.SudaFundraisingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SudaFundraising *SudaFundraisingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaFundraising.Contract.SudaFundraisingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SudaFundraising *SudaFundraisingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SudaFundraising.Contract.SudaFundraisingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SudaFundraising *SudaFundraisingCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SudaFundraising.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SudaFundraising *SudaFundraisingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaFundraising.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SudaFundraising *SudaFundraisingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SudaFundraising.Contract.contract.Transact(opts, method, params...)
}

// Contributions is a free data retrieval call binding the contract method 0x3d891f59.
//
// Solidity: function contributions(uint256 , address ) view returns(uint256)
func (_SudaFundraising *SudaFundraisingCaller) Contributions(opts *bind.CallOpts, arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SudaFundraising.contract.Call(opts, &out, "contributions", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Contributions is a free data retrieval call binding the contract method 0x3d891f59.
//
// Solidity: function contributions(uint256 , address ) view returns(uint256)
func (_SudaFundraising *SudaFundraisingSession) Contributions(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _SudaFundraising.Contract.Contributions(&_SudaFundraising.CallOpts, arg0, arg1)
}

// Contributions is a free data retrieval call binding the contract method 0x3d891f59.
//
// Solidity: function contributions(uint256 , address ) view returns(uint256)
func (_SudaFundraising *SudaFundraisingCallerSession) Contributions(arg0 *big.Int, arg1 common.Address) (*big.Int, error) {
	return _SudaFundraising.Contract.Contributions(&_SudaFundraising.CallOpts, arg0, arg1)
}

// Fundraisers is a free data retrieval call binding the contract method 0xf38f72fb.
//
// Solidity: function fundraisers(uint256 ) view returns(address creator, uint256 goal, uint256 raised, uint256 deadline, uint8 status)
func (_SudaFundraising *SudaFundraisingCaller) Fundraisers(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Creator  common.Address
	Goal     *big.Int
	Raised   *big.Int
	Deadline *big.Int
	Status   uint8
}, error) {
	var out []interface{}
	err := _SudaFundraising.contract.Call(opts, &out, "fundraisers", arg0)

	outstruct := new(struct {
		Creator  common.Address
		Goal     *big.Int
		Raised   *big.Int
		Deadline *big.Int
		Status   uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Creator = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Goal = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Raised = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Deadline = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.Status = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

// Fundraisers is a free data retrieval call binding the contract method 0xf38f72fb.
//
// Solidity: function fundraisers(uint256 ) view returns(address creator, uint256 goal, uint256 raised, uint256 deadline, uint8 status)
func (_SudaFundraising *SudaFundraisingSession) Fundraisers(arg0 *big.Int) (struct {
	Creator  common.Address
	Goal     *big.Int
	Raised   *big.Int
	Deadline *big.Int
	Status   uint8
}, error) {
	return _SudaFundraising.Contract.Fundraisers(&_SudaFundraising.CallOpts, arg0)
}

// Fundraisers is a free data retrieval call binding the contract method 0xf38f72fb.
//
// Solidity: function fundraisers(uint256 ) view returns(address creator, uint256 goal, uint256 raised, uint256 deadline, uint8 status)
func (_SudaFundraising *SudaFundraisingCallerSession) Fundraisers(arg0 *big.Int) (struct {
	Creator  common.Address
	Goal     *big.Int
	Raised   *big.Int
	Deadline *big.Int
	Status   uint8
}, error) {
	return _SudaFundraising.Contract.Fundraisers(&_SudaFundraising.CallOpts, arg0)
}

// NextFundraiserId is a free data retrieval call binding the contract method 0x444480da.
//
// Solidity: function nextFundraiserId() view returns(uint256)
func (_SudaFundraising *SudaFundraisingCaller) NextFundraiserId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SudaFundraising.contract.Call(opts, &out, "nextFundraiserId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NextFundraiserId is a free data retrieval call binding the contract method 0x444480da.
//
// Solidity: function nextFundraiserId() view returns(uint256)
func (_SudaFundraising *SudaFundraisingSession) NextFundraiserId() (*big.Int, error) {
	return _SudaFundraising.Contract.NextFundraiserId(&_SudaFundraising.CallOpts)
}

// NextFundraiserId is a free data retrieval call binding the contract method 0x444480da.
//
// Solidity: function nextFundraiserId() view returns(uint256)
func (_SudaFundraising *SudaFundraisingCallerSession) NextFundraiserId() (*big.Int, error) {
	return _SudaFundraising.Contract.NextFundraiserId(&_SudaFundraising.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SudaFundraising *SudaFundraisingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SudaFundraising.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SudaFundraising *SudaFundraisingSession) Owner() (common.Address, error) {
	return _SudaFundraising.Contract.Owner(&_SudaFundraising.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SudaFundraising *SudaFundraisingCallerSession) Owner() (common.Address, error) {
	return _SudaFundraising.Contract.Owner(&_SudaFundraising.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SudaFundraising *SudaFundraisingCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SudaFundraising.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SudaFundraising *SudaFundraisingSession) Paused() (bool, error) {
	return _SudaFundraising.Contract.Paused(&_SudaFundraising.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_SudaFundraising *SudaFundraisingCallerSession) Paused() (bool, error) {
	return _SudaFundraising.Contract.Paused(&_SudaFundraising.CallOpts)
}

// SudaToken is a free data retrieval call binding the contract method 0xda71113f.
//
// Solidity: function sudaToken() view returns(address)
func (_SudaFundraising *SudaFundraisingCaller) SudaToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SudaFundraising.contract.Call(opts, &out, "sudaToken")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SudaToken is a free data retrieval call binding the contract method 0xda71113f.
//
// Solidity: function sudaToken() view returns(address)
func (_SudaFundraising *SudaFundraisingSession) SudaToken() (common.Address, error) {
	return _SudaFundraising.Contract.SudaToken(&_SudaFundraising.CallOpts)
}

// SudaToken is a free data retrieval call binding the contract method 0xda71113f.
//
// Solidity: function sudaToken() view returns(address)
func (_SudaFundraising *SudaFundraisingCallerSession) SudaToken() (common.Address, error) {
	return _SudaFundraising.Contract.SudaToken(&_SudaFundraising.CallOpts)
}

// Create is a paid mutator transaction binding the contract method 0x9f7b4579.
//
// Solidity: function create(uint256 goal, uint256 deadline) returns(uint256 fundraiserId)
func (_SudaFundraising *SudaFundraisingTransactor) Create(opts *bind.TransactOpts, goal *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.contract.Transact(opts, "create", goal, deadline)
}

// Create is a paid mutator transaction binding the contract method 0x9f7b4579.
//
// Solidity: function create(uint256 goal, uint256 deadline) returns(uint256 fundraiserId)
func (_SudaFundraising *SudaFundraisingSession) Create(goal *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.Contract.Create(&_SudaFundraising.TransactOpts, goal, deadline)
}

// Create is a paid mutator transaction binding the contract method 0x9f7b4579.
//
// Solidity: function create(uint256 goal, uint256 deadline) returns(uint256 fundraiserId)
func (_SudaFundraising *SudaFundraisingTransactorSession) Create(goal *big.Int, deadline *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.Contract.Create(&_SudaFundraising.TransactOpts, goal, deadline)
}

// Donate is a paid mutator transaction binding the contract method 0x0cdd53f6.
//
// Solidity: function donate(uint256 fundraiserId, uint256 amount) returns()
func (_SudaFundraising *SudaFundraisingTransactor) Donate(opts *bind.TransactOpts, fundraiserId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.contract.Transact(opts, "donate", fundraiserId, amount)
}

// Donate is a paid mutator transaction binding the contract method 0x0cdd53f6.
//
// Solidity: function donate(uint256 fundraiserId, uint256 amount) returns()
func (_SudaFundraising *SudaFundraisingSession) Donate(fundraiserId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.Contract.Donate(&_SudaFundraising.TransactOpts, fundraiserId, amount)
}

// Donate is a paid mutator transaction binding the contract method 0x0cdd53f6.
//
// Solidity: function donate(uint256 fundraiserId, uint256 amount) returns()
func (_SudaFundraising *SudaFundraisingTransactorSession) Donate(fundraiserId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.Contract.Donate(&_SudaFundraising.TransactOpts, fundraiserId, amount)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SudaFundraising *SudaFundraisingTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaFundraising.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SudaFundraising *SudaFundraisingSession) Pause() (*types.Transaction, error) {
	return _SudaFundraising.Contract.Pause(&_SudaFundraising.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_SudaFundraising *SudaFundraisingTransactorSession) Pause() (*types.Transaction, error) {
	return _SudaFundraising.Contract.Pause(&_SudaFundraising.TransactOpts)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 fundraiserId) returns()
func (_SudaFundraising *SudaFundraisingTransactor) Refund(opts *bind.TransactOpts, fundraiserId *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.contract.Transact(opts, "refund", fundraiserId)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 fundraiserId) returns()
func (_SudaFundraising *SudaFundraisingSession) Refund(fundraiserId *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.Contract.Refund(&_SudaFundraising.TransactOpts, fundraiserId)
}

// Refund is a paid mutator transaction binding the contract method 0x278ecde1.
//
// Solidity: function refund(uint256 fundraiserId) returns()
func (_SudaFundraising *SudaFundraisingTransactorSession) Refund(fundraiserId *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.Contract.Refund(&_SudaFundraising.TransactOpts, fundraiserId)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SudaFundraising *SudaFundraisingTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaFundraising.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SudaFundraising *SudaFundraisingSession) RenounceOwnership() (*types.Transaction, error) {
	return _SudaFundraising.Contract.RenounceOwnership(&_SudaFundraising.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SudaFundraising *SudaFundraisingTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SudaFundraising.Contract.RenounceOwnership(&_SudaFundraising.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SudaFundraising *SudaFundraisingTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SudaFundraising.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SudaFundraising *SudaFundraisingSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SudaFundraising.Contract.TransferOwnership(&_SudaFundraising.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SudaFundraising *SudaFundraisingTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SudaFundraising.Contract.TransferOwnership(&_SudaFundraising.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SudaFundraising *SudaFundraisingTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SudaFundraising.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SudaFundraising *SudaFundraisingSession) Unpause() (*types.Transaction, error) {
	return _SudaFundraising.Contract.Unpause(&_SudaFundraising.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_SudaFundraising *SudaFundraisingTransactorSession) Unpause() (*types.Transaction, error) {
	return _SudaFundraising.Contract.Unpause(&_SudaFundraising.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 fundraiserId) returns()
func (_SudaFundraising *SudaFundraisingTransactor) Withdraw(opts *bind.TransactOpts, fundraiserId *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.contract.Transact(opts, "withdraw", fundraiserId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 fundraiserId) returns()
func (_SudaFundraising *SudaFundraisingSession) Withdraw(fundraiserId *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.Contract.Withdraw(&_SudaFundraising.TransactOpts, fundraiserId)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 fundraiserId) returns()
func (_SudaFundraising *SudaFundraisingTransactorSession) Withdraw(fundraiserId *big.Int) (*types.Transaction, error) {
	return _SudaFundraising.Contract.Withdraw(&_SudaFundraising.TransactOpts, fundraiserId)
}

// SudaFundraisingDonatedIterator is returned from FilterDonated and is used to iterate over the raw logs and unpacked data for Donated events raised by the SudaFundraising contract.
type SudaFundraisingDonatedIterator struct {
	Event *SudaFundraisingDonated // Event containing the contract specifics and raw log

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
func (it *SudaFundraisingDonatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaFundraisingDonated)
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
		it.Event = new(SudaFundraisingDonated)
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
func (it *SudaFundraisingDonatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaFundraisingDonatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaFundraisingDonated represents a Donated event raised by the SudaFundraising contract.
type SudaFundraisingDonated struct {
	FundraiserId *big.Int
	Donor        common.Address
	Amount       *big.Int
	TotalRaised  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDonated is a free log retrieval operation binding the contract event 0x27e413b8931dc4cb7edc879292ebc482a3047336d03423ea91addae7a2cf850a.
//
// Solidity: event Donated(uint256 indexed fundraiserId, address indexed donor, uint256 amount, uint256 totalRaised)
func (_SudaFundraising *SudaFundraisingFilterer) FilterDonated(opts *bind.FilterOpts, fundraiserId []*big.Int, donor []common.Address) (*SudaFundraisingDonatedIterator, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}
	var donorRule []interface{}
	for _, donorItem := range donor {
		donorRule = append(donorRule, donorItem)
	}

	logs, sub, err := _SudaFundraising.contract.FilterLogs(opts, "Donated", fundraiserIdRule, donorRule)
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingDonatedIterator{contract: _SudaFundraising.contract, event: "Donated", logs: logs, sub: sub}, nil
}

// WatchDonated is a free log subscription operation binding the contract event 0x27e413b8931dc4cb7edc879292ebc482a3047336d03423ea91addae7a2cf850a.
//
// Solidity: event Donated(uint256 indexed fundraiserId, address indexed donor, uint256 amount, uint256 totalRaised)
func (_SudaFundraising *SudaFundraisingFilterer) WatchDonated(opts *bind.WatchOpts, sink chan<- *SudaFundraisingDonated, fundraiserId []*big.Int, donor []common.Address) (event.Subscription, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}
	var donorRule []interface{}
	for _, donorItem := range donor {
		donorRule = append(donorRule, donorItem)
	}

	logs, sub, err := _SudaFundraising.contract.WatchLogs(opts, "Donated", fundraiserIdRule, donorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaFundraisingDonated)
				if err := _SudaFundraising.contract.UnpackLog(event, "Donated", log); err != nil {
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

// ParseDonated is a log parse operation binding the contract event 0x27e413b8931dc4cb7edc879292ebc482a3047336d03423ea91addae7a2cf850a.
//
// Solidity: event Donated(uint256 indexed fundraiserId, address indexed donor, uint256 amount, uint256 totalRaised)
func (_SudaFundraising *SudaFundraisingFilterer) ParseDonated(log types.Log) (*SudaFundraisingDonated, error) {
	event := new(SudaFundraisingDonated)
	if err := _SudaFundraising.contract.UnpackLog(event, "Donated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaFundraisingExpiredIterator is returned from FilterExpired and is used to iterate over the raw logs and unpacked data for Expired events raised by the SudaFundraising contract.
type SudaFundraisingExpiredIterator struct {
	Event *SudaFundraisingExpired // Event containing the contract specifics and raw log

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
func (it *SudaFundraisingExpiredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaFundraisingExpired)
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
		it.Event = new(SudaFundraisingExpired)
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
func (it *SudaFundraisingExpiredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaFundraisingExpiredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaFundraisingExpired represents a Expired event raised by the SudaFundraising contract.
type SudaFundraisingExpired struct {
	FundraiserId *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterExpired is a free log retrieval operation binding the contract event 0xf80dbaea4785589e52984ca36a31de106adc77759539a5c7d92883bf49692fe9.
//
// Solidity: event Expired(uint256 indexed fundraiserId)
func (_SudaFundraising *SudaFundraisingFilterer) FilterExpired(opts *bind.FilterOpts, fundraiserId []*big.Int) (*SudaFundraisingExpiredIterator, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}

	logs, sub, err := _SudaFundraising.contract.FilterLogs(opts, "Expired", fundraiserIdRule)
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingExpiredIterator{contract: _SudaFundraising.contract, event: "Expired", logs: logs, sub: sub}, nil
}

// WatchExpired is a free log subscription operation binding the contract event 0xf80dbaea4785589e52984ca36a31de106adc77759539a5c7d92883bf49692fe9.
//
// Solidity: event Expired(uint256 indexed fundraiserId)
func (_SudaFundraising *SudaFundraisingFilterer) WatchExpired(opts *bind.WatchOpts, sink chan<- *SudaFundraisingExpired, fundraiserId []*big.Int) (event.Subscription, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}

	logs, sub, err := _SudaFundraising.contract.WatchLogs(opts, "Expired", fundraiserIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaFundraisingExpired)
				if err := _SudaFundraising.contract.UnpackLog(event, "Expired", log); err != nil {
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

// ParseExpired is a log parse operation binding the contract event 0xf80dbaea4785589e52984ca36a31de106adc77759539a5c7d92883bf49692fe9.
//
// Solidity: event Expired(uint256 indexed fundraiserId)
func (_SudaFundraising *SudaFundraisingFilterer) ParseExpired(log types.Log) (*SudaFundraisingExpired, error) {
	event := new(SudaFundraisingExpired)
	if err := _SudaFundraising.contract.UnpackLog(event, "Expired", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaFundraisingFundraiserCreatedIterator is returned from FilterFundraiserCreated and is used to iterate over the raw logs and unpacked data for FundraiserCreated events raised by the SudaFundraising contract.
type SudaFundraisingFundraiserCreatedIterator struct {
	Event *SudaFundraisingFundraiserCreated // Event containing the contract specifics and raw log

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
func (it *SudaFundraisingFundraiserCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaFundraisingFundraiserCreated)
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
		it.Event = new(SudaFundraisingFundraiserCreated)
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
func (it *SudaFundraisingFundraiserCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaFundraisingFundraiserCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaFundraisingFundraiserCreated represents a FundraiserCreated event raised by the SudaFundraising contract.
type SudaFundraisingFundraiserCreated struct {
	FundraiserId *big.Int
	Creator      common.Address
	Goal         *big.Int
	Deadline     *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterFundraiserCreated is a free log retrieval operation binding the contract event 0x4b54696510e948513db22b55242041175ae8be7658712a4da84b83eee7c336fd.
//
// Solidity: event FundraiserCreated(uint256 indexed fundraiserId, address indexed creator, uint256 goal, uint256 deadline)
func (_SudaFundraising *SudaFundraisingFilterer) FilterFundraiserCreated(opts *bind.FilterOpts, fundraiserId []*big.Int, creator []common.Address) (*SudaFundraisingFundraiserCreatedIterator, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SudaFundraising.contract.FilterLogs(opts, "FundraiserCreated", fundraiserIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingFundraiserCreatedIterator{contract: _SudaFundraising.contract, event: "FundraiserCreated", logs: logs, sub: sub}, nil
}

// WatchFundraiserCreated is a free log subscription operation binding the contract event 0x4b54696510e948513db22b55242041175ae8be7658712a4da84b83eee7c336fd.
//
// Solidity: event FundraiserCreated(uint256 indexed fundraiserId, address indexed creator, uint256 goal, uint256 deadline)
func (_SudaFundraising *SudaFundraisingFilterer) WatchFundraiserCreated(opts *bind.WatchOpts, sink chan<- *SudaFundraisingFundraiserCreated, fundraiserId []*big.Int, creator []common.Address) (event.Subscription, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SudaFundraising.contract.WatchLogs(opts, "FundraiserCreated", fundraiserIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaFundraisingFundraiserCreated)
				if err := _SudaFundraising.contract.UnpackLog(event, "FundraiserCreated", log); err != nil {
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

// ParseFundraiserCreated is a log parse operation binding the contract event 0x4b54696510e948513db22b55242041175ae8be7658712a4da84b83eee7c336fd.
//
// Solidity: event FundraiserCreated(uint256 indexed fundraiserId, address indexed creator, uint256 goal, uint256 deadline)
func (_SudaFundraising *SudaFundraisingFilterer) ParseFundraiserCreated(log types.Log) (*SudaFundraisingFundraiserCreated, error) {
	event := new(SudaFundraisingFundraiserCreated)
	if err := _SudaFundraising.contract.UnpackLog(event, "FundraiserCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaFundraisingGoalReachedIterator is returned from FilterGoalReached and is used to iterate over the raw logs and unpacked data for GoalReached events raised by the SudaFundraising contract.
type SudaFundraisingGoalReachedIterator struct {
	Event *SudaFundraisingGoalReached // Event containing the contract specifics and raw log

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
func (it *SudaFundraisingGoalReachedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaFundraisingGoalReached)
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
		it.Event = new(SudaFundraisingGoalReached)
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
func (it *SudaFundraisingGoalReachedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaFundraisingGoalReachedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaFundraisingGoalReached represents a GoalReached event raised by the SudaFundraising contract.
type SudaFundraisingGoalReached struct {
	FundraiserId *big.Int
	TotalRaised  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterGoalReached is a free log retrieval operation binding the contract event 0x85b3ed4e45559c5f41fb220aa4ac86a440dfc741f219089de694242940aaa09c.
//
// Solidity: event GoalReached(uint256 indexed fundraiserId, uint256 totalRaised)
func (_SudaFundraising *SudaFundraisingFilterer) FilterGoalReached(opts *bind.FilterOpts, fundraiserId []*big.Int) (*SudaFundraisingGoalReachedIterator, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}

	logs, sub, err := _SudaFundraising.contract.FilterLogs(opts, "GoalReached", fundraiserIdRule)
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingGoalReachedIterator{contract: _SudaFundraising.contract, event: "GoalReached", logs: logs, sub: sub}, nil
}

// WatchGoalReached is a free log subscription operation binding the contract event 0x85b3ed4e45559c5f41fb220aa4ac86a440dfc741f219089de694242940aaa09c.
//
// Solidity: event GoalReached(uint256 indexed fundraiserId, uint256 totalRaised)
func (_SudaFundraising *SudaFundraisingFilterer) WatchGoalReached(opts *bind.WatchOpts, sink chan<- *SudaFundraisingGoalReached, fundraiserId []*big.Int) (event.Subscription, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}

	logs, sub, err := _SudaFundraising.contract.WatchLogs(opts, "GoalReached", fundraiserIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaFundraisingGoalReached)
				if err := _SudaFundraising.contract.UnpackLog(event, "GoalReached", log); err != nil {
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

// ParseGoalReached is a log parse operation binding the contract event 0x85b3ed4e45559c5f41fb220aa4ac86a440dfc741f219089de694242940aaa09c.
//
// Solidity: event GoalReached(uint256 indexed fundraiserId, uint256 totalRaised)
func (_SudaFundraising *SudaFundraisingFilterer) ParseGoalReached(log types.Log) (*SudaFundraisingGoalReached, error) {
	event := new(SudaFundraisingGoalReached)
	if err := _SudaFundraising.contract.UnpackLog(event, "GoalReached", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaFundraisingOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SudaFundraising contract.
type SudaFundraisingOwnershipTransferredIterator struct {
	Event *SudaFundraisingOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SudaFundraisingOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaFundraisingOwnershipTransferred)
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
		it.Event = new(SudaFundraisingOwnershipTransferred)
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
func (it *SudaFundraisingOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaFundraisingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaFundraisingOwnershipTransferred represents a OwnershipTransferred event raised by the SudaFundraising contract.
type SudaFundraisingOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SudaFundraising *SudaFundraisingFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SudaFundraisingOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SudaFundraising.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingOwnershipTransferredIterator{contract: _SudaFundraising.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SudaFundraising *SudaFundraisingFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SudaFundraisingOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SudaFundraising.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaFundraisingOwnershipTransferred)
				if err := _SudaFundraising.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SudaFundraising *SudaFundraisingFilterer) ParseOwnershipTransferred(log types.Log) (*SudaFundraisingOwnershipTransferred, error) {
	event := new(SudaFundraisingOwnershipTransferred)
	if err := _SudaFundraising.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaFundraisingPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the SudaFundraising contract.
type SudaFundraisingPausedIterator struct {
	Event *SudaFundraisingPaused // Event containing the contract specifics and raw log

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
func (it *SudaFundraisingPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaFundraisingPaused)
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
		it.Event = new(SudaFundraisingPaused)
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
func (it *SudaFundraisingPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaFundraisingPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaFundraisingPaused represents a Paused event raised by the SudaFundraising contract.
type SudaFundraisingPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SudaFundraising *SudaFundraisingFilterer) FilterPaused(opts *bind.FilterOpts) (*SudaFundraisingPausedIterator, error) {

	logs, sub, err := _SudaFundraising.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingPausedIterator{contract: _SudaFundraising.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_SudaFundraising *SudaFundraisingFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *SudaFundraisingPaused) (event.Subscription, error) {

	logs, sub, err := _SudaFundraising.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaFundraisingPaused)
				if err := _SudaFundraising.contract.UnpackLog(event, "Paused", log); err != nil {
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
func (_SudaFundraising *SudaFundraisingFilterer) ParsePaused(log types.Log) (*SudaFundraisingPaused, error) {
	event := new(SudaFundraisingPaused)
	if err := _SudaFundraising.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaFundraisingRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the SudaFundraising contract.
type SudaFundraisingRefundedIterator struct {
	Event *SudaFundraisingRefunded // Event containing the contract specifics and raw log

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
func (it *SudaFundraisingRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaFundraisingRefunded)
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
		it.Event = new(SudaFundraisingRefunded)
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
func (it *SudaFundraisingRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaFundraisingRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaFundraisingRefunded represents a Refunded event raised by the SudaFundraising contract.
type SudaFundraisingRefunded struct {
	FundraiserId *big.Int
	Donor        common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0x7ca5472b7ea78c2c0141c5a12ee6d170cf4ce8ed06be3d22c8252ddfc7a6a2c4.
//
// Solidity: event Refunded(uint256 indexed fundraiserId, address indexed donor, uint256 amount)
func (_SudaFundraising *SudaFundraisingFilterer) FilterRefunded(opts *bind.FilterOpts, fundraiserId []*big.Int, donor []common.Address) (*SudaFundraisingRefundedIterator, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}
	var donorRule []interface{}
	for _, donorItem := range donor {
		donorRule = append(donorRule, donorItem)
	}

	logs, sub, err := _SudaFundraising.contract.FilterLogs(opts, "Refunded", fundraiserIdRule, donorRule)
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingRefundedIterator{contract: _SudaFundraising.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0x7ca5472b7ea78c2c0141c5a12ee6d170cf4ce8ed06be3d22c8252ddfc7a6a2c4.
//
// Solidity: event Refunded(uint256 indexed fundraiserId, address indexed donor, uint256 amount)
func (_SudaFundraising *SudaFundraisingFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *SudaFundraisingRefunded, fundraiserId []*big.Int, donor []common.Address) (event.Subscription, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}
	var donorRule []interface{}
	for _, donorItem := range donor {
		donorRule = append(donorRule, donorItem)
	}

	logs, sub, err := _SudaFundraising.contract.WatchLogs(opts, "Refunded", fundraiserIdRule, donorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaFundraisingRefunded)
				if err := _SudaFundraising.contract.UnpackLog(event, "Refunded", log); err != nil {
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

// ParseRefunded is a log parse operation binding the contract event 0x7ca5472b7ea78c2c0141c5a12ee6d170cf4ce8ed06be3d22c8252ddfc7a6a2c4.
//
// Solidity: event Refunded(uint256 indexed fundraiserId, address indexed donor, uint256 amount)
func (_SudaFundraising *SudaFundraisingFilterer) ParseRefunded(log types.Log) (*SudaFundraisingRefunded, error) {
	event := new(SudaFundraisingRefunded)
	if err := _SudaFundraising.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaFundraisingUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the SudaFundraising contract.
type SudaFundraisingUnpausedIterator struct {
	Event *SudaFundraisingUnpaused // Event containing the contract specifics and raw log

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
func (it *SudaFundraisingUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaFundraisingUnpaused)
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
		it.Event = new(SudaFundraisingUnpaused)
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
func (it *SudaFundraisingUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaFundraisingUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaFundraisingUnpaused represents a Unpaused event raised by the SudaFundraising contract.
type SudaFundraisingUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SudaFundraising *SudaFundraisingFilterer) FilterUnpaused(opts *bind.FilterOpts) (*SudaFundraisingUnpausedIterator, error) {

	logs, sub, err := _SudaFundraising.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingUnpausedIterator{contract: _SudaFundraising.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_SudaFundraising *SudaFundraisingFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *SudaFundraisingUnpaused) (event.Subscription, error) {

	logs, sub, err := _SudaFundraising.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaFundraisingUnpaused)
				if err := _SudaFundraising.contract.UnpackLog(event, "Unpaused", log); err != nil {
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
func (_SudaFundraising *SudaFundraisingFilterer) ParseUnpaused(log types.Log) (*SudaFundraisingUnpaused, error) {
	event := new(SudaFundraisingUnpaused)
	if err := _SudaFundraising.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SudaFundraisingWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the SudaFundraising contract.
type SudaFundraisingWithdrawnIterator struct {
	Event *SudaFundraisingWithdrawn // Event containing the contract specifics and raw log

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
func (it *SudaFundraisingWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SudaFundraisingWithdrawn)
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
		it.Event = new(SudaFundraisingWithdrawn)
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
func (it *SudaFundraisingWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SudaFundraisingWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SudaFundraisingWithdrawn represents a Withdrawn event raised by the SudaFundraising contract.
type SudaFundraisingWithdrawn struct {
	FundraiserId *big.Int
	Creator      common.Address
	Amount       *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0xcf7d23a3cbe4e8b36ff82fd1b05b1b17373dc7804b4ebbd6e2356716ef202372.
//
// Solidity: event Withdrawn(uint256 indexed fundraiserId, address indexed creator, uint256 amount)
func (_SudaFundraising *SudaFundraisingFilterer) FilterWithdrawn(opts *bind.FilterOpts, fundraiserId []*big.Int, creator []common.Address) (*SudaFundraisingWithdrawnIterator, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SudaFundraising.contract.FilterLogs(opts, "Withdrawn", fundraiserIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return &SudaFundraisingWithdrawnIterator{contract: _SudaFundraising.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0xcf7d23a3cbe4e8b36ff82fd1b05b1b17373dc7804b4ebbd6e2356716ef202372.
//
// Solidity: event Withdrawn(uint256 indexed fundraiserId, address indexed creator, uint256 amount)
func (_SudaFundraising *SudaFundraisingFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *SudaFundraisingWithdrawn, fundraiserId []*big.Int, creator []common.Address) (event.Subscription, error) {

	var fundraiserIdRule []interface{}
	for _, fundraiserIdItem := range fundraiserId {
		fundraiserIdRule = append(fundraiserIdRule, fundraiserIdItem)
	}
	var creatorRule []interface{}
	for _, creatorItem := range creator {
		creatorRule = append(creatorRule, creatorItem)
	}

	logs, sub, err := _SudaFundraising.contract.WatchLogs(opts, "Withdrawn", fundraiserIdRule, creatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SudaFundraisingWithdrawn)
				if err := _SudaFundraising.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0xcf7d23a3cbe4e8b36ff82fd1b05b1b17373dc7804b4ebbd6e2356716ef202372.
//
// Solidity: event Withdrawn(uint256 indexed fundraiserId, address indexed creator, uint256 amount)
func (_SudaFundraising *SudaFundraisingFilterer) ParseWithdrawn(log types.Log) (*SudaFundraisingWithdrawn, error) {
	event := new(SudaFundraisingWithdrawn)
	if err := _SudaFundraising.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
