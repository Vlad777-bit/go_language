package main

import (
	"errors"
	"fmt"
)

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        string
}

var transactions = make([]Transaction, 0)

func AddTransaction(tx Transaction) error {
	if tx.Amount == 0 {
		return errors.New("transaction amount must not be zero")
	}

	tx.ID = len(transactions) + 1
	transactions = append(transactions, tx)

	return nil
}

func ListTransactions() []Transaction {
	res := make([]Transaction, len(transactions))
	copy(res, transactions)

	return res
}

func main() {
	fmt.Println("Ledger service started")

	testTransactions := []Transaction{
		{
			Amount:      25.50,
			Category:    "Food",
			Description: "Lunch",
			Date:        "2026-09-14",
		},
		{
			Amount:      60,
			Category:    "Transport",
			Description: "Train ticket",
			Date:        "2026-09-14",
		},
		{
			Amount:      120.99,
			Category:    "Shopping",
			Description: "Headphones",
			Date:        "2026-09-14",
		},
	}

	for _, tx := range testTransactions {
		if err := AddTransaction(tx); err != nil {
			fmt.Printf("failed to add transaction: %v\n", err)
		}
	}

	fmt.Println("Transactions:")

	for _, tx := range ListTransactions() {
		fmt.Printf(
			"ID: %d, Amount: %.2f, Category: %s, Description: %s, Date: %s\n",
			tx.ID,
			tx.Amount,
			tx.Category,
			tx.Description,
			tx.Date,
		)
	}
}
