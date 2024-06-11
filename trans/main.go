package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	url  = "https://bsc-testnet.blockpi.network/v1/rpc/public"
	murl = ""
)

func main() {

	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	RecipientAddress := common.HexToAddress("0x34d9cd385f4cb710f9608620476933a5e9706bec")

	privateKey, err := crypto.HexToECDSA("0x34d9cd385f4cb710f9608620476933a5e9706bec")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Public Key Error")
	}

	SenderAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), SenderAddress)
	if err != nil {

		log.Println(err)
	}

	amount := big.NewInt(1000000000000000000)
	gasLimit := 3600
	gas, err := client.SuggestGasPrice(context.Background())

	if err != nil {
		log.Println(err)
	}

	ChainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Println(err)
	}

	transaction := types.NewTransaction(nonce, RecipientAddress, amount, uint64(gasLimit), gas, nil)
	signedTx, err := types.SignTx(transaction, types.NewEIP155Signer(ChainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("transaction sent: %s", signedTx.Hash().Hex())

}
