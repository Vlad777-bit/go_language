package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var ErrBudgetExceeded = errors.New("budget exceeded")

// Transaction описывает финансовую транзакцию.
type Transaction struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

// Budget описывает бюджет для конкретной категории.
type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period,omitempty"`
}

// Ledger хранит транзакции и бюджеты.
type Ledger struct {
	Transactions []Transaction
	Budgets      map[string]Budget
}

// NewLedger создаёт новый Ledger.
func NewLedger() *Ledger {
	return &Ledger{
		Transactions: make([]Transaction, 0),
		Budgets:      make(map[string]Budget),
	}
}

// SetBudget добавляет новый бюджет или обновляет существующий.
func (l *Ledger) SetBudget(b Budget) {
	l.Budgets[b.Category] = b
}

// AddTransaction добавляет транзакцию,
// если после её добавления не будет превышен бюджет категории.
func (l *Ledger) AddTransaction(tx Transaction) error {
	budget, exists := l.Budgets[tx.Category]

	// Если бюджет для категории установлен,
	// считаем уже потраченную сумму.
	if exists {
		var currentAmount float64

		for _, transaction := range l.Transactions {
			if transaction.Category == tx.Category {
				currentAmount += transaction.Amount
			}
		}

		// Проверяем, не будет ли превышен лимит.
		if currentAmount+tx.Amount > budget.Limit {
			return fmt.Errorf(
				"%w: category=%s, limit=%.2f, current=%.2f, transaction=%.2f",
				ErrBudgetExceeded,
				tx.Category,
				budget.Limit,
				currentAmount,
				tx.Amount,
			)
		}
	}

	// Транзакция сохраняется только после успешной проверки.
	l.Transactions = append(l.Transactions, tx)

	return nil
}

// LoadBudgets загружает бюджеты из JSON.
func (l *Ledger) LoadBudgets(r io.Reader) error {
	var budgets []Budget

	decoder := json.NewDecoder(r)

	if err := decoder.Decode(&budgets); err != nil {
		return fmt.Errorf("failed to parse budgets JSON: %w", err)
	}

	for _, budget := range budgets {
		if budget.Category == "" {
			return errors.New("budget category cannot be empty")
		}

		if budget.Limit < 0 {
			return fmt.Errorf(
				"budget limit cannot be negative for category %q",
				budget.Category,
			)
		}

		l.SetBudget(budget)
	}

	return nil
}
