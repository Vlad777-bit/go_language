package main

import (
	"errors"
	"strings"
	"testing"
)

func TestAddTransactionWithinBudget(t *testing.T) {
	ledger := NewLedger()

	ledger.SetBudget(Budget{
		Category: "еда",
		Limit:    5000,
	})

	tx := Transaction{
		Category: "еда",
		Amount:   3000,
	}

	err := ledger.AddTransaction(tx)
	if err != nil {
		t.Fatalf("expected transaction to be added, got error: %v", err)
	}

	if len(ledger.Transactions) != 1 {
		t.Fatalf(
			"expected 1 transaction, got %d",
			len(ledger.Transactions),
		)
	}
}

func TestAddTransactionExceedsBudget(t *testing.T) {
	ledger := NewLedger()

	ledger.SetBudget(Budget{
		Category: "еда",
		Limit:    5000,
	})

	err := ledger.AddTransaction(Transaction{
		Category: "еда",
		Amount:   3000,
	})
	if err != nil {
		t.Fatalf(
			"unexpected error for first transaction: %v",
			err,
		)
	}

	err = ledger.AddTransaction(Transaction{
		Category: "еда",
		Amount:   2500,
	})

	if err == nil {
		t.Fatal("expected budget exceeded error, got nil")
	}

	if !errors.Is(err, ErrBudgetExceeded) {
		t.Fatalf(
			"expected ErrBudgetExceeded, got: %v",
			err,
		)
	}

	// Проверяем самое важное:
	// отклонённая транзакция не должна сохраниться.
	if len(ledger.Transactions) != 1 {
		t.Fatalf(
			"expected rejected transaction not to be saved; got %d transactions",
			len(ledger.Transactions),
		)
	}
}

func TestAddTransactionWithoutBudget(t *testing.T) {
	ledger := NewLedger()

	tx := Transaction{
		Category: "книги",
		Amount:   10000,
	}

	err := ledger.AddTransaction(tx)
	if err != nil {
		t.Fatalf(
			"transaction without budget should be allowed, got error: %v",
			err,
		)
	}

	if len(ledger.Transactions) != 1 {
		t.Fatalf(
			"expected 1 transaction, got %d",
			len(ledger.Transactions),
		)
	}
}

func TestSetBudgetUpdatesExistingBudget(t *testing.T) {
	ledger := NewLedger()

	ledger.SetBudget(Budget{
		Category: "еда",
		Limit:    5000,
	})

	ledger.SetBudget(Budget{
		Category: "еда",
		Limit:    7000,
	})

	budget, exists := ledger.Budgets["еда"]

	if !exists {
		t.Fatal("expected budget for category еда")
	}

	if budget.Limit != 7000 {
		t.Fatalf(
			"expected updated limit 7000, got %.2f",
			budget.Limit,
		)
	}
}

func TestLoadBudgets(t *testing.T) {
	ledger := NewLedger()

	jsonData := `
	[
		{
			"category": "еда",
			"limit": 5000,
			"period": "month"
		},
		{
			"category": "такси",
			"limit": 3000,
			"period": "month"
		}
	]
	`

	err := ledger.LoadBudgets(strings.NewReader(jsonData))
	if err != nil {
		t.Fatalf(
			"expected budgets to load successfully, got: %v",
			err,
		)
	}

	if len(ledger.Budgets) != 2 {
		t.Fatalf(
			"expected 2 budgets, got %d",
			len(ledger.Budgets),
		)
	}

	if ledger.Budgets["еда"].Limit != 5000 {
		t.Fatalf(
			"expected food budget 5000, got %.2f",
			ledger.Budgets["еда"].Limit,
		)
	}
}

func TestLoadBudgetsInvalidJSON(t *testing.T) {
	ledger := NewLedger()

	invalidJSON := `
	[
		{
			"category": "еда",
			"limit":
		}
	]
	`

	err := ledger.LoadBudgets(strings.NewReader(invalidJSON))

	if err == nil {
		t.Fatal("expected JSON parsing error, got nil")
	}
}
