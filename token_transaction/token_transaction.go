package token_transaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/crypto/sha3"
	"log"
	"math/big"
)

func TokenTransaction() {

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
	paddedToAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	tokenAddress := common.HexToAddress("0x28b149020d2152179873ec60bed6bf7cd705775d")

	transferFnSignature := []byte("transfer(address,uint256)")

	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	fmt.Println(methodID)

	amount := new(big.Int)
	amount.SetString("10000000000000000000000", 10)
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	fmt.Println(hexutil.Encode(paddedAmount))

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedToAddress...)
	data = append(data, paddedAmount...)

	gasLimit := uint64(30000000)

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &tokenAddress,
		Value:    big.NewInt(0),
		GasPrice: gasPrice,
		Gas:      gasLimit,
		Data:     data,
	})

	signTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKeyEcdsa)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(signTx.Hash().Hex())
}
