package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"regexp"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

var (
	url = "https://bsc-testnet.blockpi.network/v1/rpc/public"
	re  = regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
)

func checkValidAdress(address string) bool {
	return re.MatchString(address)
}

func main() {
	// connecting to ethereum node
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var address string
	fmt.Println("Nhập địa chỉ ví của bạn")
	fmt.Scanln(&address)
	for !checkValidAdress(address) {
		fmt.Println("Địa chỉ ví không đúng")
		fmt.Scanln(&address)
	}
	youraccount := common.HexToAddress(address)

	// Get header
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Phiên giao dịch của ban (Block Number):", header.Number.String())
	}
	balance, err := client.BalanceAt(context.Background(), youraccount, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("1: Xem số dư")
	fmt.Println("2: Chuyển tiền")
	fmt.Println("3: Xem block Number")
	fmt.Println("4. Kiểm tra wallet")
	fmt.Println("5: Swap Ethreum value to wei")
	fmt.Println("6 Giao dịch ERC-20")
	fmt.Println("")

	var input int
	fmt.Scanln(&input)

	switch input {
	case 1:
		fmt.Println(balance)
	case 2:
		var input2 string
		fmt.Println("Nhập địa chỉ ví cần chuyển")
		fmt.Scanln(&input2)
		for !checkValidAdress(input2) {
			log.Fatal("Địa chỉ ví không đúng")
			fmt.Scanln(&input2)
		}
		account2 := common.HexToAddress(input2)
		var amount int64
		fmt.Println("Nhập số tiền cần chuyển")

		fmt.Scanln(&amount)
		amountsend := big.NewInt(amount)
		if amountsend.Cmp(balance) >= 0 {
			log.Fatal("Không đủ số dư")
		}
		nonce, err := client.PendingNonceAt(context.Background(), youraccount)
		if err != nil {
			log.Fatal(err)
		}
		gasPrice, err := client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		tx := types.NewTransaction(nonce, account2, amountsend, 21000, gasPrice, nil)
		chainID, err := client.NetworkID(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		var pvkey string
		fmt.Println("Nhập private key")
		fmt.Scanln(&pvkey)

		privateKey, err := crypto.HexToECDSA(pvkey)
		if err != nil {
			log.Fatal(err)
		}
		tx, err = types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
		if err != nil {
			log.Fatal(err)
		}
		err = client.SendTransaction(context.Background(), tx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("tx sent: %s \n", tx.Hash().Hex())
		//check Transaction success or false
		receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(receipt.Status)
		fmt.Println(receipt.Logs)

	case 3:
		fmt.Println("Phiên giao dịch của ban (Block Number):", header.Number.String())
	case 4:
		mnemonic := ""
		fmt.Println("Nhập cụm từ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		mnemonic = scanner.Text()
		wallet, err := hdwallet.NewFromMnemonic(mnemonic)
		if err != nil {
			log.Fatal(err)
		}
		path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
		acc, err := wallet.Derive(path, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(acc.Address.Hex())
	case 5:
		input5 := ""
		fmt.Println("Nhập giá trị Ethreum")
		fmt.Scanln(&input5)
		fbalance := new(big.Float)
		fbalance.SetString(input5)
		weiValue := new(big.Float).Mul(fbalance, big.NewFloat(math.Pow10(18)))
		v, _ := weiValue.Int64()

		fmt.Printf("Giá trị Wei %v \n", v)

	default:
		fmt.Println("Invalid")
	}

}
