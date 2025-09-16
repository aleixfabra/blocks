package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mempool/internal/client"
	"mempool/internal/domain"
	"mempool/internal/mempool"
	"net/http"
)

var mp = mempool.NewMempool()

func fetchTransactions() {
	for {
		for _, t := range mp.GetTransactionsToProcess() {
			err := fetchTransaction(t)
			if err != nil {
				fmt.Println("Error fetching transaction:", err)
			}
		}
	}
}

func fetchTransaction(t domain.Transaction) error {
	tp, err := client.GetCurrentPrice()
	if err != nil {
		return err
	}

	tx := domain.Transaction{
		ID:       t.ID,
		GasPrice: tp.GasPrice,
		Fee:      tp.Fee,
	}

	mp.AddTransactionToSubmit(tx)

	return nil
}

func submitTransactions() {
	for {
		err := submitTransaction()
		if err != nil {
			fmt.Println("Error submitting transactions batch:", err)
		}
	}
}

func submitTransaction() error {
	tb := mp.GetTransactionsBatchToSubmit()
	if tb == nil {
		return nil
	}

	simulateBlockResponse, err := client.SimulateBlock(tb)
	if err != nil {
		return err
	}

	mp.RemoveTransactions(simulateBlockResponse.Transactions)

	return nil
}

func submitTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	var transactions client.SimulateBlockRequest
	if err := json.NewDecoder(r.Body).Decode(&transactions); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	mp.AddTransactionsToProcess(transactions.Transactions)
}

func main() {
	http.HandleFunc("/submitTransactions", submitTransactionsHandler)

	go fetchTransactions()
	go submitTransactions()

	log.Println("Starting server on :9090...")

	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}
}
