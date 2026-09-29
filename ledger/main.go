package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func main() {
	ledger := NewLedger()

	// Начальные бюджеты.
	ledger.SetBudget(Budget{
		Category: "еда",
		Limit:    5000,
		Period:   "month",
	})

	ledger.SetBudget(Budget{
		Category: "такси",
		Limit:    3000,
		Period:   "month",
	})

	fmt.Println("Начальные бюджеты:")
	for category, budget := range ledger.Budgets {
		fmt.Printf(
			"%s: %.2f руб.\n",
			category,
			budget.Limit,
		)
	}

	fmt.Println()

	// Загружаем бюджеты из JSON-файла.
	file, err := os.Open("budgets.json")
	if err != nil {
		fmt.Printf("Не удалось открыть budgets.json: %v\n", err)
	} else {
		defer file.Close()

		reader := bufio.NewReader(file)

		if err = ledger.LoadBudgets(reader); err != nil {
			fmt.Printf("Ошибка загрузки бюджетов: %v\n", err)
		} else {
			fmt.Println("Бюджеты из budgets.json успешно загружены.")
		}
	}

	fmt.Println()

	// Транзакция в пределах бюджета.
	tx1 := Transaction{
		Category: "еда",
		Amount:   3000,
	}

	err = ledger.AddTransaction(tx1)
	if err != nil {
		fmt.Printf("Транзакция отклонена: %v\n", err)
	} else {
		fmt.Printf(
			"Транзакция успешно добавлена: %+v\n",
			tx1,
		)
	}

	// Эта транзакция должна превысить бюджет.
	tx2 := Transaction{
		Category: "еда",
		Amount:   2500,
	}

	err = ledger.AddTransaction(tx2)
	if err != nil {
		if errors.Is(err, ErrBudgetExceeded) {
			fmt.Printf(
				"Транзакция отклонена из-за превышения бюджета: %v\n",
				err,
			)
		} else {
			fmt.Printf("Ошибка добавления транзакции: %v\n", err)
		}
	} else {
		fmt.Printf(
			"Транзакция успешно добавлена: %+v\n",
			tx2,
		)
	}

	fmt.Println()
	fmt.Println("Сохранённые транзакции:")

	for _, tx := range ledger.Transactions {
		fmt.Printf(
			"Категория: %s, сумма: %.2f руб.\n",
			tx.Category,
			tx.Amount,
		)
	}
}
