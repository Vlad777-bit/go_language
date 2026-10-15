package ledger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var ErrBudgetExceeded = errors.New("budget exceeded")

// Validatable задаёт единый контракт валидации данных.
type Validatable interface {
	Validate() error
}

// Transaction описывает финансовую транзакцию.
type Transaction struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

// Validate проверяет корректность транзакции.
func (t Transaction) Validate() error {
	if t.Category == "" {
		return errors.New("transaction category cannot be empty")
	}

	if t.Amount <= 0 {
		return errors.New("transaction amount must be greater than zero")
	}

	return nil
}

// Budget описывает бюджет для категории.
type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period,omitempty"`
}

// Validate проверяет корректность бюджета.
func (b Budget) Validate() error {
	if b.Category == "" {
		return errors.New("budget category cannot be empty")
	}

	if b.Limit <= 0 {
		return errors.New("budget limit must be greater than zero")
	}

	return nil
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

// CheckValid демонстрирует полиморфизм через интерфейс Validatable.
func CheckValid(v Validatable) error {
	return v.Validate()
}

// SetBudget добавляет новый бюджет или обновляет существующий.
func (l *Ledger) SetBudget(b Budget) error {
	if err := b.Validate(); err != nil {
		return fmt.Errorf("invalid budget: %w", err)
	}

	l.Budgets[b.Category] = b

	return nil
}

// AddTransaction добавляет транзакцию после валидации
// и проверки бюджета.
func (l *Ledger) AddTransaction(tx Transaction) error {
	if err := tx.Validate(); err != nil {
		return fmt.Errorf("invalid transaction: %w", err)
	}

	budget, exists := l.Budgets[tx.Category]

	if exists {
		var currentAmount float64

		for _, transaction := range l.Transactions {
			if transaction.Category == tx.Category {
				currentAmount += transaction.Amount
			}
		}

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
		if err := l.SetBudget(budget); err != nil {
			return fmt.Errorf(
				"failed to add budget for category %q: %w",
				budget.Category,
				err,
			)
		}
	}

	return nil
}

func (l *Ledger) ListTransactions() []Transaction {
	result := make([]Transaction, len(l.Transactions))
	copy(result, l.Transactions)

	return result
}

func (l *Ledger) ListBudgets() []Budget {
	result := make([]Budget, 0, len(l.Budgets))

	for _, budget := range l.Budgets {
		result = append(result, budget)
	}

	return result
}

func (l *Ledger) Reset() {
	l.Transactions = make([]Transaction, 0)
	l.Budgets = make(map[string]Budget)
}
