package mempool

import (
	"fmt"
	"mempool/config"
	"mempool/internal/domain"
	"sort"
	"sync"

	"github.com/google/uuid"
)

type Mempool struct {
	mu                        sync.Mutex
	TransactionsToSubmit      map[string]domain.Transaction
	TransactionIDsSortedByFee []string
	TransactionsToProcess     map[string]domain.Transaction
	Processing                bool
}

func NewMempool() *Mempool {
	return &Mempool{
		TransactionsToSubmit:      make(map[string]domain.Transaction),
		TransactionIDsSortedByFee: make([]string, 0),
		TransactionsToProcess:     make(map[string]domain.Transaction),
	}
}

func (m *Mempool) AddTransactionsToProcess(ts []domain.Transaction) {
	if len(ts) == 0 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Printf("Adding %d transactions to process\n", len(ts))

	for _, t := range ts {
		uniqueId := fmt.Sprintf("%s-%d", t.ID, uuid.New().ID)
		t.ID = uniqueId
		m.TransactionsToProcess[uniqueId] = t
	}
}

func (m *Mempool) AddTransactionToSubmit(t domain.Transaction) {
	if t.GasPrice > config.GasLimit {
		fmt.Println("Transaction gas price exceeds gas limit, discarding:", t)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	fmt.Printf("Adding transaction to submit: %+v\n", t)

	m.TransactionsToSubmit[t.ID] = t

	idx := sort.Search(len(m.TransactionIDsSortedByFee), func(i int) bool {
		return m.TransactionsToSubmit[m.TransactionIDsSortedByFee[i]].Fee <= t.Fee
	})
	m.TransactionIDsSortedByFee = append(m.TransactionIDsSortedByFee, "")
	copy(m.TransactionIDsSortedByFee[idx+1:], m.TransactionIDsSortedByFee[idx:])
	m.TransactionIDsSortedByFee[idx] = t.ID

	delete(m.TransactionsToProcess, t.ID)
}

func (m *Mempool) GetTransactionsBatchToSubmit() []domain.Transaction {
	var transactionsBatchIds []string

	if len(m.TransactionIDsSortedByFee) == 0 {
		return nil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	totalGas := 0
	batchIsFull := false
	for _, tId := range m.TransactionIDsSortedByFee {
		var tGasPrice = m.TransactionsToSubmit[tId].GasPrice
		if totalGas+tGasPrice > config.GasLimit {
			batchIsFull = true
			break
		}
		transactionsBatchIds = append(transactionsBatchIds, tId)
		totalGas += tGasPrice
	}

	if !batchIsFull {
		return nil
	}

	var transactionsBatch []domain.Transaction
	for _, tId := range transactionsBatchIds {
		transactionsBatch = append(transactionsBatch, m.TransactionsToSubmit[tId])
	}

	fmt.Printf("TransactionsToSubmit to submit in next batch: %+v\n", transactionsBatch)

	return transactionsBatch
}

func (m *Mempool) RemoveTransactions(ts []domain.Transaction) {
	if len(ts) == 0 {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	idsToRemove := make(map[string]struct{}, len(ts))
	for _, t := range ts {
		idsToRemove[t.ID] = struct{}{}
		delete(m.TransactionsToSubmit, t.ID)
	}

	filteredTransactionIDsSorted := []string{}
	for _, sortedID := range m.TransactionIDsSortedByFee {
		if _, found := idsToRemove[sortedID]; !found {
			filteredTransactionIDsSorted = append(filteredTransactionIDsSorted, sortedID)
		}
	}
	m.TransactionIDsSortedByFee = filteredTransactionIDsSorted

	fmt.Printf("Processed transactions: %+v\n", ts)
}

func (m *Mempool) GetTransactionsToProcess() map[string]domain.Transaction {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.TransactionsToProcess
}
