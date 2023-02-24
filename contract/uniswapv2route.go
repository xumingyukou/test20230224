// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

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
)

// Uniswapv2routeMetaData contains all meta data concerning the Uniswapv2route contract.
var Uniswapv2routeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"swapFeeReward\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WETH\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapMining\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tradingPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountADesired\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBDesired\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountAMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"addLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenDesired\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"addLiquidityETH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveOut\",\"type\":\"uint256\"}],\"name\":\"getAmountIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveOut\",\"type\":\"uint256\"}],\"name\":\"getAmountOut\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"}],\"name\":\"getAmountsIn\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"}],\"name\":\"getAmountsOut\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveB\",\"type\":\"uint256\"}],\"name\":\"quote\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountAMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"removeLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"removeLiquidityETH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"removeLiquidityETHSupportingFeeOnTransferTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"approveMax\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"removeLiquidityETHWithPermit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"approveMax\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"removeLiquidityETHWithPermitSupportingFeeOnTransferTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountAMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"approveMax\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"removeLiquidityWithPermit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapETHForExactTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactETHForTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactETHForTokensSupportingFeeOnTransferTokens\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForETH\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForHT\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForETHSupportingFeeOnTransferTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForTokensSupportingFeeOnTransferTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMax\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapTokensForExactETH\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactHTForTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMax\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapTokensForExactTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// Uniswapv2routeABI is the input ABI used to generate the binding from.
// Deprecated: Use Uniswapv2routeMetaData.ABI instead.
var Uniswapv2routeABI = Uniswapv2routeMetaData.ABI

// Uniswapv2route is an auto generated Go binding around an Ethereum contract.
type Uniswapv2route struct {
	Uniswapv2routeCaller     // Read-only binding to the contract
	Uniswapv2routeTransactor // Write-only binding to the contract
	Uniswapv2routeFilterer   // Log filterer for contract events
}

// Uniswapv2routeCaller is an auto generated read-only Go binding around an Ethereum contract.
type Uniswapv2routeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Uniswapv2routeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Uniswapv2routeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Uniswapv2routeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Uniswapv2routeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Uniswapv2routeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Uniswapv2routeSession struct {
	Contract     *Uniswapv2route   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Uniswapv2routeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Uniswapv2routeCallerSession struct {
	Contract *Uniswapv2routeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// Uniswapv2routeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Uniswapv2routeTransactorSession struct {
	Contract     *Uniswapv2routeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// Uniswapv2routeRaw is an auto generated low-level Go binding around an Ethereum contract.
type Uniswapv2routeRaw struct {
	Contract *Uniswapv2route // Generic contract binding to access the raw methods on
}

// Uniswapv2routeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Uniswapv2routeCallerRaw struct {
	Contract *Uniswapv2routeCaller // Generic read-only contract binding to access the raw methods on
}

// Uniswapv2routeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Uniswapv2routeTransactorRaw struct {
	Contract *Uniswapv2routeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapv2route creates a new instance of Uniswapv2route, bound to a specific deployed contract.
func NewUniswapv2route(address common.Address, backend bind.ContractBackend) (*Uniswapv2route, error) {
	contract, err := bindUniswapv2route(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Uniswapv2route{Uniswapv2routeCaller: Uniswapv2routeCaller{contract: contract}, Uniswapv2routeTransactor: Uniswapv2routeTransactor{contract: contract}, Uniswapv2routeFilterer: Uniswapv2routeFilterer{contract: contract}}, nil
}

// NewUniswapv2routeCaller creates a new read-only instance of Uniswapv2route, bound to a specific deployed contract.
func NewUniswapv2routeCaller(address common.Address, caller bind.ContractCaller) (*Uniswapv2routeCaller, error) {
	contract, err := bindUniswapv2route(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Uniswapv2routeCaller{contract: contract}, nil
}

// NewUniswapv2routeTransactor creates a new write-only instance of Uniswapv2route, bound to a specific deployed contract.
func NewUniswapv2routeTransactor(address common.Address, transactor bind.ContractTransactor) (*Uniswapv2routeTransactor, error) {
	contract, err := bindUniswapv2route(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Uniswapv2routeTransactor{contract: contract}, nil
}

// NewUniswapv2routeFilterer creates a new log filterer instance of Uniswapv2route, bound to a specific deployed contract.
func NewUniswapv2routeFilterer(address common.Address, filterer bind.ContractFilterer) (*Uniswapv2routeFilterer, error) {
	contract, err := bindUniswapv2route(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Uniswapv2routeFilterer{contract: contract}, nil
}

// bindUniswapv2route binds a generic wrapper to an already deployed contract.
func bindUniswapv2route(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Uniswapv2routeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Uniswapv2route *Uniswapv2routeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Uniswapv2route.Contract.Uniswapv2routeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Uniswapv2route *Uniswapv2routeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.Uniswapv2routeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Uniswapv2route *Uniswapv2routeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.Uniswapv2routeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Uniswapv2route *Uniswapv2routeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Uniswapv2route.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Uniswapv2route *Uniswapv2routeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Uniswapv2route *Uniswapv2routeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.contract.Transact(opts, method, params...)
}

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() pure returns(address)
func (_Uniswapv2route *Uniswapv2routeCaller) WETH(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswapv2route.contract.Call(opts, &out, "WETH")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() pure returns(address)
func (_Uniswapv2route *Uniswapv2routeSession) WETH() (common.Address, error) {
	return _Uniswapv2route.Contract.WETH(&_Uniswapv2route.CallOpts)
}

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() pure returns(address)
func (_Uniswapv2route *Uniswapv2routeCallerSession) WETH() (common.Address, error) {
	return _Uniswapv2route.Contract.WETH(&_Uniswapv2route.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() pure returns(address)
func (_Uniswapv2route *Uniswapv2routeCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswapv2route.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() pure returns(address)
func (_Uniswapv2route *Uniswapv2routeSession) Factory() (common.Address, error) {
	return _Uniswapv2route.Contract.Factory(&_Uniswapv2route.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() pure returns(address)
func (_Uniswapv2route *Uniswapv2routeCallerSession) Factory() (common.Address, error) {
	return _Uniswapv2route.Contract.Factory(&_Uniswapv2route.CallOpts)
}

// GetAmountIn is a free data retrieval call binding the contract method 0x85f8c259.
//
// Solidity: function getAmountIn(uint256 amountOut, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountIn)
func (_Uniswapv2route *Uniswapv2routeCaller) GetAmountIn(opts *bind.CallOpts, amountOut *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Uniswapv2route.contract.Call(opts, &out, "getAmountIn", amountOut, reserveIn, reserveOut)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountIn is a free data retrieval call binding the contract method 0x85f8c259.
//
// Solidity: function getAmountIn(uint256 amountOut, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountIn)
func (_Uniswapv2route *Uniswapv2routeSession) GetAmountIn(amountOut *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	return _Uniswapv2route.Contract.GetAmountIn(&_Uniswapv2route.CallOpts, amountOut, reserveIn, reserveOut)
}

// GetAmountIn is a free data retrieval call binding the contract method 0x85f8c259.
//
// Solidity: function getAmountIn(uint256 amountOut, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountIn)
func (_Uniswapv2route *Uniswapv2routeCallerSession) GetAmountIn(amountOut *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	return _Uniswapv2route.Contract.GetAmountIn(&_Uniswapv2route.CallOpts, amountOut, reserveIn, reserveOut)
}

// GetAmountOut is a free data retrieval call binding the contract method 0x054d50d4.
//
// Solidity: function getAmountOut(uint256 amountIn, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountOut)
func (_Uniswapv2route *Uniswapv2routeCaller) GetAmountOut(opts *bind.CallOpts, amountIn *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Uniswapv2route.contract.Call(opts, &out, "getAmountOut", amountIn, reserveIn, reserveOut)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountOut is a free data retrieval call binding the contract method 0x054d50d4.
//
// Solidity: function getAmountOut(uint256 amountIn, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountOut)
func (_Uniswapv2route *Uniswapv2routeSession) GetAmountOut(amountIn *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	return _Uniswapv2route.Contract.GetAmountOut(&_Uniswapv2route.CallOpts, amountIn, reserveIn, reserveOut)
}

// GetAmountOut is a free data retrieval call binding the contract method 0x054d50d4.
//
// Solidity: function getAmountOut(uint256 amountIn, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountOut)
func (_Uniswapv2route *Uniswapv2routeCallerSession) GetAmountOut(amountIn *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	return _Uniswapv2route.Contract.GetAmountOut(&_Uniswapv2route.CallOpts, amountIn, reserveIn, reserveOut)
}

// GetAmountsIn is a free data retrieval call binding the contract method 0x1f00ca74.
//
// Solidity: function getAmountsIn(uint256 amountOut, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeCaller) GetAmountsIn(opts *bind.CallOpts, amountOut *big.Int, path []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _Uniswapv2route.contract.Call(opts, &out, "getAmountsIn", amountOut, path)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAmountsIn is a free data retrieval call binding the contract method 0x1f00ca74.
//
// Solidity: function getAmountsIn(uint256 amountOut, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeSession) GetAmountsIn(amountOut *big.Int, path []common.Address) ([]*big.Int, error) {
	return _Uniswapv2route.Contract.GetAmountsIn(&_Uniswapv2route.CallOpts, amountOut, path)
}

// GetAmountsIn is a free data retrieval call binding the contract method 0x1f00ca74.
//
// Solidity: function getAmountsIn(uint256 amountOut, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeCallerSession) GetAmountsIn(amountOut *big.Int, path []common.Address) ([]*big.Int, error) {
	return _Uniswapv2route.Contract.GetAmountsIn(&_Uniswapv2route.CallOpts, amountOut, path)
}

// GetAmountsOut is a free data retrieval call binding the contract method 0xd06ca61f.
//
// Solidity: function getAmountsOut(uint256 amountIn, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeCaller) GetAmountsOut(opts *bind.CallOpts, amountIn *big.Int, path []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _Uniswapv2route.contract.Call(opts, &out, "getAmountsOut", amountIn, path)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAmountsOut is a free data retrieval call binding the contract method 0xd06ca61f.
//
// Solidity: function getAmountsOut(uint256 amountIn, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeSession) GetAmountsOut(amountIn *big.Int, path []common.Address) ([]*big.Int, error) {
	return _Uniswapv2route.Contract.GetAmountsOut(&_Uniswapv2route.CallOpts, amountIn, path)
}

// GetAmountsOut is a free data retrieval call binding the contract method 0xd06ca61f.
//
// Solidity: function getAmountsOut(uint256 amountIn, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeCallerSession) GetAmountsOut(amountIn *big.Int, path []common.Address) ([]*big.Int, error) {
	return _Uniswapv2route.Contract.GetAmountsOut(&_Uniswapv2route.CallOpts, amountIn, path)
}

// Quote is a free data retrieval call binding the contract method 0xad615dec.
//
// Solidity: function quote(uint256 amountA, uint256 reserveA, uint256 reserveB) pure returns(uint256 amountB)
func (_Uniswapv2route *Uniswapv2routeCaller) Quote(opts *bind.CallOpts, amountA *big.Int, reserveA *big.Int, reserveB *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Uniswapv2route.contract.Call(opts, &out, "quote", amountA, reserveA, reserveB)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Quote is a free data retrieval call binding the contract method 0xad615dec.
//
// Solidity: function quote(uint256 amountA, uint256 reserveA, uint256 reserveB) pure returns(uint256 amountB)
func (_Uniswapv2route *Uniswapv2routeSession) Quote(amountA *big.Int, reserveA *big.Int, reserveB *big.Int) (*big.Int, error) {
	return _Uniswapv2route.Contract.Quote(&_Uniswapv2route.CallOpts, amountA, reserveA, reserveB)
}

// Quote is a free data retrieval call binding the contract method 0xad615dec.
//
// Solidity: function quote(uint256 amountA, uint256 reserveA, uint256 reserveB) pure returns(uint256 amountB)
func (_Uniswapv2route *Uniswapv2routeCallerSession) Quote(amountA *big.Int, reserveA *big.Int, reserveB *big.Int) (*big.Int, error) {
	return _Uniswapv2route.Contract.Quote(&_Uniswapv2route.CallOpts, amountA, reserveA, reserveB)
}

// SwapFeeReward is a free data retrieval call binding the contract method 0xa3ddb30b.
//
// Solidity: function swapFeeReward() pure returns(address)
func (_Uniswapv2route *Uniswapv2routeCaller) SwapFeeReward(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswapv2route.contract.Call(opts, &out, "swapFeeReward")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SwapFeeReward is a free data retrieval call binding the contract method 0xa3ddb30b.
//
// Solidity: function swapFeeReward() pure returns(address)
func (_Uniswapv2route *Uniswapv2routeSession) SwapFeeReward() (common.Address, error) {
	return _Uniswapv2route.Contract.SwapFeeReward(&_Uniswapv2route.CallOpts)
}

// SwapFeeReward is a free data retrieval call binding the contract method 0xa3ddb30b.
//
// Solidity: function swapFeeReward() pure returns(address)
func (_Uniswapv2route *Uniswapv2routeCallerSession) SwapFeeReward() (common.Address, error) {
	return _Uniswapv2route.Contract.SwapFeeReward(&_Uniswapv2route.CallOpts)
}

// SwapMining is a free data retrieval call binding the contract method 0x4e5eda11.
//
// Solidity: function swapMining() view returns(address)
func (_Uniswapv2route *Uniswapv2routeCaller) SwapMining(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswapv2route.contract.Call(opts, &out, "swapMining")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SwapMining is a free data retrieval call binding the contract method 0x4e5eda11.
//
// Solidity: function swapMining() view returns(address)
func (_Uniswapv2route *Uniswapv2routeSession) SwapMining() (common.Address, error) {
	return _Uniswapv2route.Contract.SwapMining(&_Uniswapv2route.CallOpts)
}

// SwapMining is a free data retrieval call binding the contract method 0x4e5eda11.
//
// Solidity: function swapMining() view returns(address)
func (_Uniswapv2route *Uniswapv2routeCallerSession) SwapMining() (common.Address, error) {
	return _Uniswapv2route.Contract.SwapMining(&_Uniswapv2route.CallOpts)
}

// TradingPool is a free data retrieval call binding the contract method 0x71f16a18.
//
// Solidity: function tradingPool() view returns(address)
func (_Uniswapv2route *Uniswapv2routeCaller) TradingPool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswapv2route.contract.Call(opts, &out, "tradingPool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TradingPool is a free data retrieval call binding the contract method 0x71f16a18.
//
// Solidity: function tradingPool() view returns(address)
func (_Uniswapv2route *Uniswapv2routeSession) TradingPool() (common.Address, error) {
	return _Uniswapv2route.Contract.TradingPool(&_Uniswapv2route.CallOpts)
}

// TradingPool is a free data retrieval call binding the contract method 0x71f16a18.
//
// Solidity: function tradingPool() view returns(address)
func (_Uniswapv2route *Uniswapv2routeCallerSession) TradingPool() (common.Address, error) {
	return _Uniswapv2route.Contract.TradingPool(&_Uniswapv2route.CallOpts)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0xe8e33700.
//
// Solidity: function addLiquidity(address tokenA, address tokenB, uint256 amountADesired, uint256 amountBDesired, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB, uint256 liquidity)
func (_Uniswapv2route *Uniswapv2routeTransactor) AddLiquidity(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address, amountADesired *big.Int, amountBDesired *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "addLiquidity", tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, to, deadline)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0xe8e33700.
//
// Solidity: function addLiquidity(address tokenA, address tokenB, uint256 amountADesired, uint256 amountBDesired, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB, uint256 liquidity)
func (_Uniswapv2route *Uniswapv2routeSession) AddLiquidity(tokenA common.Address, tokenB common.Address, amountADesired *big.Int, amountBDesired *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.AddLiquidity(&_Uniswapv2route.TransactOpts, tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, to, deadline)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0xe8e33700.
//
// Solidity: function addLiquidity(address tokenA, address tokenB, uint256 amountADesired, uint256 amountBDesired, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB, uint256 liquidity)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) AddLiquidity(tokenA common.Address, tokenB common.Address, amountADesired *big.Int, amountBDesired *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.AddLiquidity(&_Uniswapv2route.TransactOpts, tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, to, deadline)
}

// AddLiquidityETH is a paid mutator transaction binding the contract method 0xf305d719.
//
// Solidity: function addLiquidityETH(address token, uint256 amountTokenDesired, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) payable returns(uint256 amountToken, uint256 amountETH, uint256 liquidity)
func (_Uniswapv2route *Uniswapv2routeTransactor) AddLiquidityETH(opts *bind.TransactOpts, token common.Address, amountTokenDesired *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "addLiquidityETH", token, amountTokenDesired, amountTokenMin, amountETHMin, to, deadline)
}

// AddLiquidityETH is a paid mutator transaction binding the contract method 0xf305d719.
//
// Solidity: function addLiquidityETH(address token, uint256 amountTokenDesired, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) payable returns(uint256 amountToken, uint256 amountETH, uint256 liquidity)
func (_Uniswapv2route *Uniswapv2routeSession) AddLiquidityETH(token common.Address, amountTokenDesired *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.AddLiquidityETH(&_Uniswapv2route.TransactOpts, token, amountTokenDesired, amountTokenMin, amountETHMin, to, deadline)
}

// AddLiquidityETH is a paid mutator transaction binding the contract method 0xf305d719.
//
// Solidity: function addLiquidityETH(address token, uint256 amountTokenDesired, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) payable returns(uint256 amountToken, uint256 amountETH, uint256 liquidity)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) AddLiquidityETH(token common.Address, amountTokenDesired *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.AddLiquidityETH(&_Uniswapv2route.TransactOpts, token, amountTokenDesired, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xbaa2abde.
//
// Solidity: function removeLiquidity(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2route *Uniswapv2routeTransactor) RemoveLiquidity(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "removeLiquidity", tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xbaa2abde.
//
// Solidity: function removeLiquidity(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2route *Uniswapv2routeSession) RemoveLiquidity(tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidity(&_Uniswapv2route.TransactOpts, tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xbaa2abde.
//
// Solidity: function removeLiquidity(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) RemoveLiquidity(tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidity(&_Uniswapv2route.TransactOpts, tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline)
}

// RemoveLiquidityETH is a paid mutator transaction binding the contract method 0x02751cec.
//
// Solidity: function removeLiquidityETH(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeTransactor) RemoveLiquidityETH(opts *bind.TransactOpts, token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "removeLiquidityETH", token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETH is a paid mutator transaction binding the contract method 0x02751cec.
//
// Solidity: function removeLiquidityETH(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeSession) RemoveLiquidityETH(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidityETH(&_Uniswapv2route.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETH is a paid mutator transaction binding the contract method 0x02751cec.
//
// Solidity: function removeLiquidityETH(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) RemoveLiquidityETH(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidityETH(&_Uniswapv2route.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xaf2979eb.
//
// Solidity: function removeLiquidityETHSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeTransactor) RemoveLiquidityETHSupportingFeeOnTransferTokens(opts *bind.TransactOpts, token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "removeLiquidityETHSupportingFeeOnTransferTokens", token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xaf2979eb.
//
// Solidity: function removeLiquidityETHSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeSession) RemoveLiquidityETHSupportingFeeOnTransferTokens(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidityETHSupportingFeeOnTransferTokens(&_Uniswapv2route.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xaf2979eb.
//
// Solidity: function removeLiquidityETHSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) RemoveLiquidityETHSupportingFeeOnTransferTokens(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidityETHSupportingFeeOnTransferTokens(&_Uniswapv2route.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETHWithPermit is a paid mutator transaction binding the contract method 0xded9382a.
//
// Solidity: function removeLiquidityETHWithPermit(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeTransactor) RemoveLiquidityETHWithPermit(opts *bind.TransactOpts, token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "removeLiquidityETHWithPermit", token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityETHWithPermit is a paid mutator transaction binding the contract method 0xded9382a.
//
// Solidity: function removeLiquidityETHWithPermit(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeSession) RemoveLiquidityETHWithPermit(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidityETHWithPermit(&_Uniswapv2route.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityETHWithPermit is a paid mutator transaction binding the contract method 0xded9382a.
//
// Solidity: function removeLiquidityETHWithPermit(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) RemoveLiquidityETHWithPermit(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidityETHWithPermit(&_Uniswapv2route.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5b0d5984.
//
// Solidity: function removeLiquidityETHWithPermitSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeTransactor) RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(opts *bind.TransactOpts, token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "removeLiquidityETHWithPermitSupportingFeeOnTransferTokens", token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5b0d5984.
//
// Solidity: function removeLiquidityETHWithPermitSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeSession) RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(&_Uniswapv2route.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5b0d5984.
//
// Solidity: function removeLiquidityETHWithPermitSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountETH)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(&_Uniswapv2route.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityWithPermit is a paid mutator transaction binding the contract method 0x2195995c.
//
// Solidity: function removeLiquidityWithPermit(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2route *Uniswapv2routeTransactor) RemoveLiquidityWithPermit(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "removeLiquidityWithPermit", tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityWithPermit is a paid mutator transaction binding the contract method 0x2195995c.
//
// Solidity: function removeLiquidityWithPermit(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2route *Uniswapv2routeSession) RemoveLiquidityWithPermit(tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidityWithPermit(&_Uniswapv2route.TransactOpts, tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityWithPermit is a paid mutator transaction binding the contract method 0x2195995c.
//
// Solidity: function removeLiquidityWithPermit(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) RemoveLiquidityWithPermit(tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.RemoveLiquidityWithPermit(&_Uniswapv2route.TransactOpts, tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline, approveMax, v, r, s)
}

// SwapETHForExactTokens is a paid mutator transaction binding the contract method 0xfb3bdb41.
//
// Solidity: function swapETHForExactTokens(uint256 amountOut, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactor) SwapETHForExactTokens(opts *bind.TransactOpts, amountOut *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "swapETHForExactTokens", amountOut, path, to, deadline)
}

// SwapETHForExactTokens is a paid mutator transaction binding the contract method 0xfb3bdb41.
//
// Solidity: function swapETHForExactTokens(uint256 amountOut, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeSession) SwapETHForExactTokens(amountOut *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapETHForExactTokens(&_Uniswapv2route.TransactOpts, amountOut, path, to, deadline)
}

// SwapETHForExactTokens is a paid mutator transaction binding the contract method 0xfb3bdb41.
//
// Solidity: function swapETHForExactTokens(uint256 amountOut, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) SwapETHForExactTokens(amountOut *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapETHForExactTokens(&_Uniswapv2route.TransactOpts, amountOut, path, to, deadline)
}

// SwapExactETHForTokens is a paid mutator transaction binding the contract method 0x7ff36ab5.
//
// Solidity: function swapExactETHForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactor) SwapExactETHForTokens(opts *bind.TransactOpts, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "swapExactETHForTokens", amountOutMin, path, to, deadline)
}

// SwapExactETHForTokens is a paid mutator transaction binding the contract method 0x7ff36ab5.
//
// Solidity: function swapExactETHForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeSession) SwapExactETHForTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactETHForTokens(&_Uniswapv2route.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactETHForTokens is a paid mutator transaction binding the contract method 0x7ff36ab5.
//
// Solidity: function swapExactETHForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) SwapExactETHForTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactETHForTokens(&_Uniswapv2route.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactETHForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xb6f9de95.
//
// Solidity: function swapExactETHForTokensSupportingFeeOnTransferTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns()
func (_Uniswapv2route *Uniswapv2routeTransactor) SwapExactETHForTokensSupportingFeeOnTransferTokens(opts *bind.TransactOpts, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "swapExactETHForTokensSupportingFeeOnTransferTokens", amountOutMin, path, to, deadline)
}

// SwapExactETHForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xb6f9de95.
//
// Solidity: function swapExactETHForTokensSupportingFeeOnTransferTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns()
func (_Uniswapv2route *Uniswapv2routeSession) SwapExactETHForTokensSupportingFeeOnTransferTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactETHForTokensSupportingFeeOnTransferTokens(&_Uniswapv2route.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactETHForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xb6f9de95.
//
// Solidity: function swapExactETHForTokensSupportingFeeOnTransferTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns()
func (_Uniswapv2route *Uniswapv2routeTransactorSession) SwapExactETHForTokensSupportingFeeOnTransferTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactETHForTokensSupportingFeeOnTransferTokens(&_Uniswapv2route.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactHTForTokens is a paid mutator transaction binding the contract method 0xde6b8256.
//
// Solidity: function swapExactHTForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactor) SwapExactHTForTokens(opts *bind.TransactOpts, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "swapExactHTForTokens", amountOutMin, path, to, deadline)
}

// SwapExactHTForTokens is a paid mutator transaction binding the contract method 0xde6b8256.
//
// Solidity: function swapExactHTForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeSession) SwapExactHTForTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactHTForTokens(&_Uniswapv2route.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactHTForTokens is a paid mutator transaction binding the contract method 0xde6b8256.
//
// Solidity: function swapExactHTForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) SwapExactHTForTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactHTForTokens(&_Uniswapv2route.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETH is a paid mutator transaction binding the contract method 0x18cbafe5.
//
// Solidity: function swapExactTokensForETH(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactor) SwapExactTokensForETH(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "swapExactTokensForETH", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETH is a paid mutator transaction binding the contract method 0x18cbafe5.
//
// Solidity: function swapExactTokensForETH(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeSession) SwapExactTokensForETH(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactTokensForETH(&_Uniswapv2route.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETH is a paid mutator transaction binding the contract method 0x18cbafe5.
//
// Solidity: function swapExactTokensForETH(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) SwapExactTokensForETH(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactTokensForETH(&_Uniswapv2route.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x791ac947.
//
// Solidity: function swapExactTokensForETHSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2route *Uniswapv2routeTransactor) SwapExactTokensForETHSupportingFeeOnTransferTokens(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "swapExactTokensForETHSupportingFeeOnTransferTokens", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x791ac947.
//
// Solidity: function swapExactTokensForETHSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2route *Uniswapv2routeSession) SwapExactTokensForETHSupportingFeeOnTransferTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactTokensForETHSupportingFeeOnTransferTokens(&_Uniswapv2route.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x791ac947.
//
// Solidity: function swapExactTokensForETHSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2route *Uniswapv2routeTransactorSession) SwapExactTokensForETHSupportingFeeOnTransferTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactTokensForETHSupportingFeeOnTransferTokens(&_Uniswapv2route.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForHT is a paid mutator transaction binding the contract method 0x4b361208.
//
// Solidity: function swapExactTokensForHT(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactor) SwapExactTokensForHT(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "swapExactTokensForHT", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForHT is a paid mutator transaction binding the contract method 0x4b361208.
//
// Solidity: function swapExactTokensForHT(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeSession) SwapExactTokensForHT(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactTokensForHT(&_Uniswapv2route.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForHT is a paid mutator transaction binding the contract method 0x4b361208.
//
// Solidity: function swapExactTokensForHT(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) SwapExactTokensForHT(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactTokensForHT(&_Uniswapv2route.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0x38ed1739.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactor) SwapExactTokensForTokens(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "swapExactTokensForTokens", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0x38ed1739.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeSession) SwapExactTokensForTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactTokensForTokens(&_Uniswapv2route.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0x38ed1739.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) SwapExactTokensForTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactTokensForTokens(&_Uniswapv2route.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5c11d795.
//
// Solidity: function swapExactTokensForTokensSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2route *Uniswapv2routeTransactor) SwapExactTokensForTokensSupportingFeeOnTransferTokens(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "swapExactTokensForTokensSupportingFeeOnTransferTokens", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5c11d795.
//
// Solidity: function swapExactTokensForTokensSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2route *Uniswapv2routeSession) SwapExactTokensForTokensSupportingFeeOnTransferTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactTokensForTokensSupportingFeeOnTransferTokens(&_Uniswapv2route.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5c11d795.
//
// Solidity: function swapExactTokensForTokensSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2route *Uniswapv2routeTransactorSession) SwapExactTokensForTokensSupportingFeeOnTransferTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapExactTokensForTokensSupportingFeeOnTransferTokens(&_Uniswapv2route.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapTokensForExactETH is a paid mutator transaction binding the contract method 0x4a25d94a.
//
// Solidity: function swapTokensForExactETH(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactor) SwapTokensForExactETH(opts *bind.TransactOpts, amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "swapTokensForExactETH", amountOut, amountInMax, path, to, deadline)
}

// SwapTokensForExactETH is a paid mutator transaction binding the contract method 0x4a25d94a.
//
// Solidity: function swapTokensForExactETH(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeSession) SwapTokensForExactETH(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapTokensForExactETH(&_Uniswapv2route.TransactOpts, amountOut, amountInMax, path, to, deadline)
}

// SwapTokensForExactETH is a paid mutator transaction binding the contract method 0x4a25d94a.
//
// Solidity: function swapTokensForExactETH(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) SwapTokensForExactETH(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapTokensForExactETH(&_Uniswapv2route.TransactOpts, amountOut, amountInMax, path, to, deadline)
}

// SwapTokensForExactTokens is a paid mutator transaction binding the contract method 0x8803dbee.
//
// Solidity: function swapTokensForExactTokens(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactor) SwapTokensForExactTokens(opts *bind.TransactOpts, amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.contract.Transact(opts, "swapTokensForExactTokens", amountOut, amountInMax, path, to, deadline)
}

// SwapTokensForExactTokens is a paid mutator transaction binding the contract method 0x8803dbee.
//
// Solidity: function swapTokensForExactTokens(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeSession) SwapTokensForExactTokens(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapTokensForExactTokens(&_Uniswapv2route.TransactOpts, amountOut, amountInMax, path, to, deadline)
}

// SwapTokensForExactTokens is a paid mutator transaction binding the contract method 0x8803dbee.
//
// Solidity: function swapTokensForExactTokens(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2route *Uniswapv2routeTransactorSession) SwapTokensForExactTokens(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2route.Contract.SwapTokensForExactTokens(&_Uniswapv2route.TransactOpts, amountOut, amountInMax, path, to, deadline)
}
