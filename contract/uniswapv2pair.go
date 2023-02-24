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

// Uniswapv2pairMetaData contains all meta data concerning the Uniswapv2pair contract.
var Uniswapv2pairMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"swapFeeReward\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"WETH\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"swapMining\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"tradingPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountADesired\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBDesired\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountAMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"addLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenDesired\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"addLiquidityETH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"factory\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveOut\",\"type\":\"uint256\"}],\"name\":\"getAmountIn\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveOut\",\"type\":\"uint256\"}],\"name\":\"getAmountOut\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"}],\"name\":\"getAmountsIn\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"}],\"name\":\"getAmountsOut\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveB\",\"type\":\"uint256\"}],\"name\":\"quote\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountAMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"removeLiquidity\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"removeLiquidityETH\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"removeLiquidityETHSupportingFeeOnTransferTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"approveMax\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"removeLiquidityETHWithPermit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountToken\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountTokenMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountETHMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"approveMax\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"removeLiquidityETHWithPermitSupportingFeeOnTransferTokens\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountETH\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"tokenA\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"tokenB\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"liquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountAMin\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountBMin\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"approveMax\",\"type\":\"bool\"},{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"removeLiquidityWithPermit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountA\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountB\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapETHForExactTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactETHForTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactETHForTokensSupportingFeeOnTransferTokens\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForETH\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForHT\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForETHSupportingFeeOnTransferTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountIn\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactTokensForTokensSupportingFeeOnTransferTokens\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMax\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapTokensForExactETH\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOutMin\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapExactHTForTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amountOut\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amountInMax\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"path\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"swapTokensForExactTokens\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// Uniswapv2pairABI is the input ABI used to generate the binding from.
// Deprecated: Use Uniswapv2pairMetaData.ABI instead.
var Uniswapv2pairABI = Uniswapv2pairMetaData.ABI

// Uniswapv2pair is an auto generated Go binding around an Ethereum contract.
type Uniswapv2pair struct {
	Uniswapv2pairCaller     // Read-only binding to the contract
	Uniswapv2pairTransactor // Write-only binding to the contract
	Uniswapv2pairFilterer   // Log filterer for contract events
}

// Uniswapv2pairCaller is an auto generated read-only Go binding around an Ethereum contract.
type Uniswapv2pairCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Uniswapv2pairTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Uniswapv2pairTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Uniswapv2pairFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Uniswapv2pairFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Uniswapv2pairSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Uniswapv2pairSession struct {
	Contract     *Uniswapv2pair    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Uniswapv2pairCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Uniswapv2pairCallerSession struct {
	Contract *Uniswapv2pairCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// Uniswapv2pairTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Uniswapv2pairTransactorSession struct {
	Contract     *Uniswapv2pairTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// Uniswapv2pairRaw is an auto generated low-level Go binding around an Ethereum contract.
type Uniswapv2pairRaw struct {
	Contract *Uniswapv2pair // Generic contract binding to access the raw methods on
}

// Uniswapv2pairCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Uniswapv2pairCallerRaw struct {
	Contract *Uniswapv2pairCaller // Generic read-only contract binding to access the raw methods on
}

// Uniswapv2pairTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Uniswapv2pairTransactorRaw struct {
	Contract *Uniswapv2pairTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUniswapv2pair creates a new instance of Uniswapv2pair, bound to a specific deployed contract.
func NewUniswapv2pair(address common.Address, backend bind.ContractBackend) (*Uniswapv2pair, error) {
	contract, err := bindUniswapv2pair(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Uniswapv2pair{Uniswapv2pairCaller: Uniswapv2pairCaller{contract: contract}, Uniswapv2pairTransactor: Uniswapv2pairTransactor{contract: contract}, Uniswapv2pairFilterer: Uniswapv2pairFilterer{contract: contract}}, nil
}

// NewUniswapv2pairCaller creates a new read-only instance of Uniswapv2pair, bound to a specific deployed contract.
func NewUniswapv2pairCaller(address common.Address, caller bind.ContractCaller) (*Uniswapv2pairCaller, error) {
	contract, err := bindUniswapv2pair(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Uniswapv2pairCaller{contract: contract}, nil
}

// NewUniswapv2pairTransactor creates a new write-only instance of Uniswapv2pair, bound to a specific deployed contract.
func NewUniswapv2pairTransactor(address common.Address, transactor bind.ContractTransactor) (*Uniswapv2pairTransactor, error) {
	contract, err := bindUniswapv2pair(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Uniswapv2pairTransactor{contract: contract}, nil
}

// NewUniswapv2pairFilterer creates a new log filterer instance of Uniswapv2pair, bound to a specific deployed contract.
func NewUniswapv2pairFilterer(address common.Address, filterer bind.ContractFilterer) (*Uniswapv2pairFilterer, error) {
	contract, err := bindUniswapv2pair(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Uniswapv2pairFilterer{contract: contract}, nil
}

// bindUniswapv2pair binds a generic wrapper to an already deployed contract.
func bindUniswapv2pair(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Uniswapv2pairABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Uniswapv2pair *Uniswapv2pairRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Uniswapv2pair.Contract.Uniswapv2pairCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Uniswapv2pair *Uniswapv2pairRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.Uniswapv2pairTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Uniswapv2pair *Uniswapv2pairRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.Uniswapv2pairTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Uniswapv2pair *Uniswapv2pairCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Uniswapv2pair.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Uniswapv2pair *Uniswapv2pairTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Uniswapv2pair *Uniswapv2pairTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.contract.Transact(opts, method, params...)
}

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() pure returns(address)
func (_Uniswapv2pair *Uniswapv2pairCaller) WETH(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswapv2pair.contract.Call(opts, &out, "WETH")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() pure returns(address)
func (_Uniswapv2pair *Uniswapv2pairSession) WETH() (common.Address, error) {
	return _Uniswapv2pair.Contract.WETH(&_Uniswapv2pair.CallOpts)
}

// WETH is a free data retrieval call binding the contract method 0xad5c4648.
//
// Solidity: function WETH() pure returns(address)
func (_Uniswapv2pair *Uniswapv2pairCallerSession) WETH() (common.Address, error) {
	return _Uniswapv2pair.Contract.WETH(&_Uniswapv2pair.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() pure returns(address)
func (_Uniswapv2pair *Uniswapv2pairCaller) Factory(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswapv2pair.contract.Call(opts, &out, "factory")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() pure returns(address)
func (_Uniswapv2pair *Uniswapv2pairSession) Factory() (common.Address, error) {
	return _Uniswapv2pair.Contract.Factory(&_Uniswapv2pair.CallOpts)
}

// Factory is a free data retrieval call binding the contract method 0xc45a0155.
//
// Solidity: function factory() pure returns(address)
func (_Uniswapv2pair *Uniswapv2pairCallerSession) Factory() (common.Address, error) {
	return _Uniswapv2pair.Contract.Factory(&_Uniswapv2pair.CallOpts)
}

// GetAmountIn is a free data retrieval call binding the contract method 0x85f8c259.
//
// Solidity: function getAmountIn(uint256 amountOut, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountIn)
func (_Uniswapv2pair *Uniswapv2pairCaller) GetAmountIn(opts *bind.CallOpts, amountOut *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Uniswapv2pair.contract.Call(opts, &out, "getAmountIn", amountOut, reserveIn, reserveOut)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountIn is a free data retrieval call binding the contract method 0x85f8c259.
//
// Solidity: function getAmountIn(uint256 amountOut, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountIn)
func (_Uniswapv2pair *Uniswapv2pairSession) GetAmountIn(amountOut *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	return _Uniswapv2pair.Contract.GetAmountIn(&_Uniswapv2pair.CallOpts, amountOut, reserveIn, reserveOut)
}

// GetAmountIn is a free data retrieval call binding the contract method 0x85f8c259.
//
// Solidity: function getAmountIn(uint256 amountOut, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountIn)
func (_Uniswapv2pair *Uniswapv2pairCallerSession) GetAmountIn(amountOut *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	return _Uniswapv2pair.Contract.GetAmountIn(&_Uniswapv2pair.CallOpts, amountOut, reserveIn, reserveOut)
}

// GetAmountOut is a free data retrieval call binding the contract method 0x054d50d4.
//
// Solidity: function getAmountOut(uint256 amountIn, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountOut)
func (_Uniswapv2pair *Uniswapv2pairCaller) GetAmountOut(opts *bind.CallOpts, amountIn *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Uniswapv2pair.contract.Call(opts, &out, "getAmountOut", amountIn, reserveIn, reserveOut)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountOut is a free data retrieval call binding the contract method 0x054d50d4.
//
// Solidity: function getAmountOut(uint256 amountIn, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountOut)
func (_Uniswapv2pair *Uniswapv2pairSession) GetAmountOut(amountIn *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	return _Uniswapv2pair.Contract.GetAmountOut(&_Uniswapv2pair.CallOpts, amountIn, reserveIn, reserveOut)
}

// GetAmountOut is a free data retrieval call binding the contract method 0x054d50d4.
//
// Solidity: function getAmountOut(uint256 amountIn, uint256 reserveIn, uint256 reserveOut) pure returns(uint256 amountOut)
func (_Uniswapv2pair *Uniswapv2pairCallerSession) GetAmountOut(amountIn *big.Int, reserveIn *big.Int, reserveOut *big.Int) (*big.Int, error) {
	return _Uniswapv2pair.Contract.GetAmountOut(&_Uniswapv2pair.CallOpts, amountIn, reserveIn, reserveOut)
}

// GetAmountsIn is a free data retrieval call binding the contract method 0x1f00ca74.
//
// Solidity: function getAmountsIn(uint256 amountOut, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairCaller) GetAmountsIn(opts *bind.CallOpts, amountOut *big.Int, path []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _Uniswapv2pair.contract.Call(opts, &out, "getAmountsIn", amountOut, path)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAmountsIn is a free data retrieval call binding the contract method 0x1f00ca74.
//
// Solidity: function getAmountsIn(uint256 amountOut, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairSession) GetAmountsIn(amountOut *big.Int, path []common.Address) ([]*big.Int, error) {
	return _Uniswapv2pair.Contract.GetAmountsIn(&_Uniswapv2pair.CallOpts, amountOut, path)
}

// GetAmountsIn is a free data retrieval call binding the contract method 0x1f00ca74.
//
// Solidity: function getAmountsIn(uint256 amountOut, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairCallerSession) GetAmountsIn(amountOut *big.Int, path []common.Address) ([]*big.Int, error) {
	return _Uniswapv2pair.Contract.GetAmountsIn(&_Uniswapv2pair.CallOpts, amountOut, path)
}

// GetAmountsOut is a free data retrieval call binding the contract method 0xd06ca61f.
//
// Solidity: function getAmountsOut(uint256 amountIn, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairCaller) GetAmountsOut(opts *bind.CallOpts, amountIn *big.Int, path []common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _Uniswapv2pair.contract.Call(opts, &out, "getAmountsOut", amountIn, path)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetAmountsOut is a free data retrieval call binding the contract method 0xd06ca61f.
//
// Solidity: function getAmountsOut(uint256 amountIn, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairSession) GetAmountsOut(amountIn *big.Int, path []common.Address) ([]*big.Int, error) {
	return _Uniswapv2pair.Contract.GetAmountsOut(&_Uniswapv2pair.CallOpts, amountIn, path)
}

// GetAmountsOut is a free data retrieval call binding the contract method 0xd06ca61f.
//
// Solidity: function getAmountsOut(uint256 amountIn, address[] path) view returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairCallerSession) GetAmountsOut(amountIn *big.Int, path []common.Address) ([]*big.Int, error) {
	return _Uniswapv2pair.Contract.GetAmountsOut(&_Uniswapv2pair.CallOpts, amountIn, path)
}

// Quote is a free data retrieval call binding the contract method 0xad615dec.
//
// Solidity: function quote(uint256 amountA, uint256 reserveA, uint256 reserveB) pure returns(uint256 amountB)
func (_Uniswapv2pair *Uniswapv2pairCaller) Quote(opts *bind.CallOpts, amountA *big.Int, reserveA *big.Int, reserveB *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Uniswapv2pair.contract.Call(opts, &out, "quote", amountA, reserveA, reserveB)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Quote is a free data retrieval call binding the contract method 0xad615dec.
//
// Solidity: function quote(uint256 amountA, uint256 reserveA, uint256 reserveB) pure returns(uint256 amountB)
func (_Uniswapv2pair *Uniswapv2pairSession) Quote(amountA *big.Int, reserveA *big.Int, reserveB *big.Int) (*big.Int, error) {
	return _Uniswapv2pair.Contract.Quote(&_Uniswapv2pair.CallOpts, amountA, reserveA, reserveB)
}

// Quote is a free data retrieval call binding the contract method 0xad615dec.
//
// Solidity: function quote(uint256 amountA, uint256 reserveA, uint256 reserveB) pure returns(uint256 amountB)
func (_Uniswapv2pair *Uniswapv2pairCallerSession) Quote(amountA *big.Int, reserveA *big.Int, reserveB *big.Int) (*big.Int, error) {
	return _Uniswapv2pair.Contract.Quote(&_Uniswapv2pair.CallOpts, amountA, reserveA, reserveB)
}

// SwapFeeReward is a free data retrieval call binding the contract method 0xa3ddb30b.
//
// Solidity: function swapFeeReward() pure returns(address)
func (_Uniswapv2pair *Uniswapv2pairCaller) SwapFeeReward(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswapv2pair.contract.Call(opts, &out, "swapFeeReward")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SwapFeeReward is a free data retrieval call binding the contract method 0xa3ddb30b.
//
// Solidity: function swapFeeReward() pure returns(address)
func (_Uniswapv2pair *Uniswapv2pairSession) SwapFeeReward() (common.Address, error) {
	return _Uniswapv2pair.Contract.SwapFeeReward(&_Uniswapv2pair.CallOpts)
}

// SwapFeeReward is a free data retrieval call binding the contract method 0xa3ddb30b.
//
// Solidity: function swapFeeReward() pure returns(address)
func (_Uniswapv2pair *Uniswapv2pairCallerSession) SwapFeeReward() (common.Address, error) {
	return _Uniswapv2pair.Contract.SwapFeeReward(&_Uniswapv2pair.CallOpts)
}

// SwapMining is a free data retrieval call binding the contract method 0x4e5eda11.
//
// Solidity: function swapMining() view returns(address)
func (_Uniswapv2pair *Uniswapv2pairCaller) SwapMining(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswapv2pair.contract.Call(opts, &out, "swapMining")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SwapMining is a free data retrieval call binding the contract method 0x4e5eda11.
//
// Solidity: function swapMining() view returns(address)
func (_Uniswapv2pair *Uniswapv2pairSession) SwapMining() (common.Address, error) {
	return _Uniswapv2pair.Contract.SwapMining(&_Uniswapv2pair.CallOpts)
}

// SwapMining is a free data retrieval call binding the contract method 0x4e5eda11.
//
// Solidity: function swapMining() view returns(address)
func (_Uniswapv2pair *Uniswapv2pairCallerSession) SwapMining() (common.Address, error) {
	return _Uniswapv2pair.Contract.SwapMining(&_Uniswapv2pair.CallOpts)
}

// TradingPool is a free data retrieval call binding the contract method 0x71f16a18.
//
// Solidity: function tradingPool() view returns(address)
func (_Uniswapv2pair *Uniswapv2pairCaller) TradingPool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Uniswapv2pair.contract.Call(opts, &out, "tradingPool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TradingPool is a free data retrieval call binding the contract method 0x71f16a18.
//
// Solidity: function tradingPool() view returns(address)
func (_Uniswapv2pair *Uniswapv2pairSession) TradingPool() (common.Address, error) {
	return _Uniswapv2pair.Contract.TradingPool(&_Uniswapv2pair.CallOpts)
}

// TradingPool is a free data retrieval call binding the contract method 0x71f16a18.
//
// Solidity: function tradingPool() view returns(address)
func (_Uniswapv2pair *Uniswapv2pairCallerSession) TradingPool() (common.Address, error) {
	return _Uniswapv2pair.Contract.TradingPool(&_Uniswapv2pair.CallOpts)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0xe8e33700.
//
// Solidity: function addLiquidity(address tokenA, address tokenB, uint256 amountADesired, uint256 amountBDesired, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB, uint256 liquidity)
func (_Uniswapv2pair *Uniswapv2pairTransactor) AddLiquidity(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address, amountADesired *big.Int, amountBDesired *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "addLiquidity", tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, to, deadline)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0xe8e33700.
//
// Solidity: function addLiquidity(address tokenA, address tokenB, uint256 amountADesired, uint256 amountBDesired, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB, uint256 liquidity)
func (_Uniswapv2pair *Uniswapv2pairSession) AddLiquidity(tokenA common.Address, tokenB common.Address, amountADesired *big.Int, amountBDesired *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.AddLiquidity(&_Uniswapv2pair.TransactOpts, tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, to, deadline)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0xe8e33700.
//
// Solidity: function addLiquidity(address tokenA, address tokenB, uint256 amountADesired, uint256 amountBDesired, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB, uint256 liquidity)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) AddLiquidity(tokenA common.Address, tokenB common.Address, amountADesired *big.Int, amountBDesired *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.AddLiquidity(&_Uniswapv2pair.TransactOpts, tokenA, tokenB, amountADesired, amountBDesired, amountAMin, amountBMin, to, deadline)
}

// AddLiquidityETH is a paid mutator transaction binding the contract method 0xf305d719.
//
// Solidity: function addLiquidityETH(address token, uint256 amountTokenDesired, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) payable returns(uint256 amountToken, uint256 amountETH, uint256 liquidity)
func (_Uniswapv2pair *Uniswapv2pairTransactor) AddLiquidityETH(opts *bind.TransactOpts, token common.Address, amountTokenDesired *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "addLiquidityETH", token, amountTokenDesired, amountTokenMin, amountETHMin, to, deadline)
}

// AddLiquidityETH is a paid mutator transaction binding the contract method 0xf305d719.
//
// Solidity: function addLiquidityETH(address token, uint256 amountTokenDesired, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) payable returns(uint256 amountToken, uint256 amountETH, uint256 liquidity)
func (_Uniswapv2pair *Uniswapv2pairSession) AddLiquidityETH(token common.Address, amountTokenDesired *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.AddLiquidityETH(&_Uniswapv2pair.TransactOpts, token, amountTokenDesired, amountTokenMin, amountETHMin, to, deadline)
}

// AddLiquidityETH is a paid mutator transaction binding the contract method 0xf305d719.
//
// Solidity: function addLiquidityETH(address token, uint256 amountTokenDesired, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) payable returns(uint256 amountToken, uint256 amountETH, uint256 liquidity)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) AddLiquidityETH(token common.Address, amountTokenDesired *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.AddLiquidityETH(&_Uniswapv2pair.TransactOpts, token, amountTokenDesired, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xbaa2abde.
//
// Solidity: function removeLiquidity(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2pair *Uniswapv2pairTransactor) RemoveLiquidity(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "removeLiquidity", tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xbaa2abde.
//
// Solidity: function removeLiquidity(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2pair *Uniswapv2pairSession) RemoveLiquidity(tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidity(&_Uniswapv2pair.TransactOpts, tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xbaa2abde.
//
// Solidity: function removeLiquidity(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) RemoveLiquidity(tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidity(&_Uniswapv2pair.TransactOpts, tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline)
}

// RemoveLiquidityETH is a paid mutator transaction binding the contract method 0x02751cec.
//
// Solidity: function removeLiquidityETH(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairTransactor) RemoveLiquidityETH(opts *bind.TransactOpts, token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "removeLiquidityETH", token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETH is a paid mutator transaction binding the contract method 0x02751cec.
//
// Solidity: function removeLiquidityETH(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairSession) RemoveLiquidityETH(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidityETH(&_Uniswapv2pair.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETH is a paid mutator transaction binding the contract method 0x02751cec.
//
// Solidity: function removeLiquidityETH(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) RemoveLiquidityETH(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidityETH(&_Uniswapv2pair.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xaf2979eb.
//
// Solidity: function removeLiquidityETHSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairTransactor) RemoveLiquidityETHSupportingFeeOnTransferTokens(opts *bind.TransactOpts, token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "removeLiquidityETHSupportingFeeOnTransferTokens", token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xaf2979eb.
//
// Solidity: function removeLiquidityETHSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairSession) RemoveLiquidityETHSupportingFeeOnTransferTokens(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidityETHSupportingFeeOnTransferTokens(&_Uniswapv2pair.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xaf2979eb.
//
// Solidity: function removeLiquidityETHSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline) returns(uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) RemoveLiquidityETHSupportingFeeOnTransferTokens(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidityETHSupportingFeeOnTransferTokens(&_Uniswapv2pair.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline)
}

// RemoveLiquidityETHWithPermit is a paid mutator transaction binding the contract method 0xded9382a.
//
// Solidity: function removeLiquidityETHWithPermit(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairTransactor) RemoveLiquidityETHWithPermit(opts *bind.TransactOpts, token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "removeLiquidityETHWithPermit", token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityETHWithPermit is a paid mutator transaction binding the contract method 0xded9382a.
//
// Solidity: function removeLiquidityETHWithPermit(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairSession) RemoveLiquidityETHWithPermit(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidityETHWithPermit(&_Uniswapv2pair.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityETHWithPermit is a paid mutator transaction binding the contract method 0xded9382a.
//
// Solidity: function removeLiquidityETHWithPermit(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountToken, uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) RemoveLiquidityETHWithPermit(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidityETHWithPermit(&_Uniswapv2pair.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5b0d5984.
//
// Solidity: function removeLiquidityETHWithPermitSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairTransactor) RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(opts *bind.TransactOpts, token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "removeLiquidityETHWithPermitSupportingFeeOnTransferTokens", token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5b0d5984.
//
// Solidity: function removeLiquidityETHWithPermitSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairSession) RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(&_Uniswapv2pair.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5b0d5984.
//
// Solidity: function removeLiquidityETHWithPermitSupportingFeeOnTransferTokens(address token, uint256 liquidity, uint256 amountTokenMin, uint256 amountETHMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountETH)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(token common.Address, liquidity *big.Int, amountTokenMin *big.Int, amountETHMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidityETHWithPermitSupportingFeeOnTransferTokens(&_Uniswapv2pair.TransactOpts, token, liquidity, amountTokenMin, amountETHMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityWithPermit is a paid mutator transaction binding the contract method 0x2195995c.
//
// Solidity: function removeLiquidityWithPermit(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2pair *Uniswapv2pairTransactor) RemoveLiquidityWithPermit(opts *bind.TransactOpts, tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "removeLiquidityWithPermit", tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityWithPermit is a paid mutator transaction binding the contract method 0x2195995c.
//
// Solidity: function removeLiquidityWithPermit(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2pair *Uniswapv2pairSession) RemoveLiquidityWithPermit(tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidityWithPermit(&_Uniswapv2pair.TransactOpts, tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline, approveMax, v, r, s)
}

// RemoveLiquidityWithPermit is a paid mutator transaction binding the contract method 0x2195995c.
//
// Solidity: function removeLiquidityWithPermit(address tokenA, address tokenB, uint256 liquidity, uint256 amountAMin, uint256 amountBMin, address to, uint256 deadline, bool approveMax, uint8 v, bytes32 r, bytes32 s) returns(uint256 amountA, uint256 amountB)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) RemoveLiquidityWithPermit(tokenA common.Address, tokenB common.Address, liquidity *big.Int, amountAMin *big.Int, amountBMin *big.Int, to common.Address, deadline *big.Int, approveMax bool, v uint8, r [32]byte, s [32]byte) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.RemoveLiquidityWithPermit(&_Uniswapv2pair.TransactOpts, tokenA, tokenB, liquidity, amountAMin, amountBMin, to, deadline, approveMax, v, r, s)
}

// SwapETHForExactTokens is a paid mutator transaction binding the contract method 0xfb3bdb41.
//
// Solidity: function swapETHForExactTokens(uint256 amountOut, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactor) SwapETHForExactTokens(opts *bind.TransactOpts, amountOut *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "swapETHForExactTokens", amountOut, path, to, deadline)
}

// SwapETHForExactTokens is a paid mutator transaction binding the contract method 0xfb3bdb41.
//
// Solidity: function swapETHForExactTokens(uint256 amountOut, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairSession) SwapETHForExactTokens(amountOut *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapETHForExactTokens(&_Uniswapv2pair.TransactOpts, amountOut, path, to, deadline)
}

// SwapETHForExactTokens is a paid mutator transaction binding the contract method 0xfb3bdb41.
//
// Solidity: function swapETHForExactTokens(uint256 amountOut, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) SwapETHForExactTokens(amountOut *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapETHForExactTokens(&_Uniswapv2pair.TransactOpts, amountOut, path, to, deadline)
}

// SwapExactETHForTokens is a paid mutator transaction binding the contract method 0x7ff36ab5.
//
// Solidity: function swapExactETHForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactor) SwapExactETHForTokens(opts *bind.TransactOpts, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "swapExactETHForTokens", amountOutMin, path, to, deadline)
}

// SwapExactETHForTokens is a paid mutator transaction binding the contract method 0x7ff36ab5.
//
// Solidity: function swapExactETHForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairSession) SwapExactETHForTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactETHForTokens(&_Uniswapv2pair.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactETHForTokens is a paid mutator transaction binding the contract method 0x7ff36ab5.
//
// Solidity: function swapExactETHForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) SwapExactETHForTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactETHForTokens(&_Uniswapv2pair.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactETHForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xb6f9de95.
//
// Solidity: function swapExactETHForTokensSupportingFeeOnTransferTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns()
func (_Uniswapv2pair *Uniswapv2pairTransactor) SwapExactETHForTokensSupportingFeeOnTransferTokens(opts *bind.TransactOpts, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "swapExactETHForTokensSupportingFeeOnTransferTokens", amountOutMin, path, to, deadline)
}

// SwapExactETHForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xb6f9de95.
//
// Solidity: function swapExactETHForTokensSupportingFeeOnTransferTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns()
func (_Uniswapv2pair *Uniswapv2pairSession) SwapExactETHForTokensSupportingFeeOnTransferTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactETHForTokensSupportingFeeOnTransferTokens(&_Uniswapv2pair.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactETHForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0xb6f9de95.
//
// Solidity: function swapExactETHForTokensSupportingFeeOnTransferTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns()
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) SwapExactETHForTokensSupportingFeeOnTransferTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactETHForTokensSupportingFeeOnTransferTokens(&_Uniswapv2pair.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactHTForTokens is a paid mutator transaction binding the contract method 0xde6b8256.
//
// Solidity: function swapExactHTForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactor) SwapExactHTForTokens(opts *bind.TransactOpts, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "swapExactHTForTokens", amountOutMin, path, to, deadline)
}

// SwapExactHTForTokens is a paid mutator transaction binding the contract method 0xde6b8256.
//
// Solidity: function swapExactHTForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairSession) SwapExactHTForTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactHTForTokens(&_Uniswapv2pair.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactHTForTokens is a paid mutator transaction binding the contract method 0xde6b8256.
//
// Solidity: function swapExactHTForTokens(uint256 amountOutMin, address[] path, address to, uint256 deadline) payable returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) SwapExactHTForTokens(amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactHTForTokens(&_Uniswapv2pair.TransactOpts, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETH is a paid mutator transaction binding the contract method 0x18cbafe5.
//
// Solidity: function swapExactTokensForETH(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactor) SwapExactTokensForETH(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "swapExactTokensForETH", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETH is a paid mutator transaction binding the contract method 0x18cbafe5.
//
// Solidity: function swapExactTokensForETH(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairSession) SwapExactTokensForETH(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactTokensForETH(&_Uniswapv2pair.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETH is a paid mutator transaction binding the contract method 0x18cbafe5.
//
// Solidity: function swapExactTokensForETH(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) SwapExactTokensForETH(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactTokensForETH(&_Uniswapv2pair.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x791ac947.
//
// Solidity: function swapExactTokensForETHSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2pair *Uniswapv2pairTransactor) SwapExactTokensForETHSupportingFeeOnTransferTokens(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "swapExactTokensForETHSupportingFeeOnTransferTokens", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x791ac947.
//
// Solidity: function swapExactTokensForETHSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2pair *Uniswapv2pairSession) SwapExactTokensForETHSupportingFeeOnTransferTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactTokensForETHSupportingFeeOnTransferTokens(&_Uniswapv2pair.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForETHSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x791ac947.
//
// Solidity: function swapExactTokensForETHSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) SwapExactTokensForETHSupportingFeeOnTransferTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactTokensForETHSupportingFeeOnTransferTokens(&_Uniswapv2pair.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForHT is a paid mutator transaction binding the contract method 0x4b361208.
//
// Solidity: function swapExactTokensForHT(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactor) SwapExactTokensForHT(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "swapExactTokensForHT", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForHT is a paid mutator transaction binding the contract method 0x4b361208.
//
// Solidity: function swapExactTokensForHT(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairSession) SwapExactTokensForHT(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactTokensForHT(&_Uniswapv2pair.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForHT is a paid mutator transaction binding the contract method 0x4b361208.
//
// Solidity: function swapExactTokensForHT(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) SwapExactTokensForHT(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactTokensForHT(&_Uniswapv2pair.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0x38ed1739.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactor) SwapExactTokensForTokens(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "swapExactTokensForTokens", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0x38ed1739.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairSession) SwapExactTokensForTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactTokensForTokens(&_Uniswapv2pair.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokens is a paid mutator transaction binding the contract method 0x38ed1739.
//
// Solidity: function swapExactTokensForTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) SwapExactTokensForTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactTokensForTokens(&_Uniswapv2pair.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5c11d795.
//
// Solidity: function swapExactTokensForTokensSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2pair *Uniswapv2pairTransactor) SwapExactTokensForTokensSupportingFeeOnTransferTokens(opts *bind.TransactOpts, amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "swapExactTokensForTokensSupportingFeeOnTransferTokens", amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5c11d795.
//
// Solidity: function swapExactTokensForTokensSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2pair *Uniswapv2pairSession) SwapExactTokensForTokensSupportingFeeOnTransferTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactTokensForTokensSupportingFeeOnTransferTokens(&_Uniswapv2pair.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapExactTokensForTokensSupportingFeeOnTransferTokens is a paid mutator transaction binding the contract method 0x5c11d795.
//
// Solidity: function swapExactTokensForTokensSupportingFeeOnTransferTokens(uint256 amountIn, uint256 amountOutMin, address[] path, address to, uint256 deadline) returns()
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) SwapExactTokensForTokensSupportingFeeOnTransferTokens(amountIn *big.Int, amountOutMin *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapExactTokensForTokensSupportingFeeOnTransferTokens(&_Uniswapv2pair.TransactOpts, amountIn, amountOutMin, path, to, deadline)
}

// SwapTokensForExactETH is a paid mutator transaction binding the contract method 0x4a25d94a.
//
// Solidity: function swapTokensForExactETH(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactor) SwapTokensForExactETH(opts *bind.TransactOpts, amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "swapTokensForExactETH", amountOut, amountInMax, path, to, deadline)
}

// SwapTokensForExactETH is a paid mutator transaction binding the contract method 0x4a25d94a.
//
// Solidity: function swapTokensForExactETH(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairSession) SwapTokensForExactETH(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapTokensForExactETH(&_Uniswapv2pair.TransactOpts, amountOut, amountInMax, path, to, deadline)
}

// SwapTokensForExactETH is a paid mutator transaction binding the contract method 0x4a25d94a.
//
// Solidity: function swapTokensForExactETH(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) SwapTokensForExactETH(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapTokensForExactETH(&_Uniswapv2pair.TransactOpts, amountOut, amountInMax, path, to, deadline)
}

// SwapTokensForExactTokens is a paid mutator transaction binding the contract method 0x8803dbee.
//
// Solidity: function swapTokensForExactTokens(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactor) SwapTokensForExactTokens(opts *bind.TransactOpts, amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.contract.Transact(opts, "swapTokensForExactTokens", amountOut, amountInMax, path, to, deadline)
}

// SwapTokensForExactTokens is a paid mutator transaction binding the contract method 0x8803dbee.
//
// Solidity: function swapTokensForExactTokens(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairSession) SwapTokensForExactTokens(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapTokensForExactTokens(&_Uniswapv2pair.TransactOpts, amountOut, amountInMax, path, to, deadline)
}

// SwapTokensForExactTokens is a paid mutator transaction binding the contract method 0x8803dbee.
//
// Solidity: function swapTokensForExactTokens(uint256 amountOut, uint256 amountInMax, address[] path, address to, uint256 deadline) returns(uint256[] amounts)
func (_Uniswapv2pair *Uniswapv2pairTransactorSession) SwapTokensForExactTokens(amountOut *big.Int, amountInMax *big.Int, path []common.Address, to common.Address, deadline *big.Int) (*types.Transaction, error) {
	return _Uniswapv2pair.Contract.SwapTokensForExactTokens(&_Uniswapv2pair.TransactOpts, amountOut, amountInMax, path, to, deadline)
}
