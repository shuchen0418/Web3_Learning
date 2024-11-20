package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
	"log"
)

func main() {

	/*privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	//fmt.Println(hexutil.Encode(privateKeyBytes))
	fmt.Println(hexutil.Encode(privateKeyBytes)[2:]) //ef051aa0bad1c865c8eb1cb62a387c3ad26d6c0ee1d466200d48b5893e2aece5
	*/
	ecdsaPrivateKey, err := crypto.HexToECDSA("ef051aa0bad1c865c8eb1cb62a387c3ad26d6c0ee1d466200d48b5893e2aece5")
	if err != nil {
		log.Fatal(err)
	}
	ecdsaPrivateKeyBytes := crypto.FromECDSA(ecdsaPrivateKey)
	fmt.Println(hexutil.Encode(ecdsaPrivateKeyBytes)[2:])

	publicKey := ecdsaPrivateKey.Public()

	publicKeyBytes := crypto.FromECDSAPub(publicKey.(*ecdsa.PublicKey))

	fmt.Println(hexutil.Encode(publicKeyBytes)[4:])

	address := crypto.PubkeyToAddress(*publicKey.(*ecdsa.PublicKey)).Hex()

	fmt.Println(address)

	hash := sha3.NewLegacyKeccak256()
	h := sha3.NewLegacyKeccak256()
	h.Write(publicKeyBytes[1:])
	hash.Write(publicKeyBytes)
	fmt.Println("[1:]", hexutil.Encode(hash.Sum(nil)[12:]))
	fmt.Println(hexutil.Encode(h.Sum(nil)[12:]))
}
