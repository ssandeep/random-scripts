package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

/*

This script gets all the blocks in which state sync txns are present on Mumbai (matic pos testnet) chain from polygonscan. It then checks if there are any difference between the no of txns in these blocks on local node vs remote rpc url and reports such occurrences.

Mumbai API for getting all state sync txs:
https://api-testnet.polygonscan.com/api?module=account&action=txlist&address=0x0000000000000000000000000000000000000000&startblock=1&endblock=99999999&page=1&offset=1000&sort=asc&apikey=YourApiKeyToken


https://api-testnet.polygonscan.com/api?module=account&action=txlist&address=0x0000000000000000000000000000000000000000&startblock=1&endblock=1000&page=1&offset=1000&sort=asc



1. Get all state sync txs in a block range from polygonscan
2. For these txs, check if we have these txs in our localhost bor rpc
3. If no, append output to a file

*/

type PolygonScanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		BlockNumber       string `json:"blockNumber"`
		TimeStamp         string `json:"timeStamp"`
		Hash              string `json:"hash"`
		Nonce             string `json:"nonce"`
		BlockHash         string `json:"blockHash"`
		TransactionIndex  string `json:"transactionIndex"`
		From              string `json:"from"`
		To                string `json:"to"`
		Value             string `json:"value"`
		Gas               string `json:"gas"`
		GasPrice          string `json:"gasPrice"`
		IsError           string `json:"isError"`
		TxreceiptStatus   string `json:"txreceipt_status"`
		Input             string `json:"input"`
		ContractAddress   string `json:"contractAddress"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		GasUsed           string `json:"gasUsed"`
		Confirmations     string `json:"confirmations"`
	} `json:"result"`
}

type Tx struct {
	BlockNumber string
	Hash        string
}

func getStateSyncTxns(start, end int) []Tx {
	var txs []Tx
	psMumbaiApiUrl := "https://api-testnet.polygonscan.com/api?module=account&action=txlist&address=0x0000000000000000000000000000000000000000&startblock=" + strconv.Itoa(start) + "&endblock=" + strconv.Itoa(end) + "&sort=asc"
	fmt.Println("Fetching data from ", psMumbaiApiUrl)

	resp, err := http.Get(psMumbaiApiUrl)
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) // response body is []byte

	var result PolygonScanResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}

	// fmt.Println(PrettyPrint(result))

	for _, rec := range result.Result {
		fmt.Println(rec.BlockNumber)
		txs = append(txs, Tx{BlockNumber: rec.BlockNumber, Hash: rec.Hash})
	}
	return txs
}

// func PrettyPrint(i interface{}) string {
// 	s, _ := json.MarshalIndent(i, "", "\t")
// 	return string(s)
// }

func main() {
	var txs []Tx
	maxBlockNo := 18290000 //25000000
	currBlockNo := 18280000
	for currBlockNo < maxBlockNo {
		nextBlockNo := currBlockNo + 50000
		txs = getStateSyncTxns(currBlockNo, nextBlockNo)
		checkTxs(txs)
		currBlockNo = nextBlockNo
	}
}

func checkTxs(txs []Tx) {

	// curl localhost:8545 -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0xd96ecec3ac99e7e0f1edc62cff7d349c8c51cbfd0efc72f00662ecee6d41b14a"],"id":0}'

	url := "http://localhost:8545"

	for _, tx := range txs {

		var jsonStr = []byte(`{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["` + tx.Hash + `"],"id":0}`)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
		if `{"jsonrpc":"2.0","id":0,"result":null}` == string(body) {
			panic("~~~~~\n\n\nyayyyyyyyyyyyyyyyy\n\n\n~~~~~")
		}

	}
}