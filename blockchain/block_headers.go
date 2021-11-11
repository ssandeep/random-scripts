package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
	client, err := ethclient.Dial("https://rpc-mainnet.maticvigil.com/")
	if err != nil {
		log.Fatal(err)
	}

	latestBlockHeader, err := client.HeaderByNumber(context.Background(), nil) // most recent block
	printBlockRLP(latestBlockHeader)

	noTestCases := 10
	maxBlockNum := latestBlockHeader.Number

	for i := 0; i < noTestCases; i++ {
		// Generate cryptographically strong pseudo-random between 0 - max
		n, err := rand.Int(rand.Reader, maxBlockNum)
		if err != nil {
			log.Fatal((err))
		}
		header, err := client.HeaderByNumber(context.Background(), n)
		if err != nil {
			log.Fatal(err)
		}
		printBlockRLP(header)
	}
}

func printBlockRLP(header *types.Header) {
	data, err := rlp.EncodeToBytes(header)
	if err != nil {
		log.Fatal(err)
	}
	output := hex.EncodeToString(data)
	fmt.Printf("Block Number: %v\n", header.Number)
	fmt.Printf("0x%s\n\n", output)

	strOutput := string(data)
	fmt.Println(strOutput + "\n")
}
