package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {

	password := "password"

	b, err := ioutil.ReadFile("./wallet/UTC--2024-04-04T04-55-53.138917245Z--810c6788bc77dca79559100403c0375b0e74f452")
	if err != nil {
		log.Fatal(err)
	}
	key, err := keystore.DecryptKey(b, password)
	if err != nil {
		log.Fatal(err)
	}
	pData := crypto.FromECDSA(key.PrivateKey)
	fmt.Println("Priv", hexutil.Encode(pData))

	pData = crypto.FromECDSAPub(&key.PrivateKey.PublicKey)
	fmt.Println("Pub", hexutil.Encode(pData))

	fmt.Println("Add", crypto.PubkeyToAddress(key.PrivateKey.PublicKey).Hex())
	account := common.HexToAddress("0x246BbaC3b0Cb5e7bF137880EB771AA9c87cD5698")
	fmt.Println(account)
}
