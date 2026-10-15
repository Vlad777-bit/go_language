package main

import (
	"fmt"

	"github.com/Vlad777-bit/ledger"
)

func main() {
	l := ledger.NewLedger()

	err := l.SetBudget(ledger.Budget{
		Category: "еда",
		Limit:    5000,
	})
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Ledger запущен")
}
