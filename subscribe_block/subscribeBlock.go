package subscribe_block

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func SubscribeBlock() {
	client, err := ethclient.Dial("wss://eth-testnet.4everland.org/ws/v1/37fa9972c1b1cd5fab542c7bdd4cde2f")
	if err != nil {
		log.Fatal(err)
	}

	headers := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println(header.Hash().Hex())
			fmt.Println(header.Time)
			fmt.Println(header.Number.Uint64())
			fmt.Println(header.Nonce.Uint64())
		}
	}

}
