package main

import (
	"log"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

func main() {
	// fmt.Println(keystore.StandardScryptN, keystore.StandardScryptP)
	//Create new KeyStore
	ks := keystore.NewKeyStore("./wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	for i := 0; i < 2; i++ {
		_, err := ks.NewAccount("password")
		if err != nil {
			log.Fatal(err)
		}
	}

}
