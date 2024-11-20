package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/ZYOdzRs9UVWVsb715Ctarx9O0CxvGs-f")
	if err != nil {
		log.Fatal(err)
	}

	privateKeyEcdsa, err := crypto.HexToECDSA("f4db97a3f3fd46f7b16c8007ff2d26a7e9f755a69a947cd54124ade8d0d5e2dd")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKeyEcdsa.Public()

	publicKeyEcdsa, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("cannot assert type: publicKey is of type %T", publicKey)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyEcdsa)

	toAddress := common.HexToAddress("0x6240D5f65CB4827f8b1C4b79EB974675382fdcb0")

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	value := big.NewInt(0)

	gasLimit := uint64(21000)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      gasLimit,
		GasPrice: gasPrice,
		Data:     []byte{},
	})

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKeyEcdsa)

	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("tx sent: %s", signedTx.Hash().Hex())
}
