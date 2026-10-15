package api

import (
	"net/http"

	"github.com/Vlad777-bit/ledger"
)

func NewServer(l *ledger.Ledger) http.Handler {
	handler := NewHandler(l)

	mux := http.NewServeMux()

	mux.HandleFunc(
		"POST /api/transactions",
		handler.CreateTransaction,
	)

	mux.HandleFunc(
		"GET /api/transactions",
		handler.ListTransactions,
	)

	mux.HandleFunc(
		"POST /api/budgets",
		handler.CreateBudget,
	)

	mux.HandleFunc(
		"GET /api/budgets",
		handler.ListBudgets,
	)

	return LoggingMiddleware(mux)
}
