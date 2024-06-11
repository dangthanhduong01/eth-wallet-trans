package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// Connect to RPC
	client, err := ethclient.Dial("https://rpc.ankr.com/eth")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network: %v", err)
		return
	}
	fmt.Println("Success! you are connected to the Ethereum Network")
	fmt.Printf("EthClient: %+v\n", client)
	fmt.Println("------------------------------------------------------")

	// Call balanceOf() of USDT contract over ABI
	ABI := "[{\"constant\":true,\"inputs\":[{\"name\":\"who\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"
	tokenAbi, err := abi.JSON(strings.NewReader(ABI))
	data, err := tokenAbi.Pack("balanceOf", common.HexToAddress("0xF977814e90dA44bFA03b6295A0616a897441aceC"))
	tokenAddr := common.HexToAddress("0xdac17f958d2ee523a2206206994597c13d831ec7")
	msg := ethereum.CallMsg{
		To:   &tokenAddr,
		Data: data,
	}
	outputData, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Call balanceOf() ERROR: %v", err)
		return
	}
	balance, err := tokenAbi.Unpack("balanceOf", outputData)
	if err != nil {
		log.Fatalf("Decode output of balanceOf() ERROR: %v", err)
		return
	}
	fmt.Println("USDT:", balance[0])
	fmt.Println("------------------------------------------------------")

	// Call getReserves() using NewMethod
	funcName := "getReserves"
	var inputs abi.Arguments = nil
	uint112Type, err := abi.NewType("uint112", "uint112", nil)
	outputs := abi.Arguments{
		abi.Argument{
			Name:    "reserve0",
			Type:    uint112Type,
			Indexed: false,
		},
		abi.Argument{
			Name:    "reserve1",
			Type:    uint112Type,
			Indexed: false,
		},
	}
	method := abi.NewMethod(funcName, funcName, abi.FunctionType(3), "view", false, false, inputs, outputs)
	sig := crypto.Keccak256([]byte(method.Sig))
	sig = []byte{sig[0], sig[1], sig[2], sig[3]}
	fmt.Println("Method:", method.String(), " => Sig:", "0x"+hex.EncodeToString(sig))
	data = sig
	pairV2Addr := common.HexToAddress("0xB4e16d0168e52d35CaCD2c6185b44281Ec28C9Dc")
	msg = ethereum.CallMsg{
		To:   &pairV2Addr,
		Data: data,
	}
	outputData, err = client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatalf("Call getReserves() ERROR: %v", err)
		return
	}
	reserves, err := outputs.Unpack(outputData)
	if err != nil {
		log.Fatalf("Decode output of getReserves() ERROR: %v", err)
		return
	}
	fmt.Println("Reserves:", reserves)
	fmt.Println("------------------------------------------------------")
}
