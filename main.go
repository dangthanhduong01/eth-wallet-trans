package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	url  = "https://bsc-testnet.blockpi.network/v1/rpc/public"
	murl = ""
)

func createWallet() (string, string) {
	pvk, err := crypto.GenerateKey()
	if err != nil {
		log.Println(err)
	}
	getPublicKey := crypto.FromECDSA(pvk)
	thePublicKey := hexutil.Encode(getPublicKey)

	thePublicAddress := crypto.PubkeyToAddress(pvk.PublicKey).Hex()

	return thePublicAddress, thePublicKey
}

func main() {
	// connecting to ethereum node
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Read the Account Balance
	// DuongMobie 0x34d9cd385f4cb710f9608620476933a5e9706bec
	account1 := common.HexToAddress("0x246BbaC3b0Cb5e7bF137880EB771AA9c87cD5698")
	// DuongPC 0x2376b1808817b1833BfA97532118383d67974D09
	account2 := common.HexToAddress("0x2376b1808817b1833BfA97532118383d67974D09")

	//Check the Balance of client instance
	balance1, err := client.BalanceAt(context.Background(), account1, nil)
	if err != nil {
		log.Fatal(err)
	}

	balance2, err := client.BalanceAt(context.Background(), account2, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance1)
	fmt.Println(balance2)

	// header, err := client.HeaderByNumber(context.Background(), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("BlockNumber:", header.Number.String())
	// }
	// nblock := header.Number.String()
	// bn, err := strconv.ParseInt(nblock, 10, 64)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// Passing the block number let's you read the account balance at the time of that block.
	// blockNumber := big.NewInt(bn)
	// balanceAt, err := client.BalanceAt(context.Background(), account1, blockNumber)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(balanceAt) // 25729324269165216042

	// fbalance := new(big.Float)
	// fbalance.SetString(balanceAt.String())
	// ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	// fmt.Println(ethValue) // 25.729324269165216041

	// pendingBalance, err := client.PendingBalanceAt(context.Background(), account1)
	// fmt.Println(pendingBalance) // 25729324269165216042

	// Số nonce của giao dịch. Mỗi giao dịch từ cùng một tài khoản phải có số nonce unique.
	nonce, err := client.PendingNonceAt(context.Background(), account1)
	if err != nil {
		log.Fatal(err)
	}

	amount := big.NewInt(100000000000000000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTransaction(nonce, account2, amount, 21000, gasPrice, nil)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := crypto.HexToECDSA("387030e12125537380d3a5f8365c03f2a6443aecbca4044cf115ecf7a6bb4c63")
	if err != nil {
		log.Fatal(err)
	}
	// b, err := ioutil.ReadFile("wallet/UTC--2024-04-04T04-55-51.779251019Z--d825c58758fab07cbed3791c67ec52a89d57fe8c")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// key, err := keystore.DecryptKey(b, "password")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	tx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s", tx)
}
