package domain

type Transaction struct {
	ID       string `json:"id"`
	GasPrice int    `json:"gasPrice"`
	Fee      int    `json:"fee"`
}

type TransactionPrice struct {
	GasPrice int `json:"gasPrice"`
	Fee      int `json:"fee"`
}
