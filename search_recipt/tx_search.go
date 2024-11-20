package main

func main() {
	//client, err := ethclient.Dial("https://ethereum-sepolia-rpc.publicnode.com")
	//if err != nil {
	//	log.Fatal(err)
	//}

	/*chainId, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	blockNumber := big.NewInt(7100527)

	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash().Hex())
		fmt.Println(tx.Value().String())
		fmt.Println(tx.Gas())
		fmt.Println(tx.GasPrice().Uint64())
		fmt.Println(tx.Nonce())
		fmt.Println(tx.Data())
		fmt.Println(tx.To().Hex())

		if sender, err := types.Sender(types.NewEIP155Signer(chainId), tx); err == nil {
			fmt.Println("sender = ", sender.Hex())
		} else {
			log.Fatal(err)
		}

		if receipt, err := client.TransactionReceipt(context.Background(), tx.Hash()); err == nil {
			fmt.Println("receipt status = ", receipt.Status)
			fmt.Println("receipt logs = ", receipt.Logs)
		}
		break
	}

	fmt.Println("==========================================")

	blockHash := common.HexToHash("0x7493c882f0507f51490e1468d83a10008d6c51a84ead650a57349348c7dc7987")

	count, err := client.TransactionCount(context.Background(), blockHash)
	if err != nil {
		log.Fatal(err)
	}

	for idx := uint(0); idx < count; idx++ {
		tx, err := client.TransactionInBlock(context.Background(), blockHash, idx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("tx hash = ", tx.Hash().Hex())
		break
	}

	fmt.Println("==========================================")

	txHash := common.HexToHash("0x89faebbbf7f40a6533776bbddd7ab3d54c517ec358a338acee42f416702b0491")

	tx, pending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("tx hash = ", tx.Hash().Hex())
	fmt.Println("tx pending = ", pending)*/
	//blockNumber := big.NewInt(7100527)
	////blockHash := common.HexToHash("0x7493c882f0507f51490e1468d83a10008d6c51a84ead650a57349348c7dc7987")
	////receiptByHash, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithHash(blockHash, false))
	//
	//receiptByNumber, err := client.BlockReceipts(context.Background(), rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(blockNumber.Int64())))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, receipt := range receiptByNumber {
	//	fmt.Println("receipt status = ", receipt.Status)
	//	fmt.Println("receipt logs = ", receipt.Logs)
	//	fmt.Println("receipt tx hash = ", receipt.TxHash.Hex())
	//	fmt.Println("receipt transaction index = ", receipt.TransactionIndex)
	//}
}
