package ledger

import (
	"errors"
	"testing"
)

func TestTransactionValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tx      Transaction
		wantErr bool
	}{
		{
			name: "valid transaction",
			tx: Transaction{
				Amount:      1000,
				Category:    "еда",
				Description: "ланч",
				Date:        "2026-10-15",
			},
			wantErr: false,
		},
		{
			name: "zero amount",
			tx: Transaction{
				Amount:   0,
				Category: "еда",
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			tx: Transaction{
				Amount:   -100,
				Category: "еда",
			},
			wantErr: true,
		},
		{
			name: "empty category",
			tx: Transaction{
				Amount:   1000,
				Category: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.tx.Validate()

			if tt.wantErr && err == nil {
				t.Fatal("expected validation error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}

func TestBudgetValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		budget  Budget
		wantErr bool
	}{
		{
			name: "valid budget",
			budget: Budget{
				Category: "еда",
				Limit:    5000,
			},
			wantErr: false,
		},
		{
			name: "zero limit",
			budget: Budget{
				Category: "еда",
				Limit:    0,
			},
			wantErr: true,
		},
		{
			name: "negative limit",
			budget: Budget{
				Category: "еда",
				Limit:    -5000,
			},
			wantErr: true,
		},
		{
			name: "empty category",
			budget: Budget{
				Category: "",
				Limit:    5000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.budget.Validate()

			if tt.wantErr && err == nil {
				t.Fatal("expected validation error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}

func TestLedgerBudgetLimit(t *testing.T) {
	l := NewLedger()
	t.Cleanup(l.Reset)

	err := l.SetBudget(Budget{
		Category: "еда",
		Limit:    5000,
	})
	if err != nil {
		t.Fatalf("set budget: %v", err)
	}

	t.Run("within limit", func(t *testing.T) {
		err := l.AddTransaction(Transaction{
			Amount:      3000,
			Category:    "еда",
			Description: "продукты",
			Date:        "2026-10-15",
		})
		if err != nil {
			t.Fatalf("expected transaction to be accepted, got %v", err)
		}

		transactions := l.ListTransactions()

		if len(transactions) != 1 {
			t.Fatalf(
				"expected 1 transaction, got %d",
				len(transactions),
			)
		}
	})

	t.Run("budget exceeded", func(t *testing.T) {
		err := l.AddTransaction(Transaction{
			Amount:      2500,
			Category:    "еда",
			Description: "ещё продукты",
			Date:        "2026-10-15",
		})

		if err == nil {
			t.Fatal("expected budget exceeded error")
		}

		if !errors.Is(err, ErrBudgetExceeded) {
			t.Fatalf(
				"expected ErrBudgetExceeded, got %v",
				err,
			)
		}

		transactions := l.ListTransactions()

		if len(transactions) != 1 {
			t.Fatalf(
				"rejected transaction must not be saved; got %d transactions",
				len(transactions),
			)
		}
	})
}

func TestSetBudget(t *testing.T) {
	tests := []struct {
		name    string
		budget  Budget
		wantErr bool
	}{
		{
			name: "valid",
			budget: Budget{
				Category: "такси",
				Limit:    3000,
			},
			wantErr: false,
		},
		{
			name: "invalid",
			budget: Budget{
				Category: "",
				Limit:    0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			l := NewLedger()

			err := l.SetBudget(tt.budget)

			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.wantErr && len(l.ListBudgets()) != 0 {
				t.Fatal("invalid budget must not be saved")
			}
		})
	}
}
