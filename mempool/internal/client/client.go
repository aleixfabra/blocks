package client

import (
	"bytes"
	"encoding/json"
	"io"
	"mempool/config"
	"mempool/internal/domain"
	"net/http"
)

type SimulateBlockRequest struct {
	Transactions []domain.Transaction `json:"transactions"`
}

type SimulateBlockResponse struct {
	Transactions          []domain.Transaction `json:"transactions"`
	TotalGas              int                  `json:"totalGas"`
	TotalFees             int                  `json:"totalFees"`
	GasLimit              int                  `json:"gasLimit"`
	ProcessingTimeSeconds float64              `json:"processingTimeSeconds"`
}

func GetCurrentPrice() (*domain.TransactionPrice, error) {
	resp, err := http.Get(config.BlocksURL + "/getCurrentPrice")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var transactionPrice domain.TransactionPrice
	if err := json.Unmarshal(body, &transactionPrice); err != nil {
		return nil, err
	}

	return &transactionPrice, nil
}

func SimulateBlock(ts []domain.Transaction) (*SimulateBlockResponse, error) {
	var simulateBlockRequest SimulateBlockRequest

	simulateBlockRequest.Transactions = ts

	payload, err := json.Marshal(simulateBlockRequest)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(config.BlocksURL+"/simulateBlock", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var simulateBlockResponse SimulateBlockResponse
	body, err := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &simulateBlockResponse); err != nil {
		return nil, err
	}

	return &simulateBlockResponse, nil
}
