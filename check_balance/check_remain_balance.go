package check_balance

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
)

func CheckRemainBalance() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/ZYOdzRs9UVWVsb715Ctarx9O0CxvGs-f")
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress("0x6240D5f65CB4827f8b1C4b79EB974675382fdcb0")

	balance, err := client.BalanceAt(context.Background(), address, nil)

	fmt.Println(balance)

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	fmt.Println(ethValue)

}
