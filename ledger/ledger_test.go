package ledger

import (
	"errors"
	"strings"
	"testing"
)

func TestTransactionValidateValid(t *testing.T) {
	tx := Transaction{
		Category: "еда",
		Amount:   1000,
		Date:     "2026-09-29",
	}

	if err := tx.Validate(); err != nil {
		t.Fatalf("expected valid transaction, got error: %v", err)
	}
}

func TestTransactionValidateEmptyCategory(t *testing.T) {
	tx := Transaction{
		Category: "",
		Amount:   1000,
	}

	if err := tx.Validate(); err == nil {
		t.Fatal("expected validation error for empty category")
	}
}

func TestTransactionValidateInvalidAmount(t *testing.T) {
	tx := Transaction{
		Category: "еда",
		Amount:   0,
	}

	if err := tx.Validate(); err == nil {
		t.Fatal("expected validation error for zero amount")
	}
}

func TestBudgetValidateValid(t *testing.T) {
	budget := Budget{
		Category: "еда",
		Limit:    5000,
	}

	if err := budget.Validate(); err != nil {
		t.Fatalf("expected valid budget, got error: %v", err)
	}
}

func TestBudgetValidateInvalidLimit(t *testing.T) {
	budget := Budget{
		Category: "еда",
		Limit:    0,
	}

	if err := budget.Validate(); err == nil {
		t.Fatal("expected validation error for invalid budget limit")
	}
}

func TestCheckValid(t *testing.T) {
	tx := Transaction{
		Category: "еда",
		Amount:   1000,
	}

	budget := Budget{
		Category: "еда",
		Limit:    5000,
	}

	if err := CheckValid(tx); err != nil {
		t.Fatalf("transaction should be valid: %v", err)
	}

	if err := CheckValid(budget); err != nil {
		t.Fatalf("budget should be valid: %v", err)
	}
}

func TestSetBudget(t *testing.T) {
	ledger := NewLedger()

	err := ledger.SetBudget(Budget{
		Category: "еда",
		Limit:    5000,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ledger.Budgets["еда"].Limit != 5000 {
		t.Fatal("budget was not saved correctly")
	}
}

func TestSetBudgetRejectsInvalidBudget(t *testing.T) {
	ledger := NewLedger()

	err := ledger.SetBudget(Budget{
		Category: "",
		Limit:    0,
	})

	if err == nil {
		t.Fatal("expected validation error")
	}

	if len(ledger.Budgets) != 0 {
		t.Fatal("invalid budget must not be saved")
	}
}

func TestAddTransactionWithinBudget(t *testing.T) {
	ledger := NewLedger()

	err := ledger.SetBudget(Budget{
		Category: "еда",
		Limit:    5000,
	})
	if err != nil {
		t.Fatalf("unexpected budget error: %v", err)
	}

	err = ledger.AddTransaction(Transaction{
		Category: "еда",
		Amount:   3000,
	})
	if err != nil {
		t.Fatalf("unexpected transaction error: %v", err)
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

	err := ledger.SetBudget(Budget{
		Category: "еда",
		Limit:    5000,
	})
	if err != nil {
		t.Fatalf("unexpected budget error: %v", err)
	}

	err = ledger.AddTransaction(Transaction{
		Category: "еда",
		Amount:   3000,
	})
	if err != nil {
		t.Fatalf("unexpected first transaction error: %v", err)
	}

	err = ledger.AddTransaction(Transaction{
		Category: "еда",
		Amount:   2500,
	})

	if err == nil {
		t.Fatal("expected budget exceeded error")
	}

	if !errors.Is(err, ErrBudgetExceeded) {
		t.Fatalf("expected ErrBudgetExceeded, got: %v", err)
	}

	if len(ledger.Transactions) != 1 {
		t.Fatal("rejected transaction must not be saved")
	}
}

func TestAddTransactionRejectsInvalidTransaction(t *testing.T) {
	ledger := NewLedger()

	err := ledger.AddTransaction(Transaction{
		Category: "",
		Amount:   -100,
	})

	if err == nil {
		t.Fatal("expected validation error")
	}

	if len(ledger.Transactions) != 0 {
		t.Fatal("invalid transaction must not be saved")
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
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ledger.Budgets) != 2 {
		t.Fatalf(
			"expected 2 budgets, got %d",
			len(ledger.Budgets),
		)
	}
}

func TestLoadBudgetsRejectsInvalidBudget(t *testing.T) {
	ledger := NewLedger()

	jsonData := `
	[
		{
			"category": "",
			"limit": 0
		}
	]
	`

	err := ledger.LoadBudgets(strings.NewReader(jsonData))

	if err == nil {
		t.Fatal("expected validation error")
	}

	if len(ledger.Budgets) != 0 {
		t.Fatal("invalid budget must not be saved")
	}
}
