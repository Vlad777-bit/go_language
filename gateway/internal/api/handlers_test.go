package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Vlad777-bit/ledger"
)

func newTestServer(t *testing.T) (*ledger.Ledger, http.Handler) {
	t.Helper()

	l := ledger.NewLedger()
	t.Cleanup(l.Reset)

	return l, NewServer(l)
}

func assertJSONContentType(
	t *testing.T,
	rec *httptest.ResponseRecorder,
) {
	t.Helper()

	const want = "application/json; charset=utf-8"

	got := rec.Header().Get("Content-Type")

	if got != want {
		t.Fatalf(
			"expected Content-Type %q, got %q",
			want,
			got,
		)
	}
}

func assertErrorResponse(
	t *testing.T,
	rec *httptest.ResponseRecorder,
) {
	t.Helper()

	var response struct {
		Error string `json:"error"`
	}

	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("decode error response: %v", err)
	}

	if response.Error == "" {
		t.Fatal("expected non-empty error message")
	}
}

func createBudget(
	t *testing.T,
	server http.Handler,
	category string,
	limit float64,
) {
	t.Helper()

	payload := CreateBudgetRequest{
		Category: category,
		Limit:    limit,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal budget request: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/budgets",
		bytes.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	rec := httptest.NewRecorder()

	server.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"create budget: expected status %d, got %d; body=%s",
			http.StatusCreated,
			rec.Code,
			rec.Body.String(),
		)
	}
}

func TestBudgetHandlers(t *testing.T) {
	t.Run("create and list", func(t *testing.T) {
		_, server := newTestServer(t)

		createBody := `{
			"category": "еда",
			"limit": 5000
		}`

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/budgets",
			strings.NewReader(createBody),
		)

		req.Header.Set(
			"Content-Type",
			"application/json",
		)

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

		assertJSONContentType(t, rec)

		var created BudgetResponse

		if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
			t.Fatalf("decode create response: %v", err)
		}

		if created.Category != "еда" {
			t.Fatalf(
				"expected category %q, got %q",
				"еда",
				created.Category,
			)
		}

		if created.Limit != 5000 {
			t.Fatalf(
				"expected limit 5000, got %.2f",
				created.Limit,
			)
		}

		listReq := httptest.NewRequest(
			http.MethodGet,
			"/api/budgets",
			nil,
		)

		listRec := httptest.NewRecorder()

		server.ServeHTTP(listRec, listReq)

		if listRec.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				listRec.Code,
			)
		}

		assertJSONContentType(t, listRec)

		var budgets []BudgetResponse

		if err := json.NewDecoder(listRec.Body).Decode(&budgets); err != nil {
			t.Fatalf("decode list response: %v", err)
		}

		if len(budgets) != 1 {
			t.Fatalf(
				"expected 1 budget, got %d",
				len(budgets),
			)
		}

		if budgets[0].Category != "еда" {
			t.Fatalf(
				"expected category %q, got %q",
				"еда",
				budgets[0].Category,
			)
		}

		if budgets[0].Limit != 5000 {
			t.Fatalf(
				"expected limit 5000, got %.2f",
				budgets[0].Limit,
			)
		}
	})

	t.Run("invalid budget", func(t *testing.T) {
		_, server := newTestServer(t)

		body := `{
			"category": "еда",
			"limit": 0
		}`

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/budgets",
			strings.NewReader(body),
		)

		req.Header.Set(
			"Content-Type",
			"application/json",
		)

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

		assertJSONContentType(t, rec)
		assertErrorResponse(t, rec)
	})

	t.Run("bad json", func(t *testing.T) {
		_, server := newTestServer(t)

		body := `{
			"category":
		}`

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/budgets",
			strings.NewReader(body),
		)

		req.Header.Set(
			"Content-Type",
			"application/json",
		)

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

		assertJSONContentType(t, rec)
		assertErrorResponse(t, rec)
	})
}

func TestTransactionHandlers(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		_, server := newTestServer(t)

		createBudget(t, server, "еда", 5000)

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

		req.Header.Set(
			"Content-Type",
			"application/json",
		)

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

		assertJSONContentType(t, rec)

		var created TransactionResponse

		if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
			t.Fatalf("decode transaction response: %v", err)
		}

		if created.Amount != 450 {
			t.Fatalf(
				"expected amount 450, got %.2f",
				created.Amount,
			)
		}

		if created.Category != "еда" {
			t.Fatalf(
				"expected category %q, got %q",
				"еда",
				created.Category,
			)
		}

		if created.Description != "ланч" {
			t.Fatalf(
				"expected description %q, got %q",
				"ланч",
				created.Description,
			)
		}

		if created.Date != "2026-10-15" {
			t.Fatalf(
				"expected date %q, got %q",
				"2026-10-15",
				created.Date,
			)
		}

		listReq := httptest.NewRequest(
			http.MethodGet,
			"/api/transactions",
			nil,
		)

		listRec := httptest.NewRecorder()

		server.ServeHTTP(listRec, listReq)

		if listRec.Code != http.StatusOK {
			t.Fatalf(
				"expected status %d, got %d",
				http.StatusOK,
				listRec.Code,
			)
		}

		assertJSONContentType(t, listRec)

		var transactions []TransactionResponse

		if err := json.NewDecoder(listRec.Body).Decode(&transactions); err != nil {
			t.Fatalf("decode transactions response: %v", err)
		}

		if len(transactions) != 1 {
			t.Fatalf(
				"expected 1 transaction, got %d",
				len(transactions),
			)
		}

		if transactions[0].Category != "еда" {
			t.Fatalf(
				"expected category %q, got %q",
				"еда",
				transactions[0].Category,
			)
		}
	})

	t.Run("exceeded", func(t *testing.T) {
		l, server := newTestServer(t)

		createBudget(t, server, "еда", 1000)

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

		firstReq.Header.Set(
			"Content-Type",
			"application/json",
		)

		firstRec := httptest.NewRecorder()

		server.ServeHTTP(firstRec, firstReq)

		if firstRec.Code != http.StatusCreated {
			t.Fatalf(
				"expected first transaction status %d, got %d; body=%s",
				http.StatusCreated,
				firstRec.Code,
				firstRec.Body.String(),
			)
		}

		exceededBody := `{
			"amount": 500,
			"category": "еда",
			"description": "вторая покупка",
			"date": "2026-10-15"
		}`

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/transactions",
			strings.NewReader(exceededBody),
		)

		req.Header.Set(
			"Content-Type",
			"application/json",
		)

		rec := httptest.NewRecorder()

		server.ServeHTTP(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf(
				"expected status %d, got %d; body=%s",
				http.StatusConflict,
				rec.Code,
				rec.Body.String(),
			)
		}

		assertJSONContentType(t, rec)

		var response struct {
			Error string `json:"error"`
		}

		if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
			t.Fatalf("decode conflict response: %v", err)
		}

		if response.Error != "budget exceeded" {
			t.Fatalf(
				"expected error %q, got %q",
				"budget exceeded",
				response.Error,
			)
		}

		transactions := l.ListTransactions()

		if len(transactions) != 1 {
			t.Fatalf(
				"expected rejected transaction not to be stored; got %d transactions",
				len(transactions),
			)
		}
	})

	t.Run("invalid transaction", func(t *testing.T) {
		_, server := newTestServer(t)

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

		req.Header.Set(
			"Content-Type",
			"application/json",
		)

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

		assertJSONContentType(t, rec)
		assertErrorResponse(t, rec)
	})

	t.Run("bad json", func(t *testing.T) {
		_, server := newTestServer(t)

		body := `{
			"amount":
		}`

		req := httptest.NewRequest(
			http.MethodPost,
			"/api/transactions",
			strings.NewReader(body),
		)

		req.Header.Set(
			"Content-Type",
			"application/json",
		)

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

		assertJSONContentType(t, rec)
		assertErrorResponse(t, rec)
	})
}
