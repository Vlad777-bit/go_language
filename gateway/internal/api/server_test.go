package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Vlad777-bit/ledger"
)

func newTestServer() http.Handler {
	l := ledger.NewLedger()

	return NewServer(l)
}

func TestListTransactionsEmpty(t *testing.T) {
	server := newTestServer()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/transactions",
		nil,
	)

	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	if rec.Header().Get("Content-Type") != "application/json; charset=utf-8" {
		t.Fatalf(
			"unexpected Content-Type: %s",
			rec.Header().Get("Content-Type"),
		)
	}

	var response []TransactionResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 0 {
		t.Fatalf(
			"expected empty transactions list, got %d items",
			len(response),
		)
	}
}

func TestCreateBudget(t *testing.T) {
	server := newTestServer()

	body := `{
		"category": "еда",
		"limit": 5000
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/budgets",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d; body=%s",
			http.StatusCreated,
			rec.Code,
			rec.Body.String(),
		)
	}

	var response BudgetResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Category != "еда" {
		t.Fatalf(
			"expected category еда, got %q",
			response.Category,
		)
	}

	if response.Limit != 5000 {
		t.Fatalf(
			"expected limit 5000, got %.2f",
			response.Limit,
		)
	}
}

func TestCreateInvalidBudget(t *testing.T) {
	server := newTestServer()

	body := `{
		"category": "",
		"limit": 0
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/budgets",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}

	var response map[string]string

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if response["error"] == "" {
		t.Fatal("expected error message")
	}
}

func TestListBudgets(t *testing.T) {
	l := ledger.NewLedger()

	err := l.SetBudget(ledger.Budget{
		Category: "еда",
		Limit:    5000,
	})
	if err != nil {
		t.Fatalf("failed to create budget: %v", err)
	}

	server := NewServer(l)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/budgets",
		nil,
	)

	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			rec.Code,
		)
	}

	var response []BudgetResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 1 {
		t.Fatalf(
			"expected 1 budget, got %d",
			len(response),
		)
	}

	if response[0].Category != "еда" {
		t.Fatalf(
			"expected category еда, got %q",
			response[0].Category,
		)
	}
}

func TestCreateTransaction(t *testing.T) {
	l := ledger.NewLedger()

	err := l.SetBudget(ledger.Budget{
		Category: "еда",
		Limit:    5000,
	})
	if err != nil {
		t.Fatalf("failed to set budget: %v", err)
	}

	server := NewServer(l)

	body := `{
		"amount": 450,
		"category": "еда",
		"description": "ланч",
		"date": "2026-10-15"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/transactions",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d; body=%s",
			http.StatusCreated,
			rec.Code,
			rec.Body.String(),
		)
	}

	var response TransactionResponse

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Amount != 450 {
		t.Fatalf(
			"expected amount 450, got %.2f",
			response.Amount,
		)
	}

	if response.Category != "еда" {
		t.Fatalf(
			"expected category еда, got %q",
			response.Category,
		)
	}

	if response.Description != "ланч" {
		t.Fatalf(
			"expected description ланч, got %q",
			response.Description,
		)
	}
}

func TestCreateInvalidTransaction(t *testing.T) {
	server := newTestServer()

	body := `{
		"amount": -100,
		"category": "еда",
		"description": "ошибка",
		"date": "2026-10-15"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/transactions",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d; body=%s",
			http.StatusBadRequest,
			rec.Code,
			rec.Body.String(),
		)
	}

	var response map[string]string

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["error"] == "" {
		t.Fatal("expected validation error")
	}
}

func TestCreateTransactionBudgetExceeded(t *testing.T) {
	l := ledger.NewLedger()

	err := l.SetBudget(ledger.Budget{
		Category: "еда",
		Limit:    1000,
	})
	if err != nil {
		t.Fatalf("failed to set budget: %v", err)
	}

	server := NewServer(l)

	firstBody := `{
		"amount": 800,
		"category": "еда",
		"description": "первая покупка",
		"date": "2026-10-15"
	}`

	firstReq := httptest.NewRequest(
		http.MethodPost,
		"/api/transactions",
		strings.NewReader(firstBody),
	)

	firstReq.Header.Set("Content-Type", "application/json")

	firstRec := httptest.NewRecorder()

	server.ServeHTTP(firstRec, firstReq)

	if firstRec.Code != http.StatusCreated {
		t.Fatalf(
			"first transaction: expected status %d, got %d; body=%s",
			http.StatusCreated,
			firstRec.Code,
			firstRec.Body.String(),
		)
	}

	secondBody := `{
		"amount": 500,
		"category": "еда",
		"description": "вторая покупка",
		"date": "2026-10-15"
	}`

	secondReq := httptest.NewRequest(
		http.MethodPost,
		"/api/transactions",
		strings.NewReader(secondBody),
	)

	secondReq.Header.Set("Content-Type", "application/json")

	secondRec := httptest.NewRecorder()

	server.ServeHTTP(secondRec, secondReq)

	if secondRec.Code != http.StatusConflict {
		t.Fatalf(
			"expected status %d, got %d; body=%s",
			http.StatusConflict,
			secondRec.Code,
			secondRec.Body.String(),
		)
	}

	var response map[string]string

	if err := json.NewDecoder(secondRec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["error"] != "budget exceeded" {
		t.Fatalf(
			"expected budget exceeded error, got %q",
			response["error"],
		)
	}

	// Проверяем главное требование:
	// отклонённая транзакция не сохранилась.
	transactions := l.ListTransactions()

	if len(transactions) != 1 {
		t.Fatalf(
			"expected only 1 saved transaction, got %d",
			len(transactions),
		)
	}
}

func TestCreateTransactionInvalidJSON(t *testing.T) {
	server := newTestServer()

	body := `{
		"amount": 100,
		"category":
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/transactions",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			rec.Code,
		)
	}

	var response map[string]string

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response["error"] != "invalid JSON body" {
		t.Fatalf(
			"expected invalid JSON body error, got %q",
			response["error"],
		)
	}
}
