package main

import (
	"bytes"
	abi2 "clients/contract"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"io/ioutil"
	"math/big"
	"os"
)

var (
	FromAddr    = common.HexToAddress("")
	ToAddr      = common.HexToAddress("")
	ContractCdd = ""
	BscTestNet  = "https://data-seed-prebsc-2-s2.binance.org:8545"
	priva1, _   = ecdsa.GenerateKey(nil, nil)
)

/*
transfer(address, uint256)
balanceOf(address)
decimals()
allowance(address, address)
symbol()
totalSupply()
name()
approve(address, uint256)
transferFrom(address, address, uint256)
*/

func callContract(client *ethclient.Client, privKey *ecdsa.PrivateKey, from, to common.Address, contract string) (string, error) {
	// 创建交易
	nonce, err := client.NonceAt(context.Background(), from, nil)
	if err != nil {
		return "", err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	var data []byte
	method := []byte("transfer(address, uint256)")
	method = crypto.Keccak256(method)[:4]
	paddedAddress := common.LeftPadBytes(to.Bytes(), 32)
	amount, _ := new(big.Int).SetString(fmt.Sprint(10^18), 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	data = append(data, method...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)
	tx := types.NewTransaction(nonce, common.HexToAddress(contract), big.NewInt(0), uint64(300000), gasPrice, data)
	// 签名交易
	signed, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(97)), priva1)
	if err != nil {
		return "", err
	}
	// 发送交易
	err = client.SendTransaction(context.Background(), signed)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

func callContractByABI(client *ethclient.Client, privKey *ecdsa.PrivateKey, from, to common.Address, contract string) (string, error) {
	// 创建交易
	nonce, err := client.NonceAt(context.Background(), from, nil)
	if err != nil {
		return "", err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	abiData, err := ioutil.ReadFile("contract.txt")
	if err != nil {
		return "", err
	}
	contractABI, err := abi.JSON(bytes.NewReader(abiData))
	if err != nil {
		return "", err
	}
	amount, _ := new(big.Int).SetString(fmt.Sprint(10^18), 10)

	callData, err := contractABI.Pack("transfer", to, amount)
	if err != nil {
		return "", err
	}
	tx := types.NewTransaction(nonce, common.HexToAddress(contract), big.NewInt(0), uint64(300000), gasPrice, callData)
	// 签名交易
	signed, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(97)), priva1)
	if err != nil {
		return "", err
	}
	// 发送交易
	err = client.SendTransaction(context.Background(), signed)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

func callTotalSupplyByABIStruct(client *ethclient.Client) {
	cdd, err := abi2.NewErc20Token(common.HexToAddress(ContractCdd), client)
	if err != nil {
		fmt.Print("new token err", err)
		os.Exit(1)
	}
	total, err := cdd.TotalSupply(nil)

	if err != nil {
		fmt.Print("get total supply err", err)
		os.Exit(1)
	}
	fmt.Println(total.String())
}
func callContractByABIStruct(client *ethclient.Client, privKey *ecdsa.PrivateKey, from, to common.Address, contract string) (string, error) {
	// 创建交易
	nonce, err := client.NonceAt(context.Background(), from, nil)
	if err != nil {
		return "", err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}
	opts, err := bind.NewKeyedTransactorWithChainID(priva1, big.NewInt(97))
	if err != nil {
		return "", err
	}
	amount, _ := new(big.Int).SetString("5000000", 10)
	opts.GasLimit = uint64(300000)
	opts.GasPrice = gasPrice
	opts.Nonce = new(big.Int).SetUint64(nonce)
	cdd, err := abi2.NewErc20Token(common.HexToAddress(ContractCdd), client)
	tx, err := cdd.Transfer(opts, ToAddr, amount)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

func main() {
	client, err := ethclient.Dial(BscTestNet)
	if err != nil {
		fmt.Print("dial client err", err)
		os.Exit(1)
	}
	tx, err := callContract(client, priva1, FromAddr, ToAddr, ContractCdd)
	if err != nil {
		fmt.Println("call contract err", err)
		os.Exit(1)
	}
	fmt.Println("tx hash:", tx)
}
