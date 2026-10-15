package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Vlad777-bit/ledger"
)

type Handler struct {
	ledger *ledger.Ledger
}

func NewHandler(l *ledger.Ledger) *Handler {
	return &Handler{
		ledger: l,
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set(
		"Content-Type",
		"application/json; charset=utf-8",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(
		w,
		status,
		map[string]string{
			"error": message,
		},
	)
}

func (h *Handler) CreateTransaction(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateTransactionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid JSON body",
		)
		return
	}

	tx := toTransaction(req)

	if err := tx.Validate(); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	if err := h.ledger.AddTransaction(tx); err != nil {
		if errors.Is(err, ledger.ErrBudgetExceeded) {
			writeError(
				w,
				http.StatusConflict,
				"budget exceeded",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		transactionToResponse(tx),
	)
}

func (h *Handler) ListTransactions(
	w http.ResponseWriter,
	r *http.Request,
) {
	transactions := h.ledger.ListTransactions()

	response := make(
		[]TransactionResponse,
		0,
		len(transactions),
	)

	for _, tx := range transactions {
		response = append(
			response,
			transactionToResponse(tx),
		)
	}

	writeJSON(
		w,
		http.StatusOK,
		response,
	)
}

func (h *Handler) CreateBudget(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req CreateBudgetRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid JSON body",
		)
		return
	}

	budget := toBudget(req)

	if err := budget.Validate(); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	if err := h.ledger.SetBudget(budget); err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"internal server error",
		)
		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		budgetToResponse(budget),
	)
}

func (h *Handler) ListBudgets(
	w http.ResponseWriter,
	r *http.Request,
) {
	budgets := h.ledger.ListBudgets()

	response := make(
		[]BudgetResponse,
		0,
		len(budgets),
	)

	for _, budget := range budgets {
		response = append(
			response,
			budgetToResponse(budget),
		)
	}

	writeJSON(
		w,
		http.StatusOK,
		response,
	)
}
