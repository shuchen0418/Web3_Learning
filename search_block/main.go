package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	client, err := ethclient.Dial("https://ethereum-rpc.publicnode.com")
	if err != nil {
		log.Fatal(err)
	}

	//通过区块数获取区块头
	blockNumber := big.NewInt(21211831)
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("区块高度", header.Number.Uint64())
	fmt.Println("创建时间", header.Time)
	fmt.Println("难度", header.Difficulty)
	fmt.Println("区块hash", header.Hash().Hex())

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("==========================================")
	fmt.Println("区块高度", block.Number().Uint64())
	fmt.Println("创建时间", block.Time())
	fmt.Println("难度", block.Difficulty().Uint64())
	fmt.Println("区块hash", block.Hash().Hex())
	fmt.Println(len(block.Transactions()))
	count, err := client.TransactionCount(context.Background(), block.Hash())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)
}
