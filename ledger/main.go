package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func main() {
	ledger := NewLedger()

	// Проверяем работу интерфейса Validatable.
	validTransaction := Transaction{
		Category: "еда",
		Amount:   1000,
		Date:     "2026-09-29",
	}

	invalidTransaction := Transaction{
		Category: "",
		Amount:   -500,
	}

	validBudget := Budget{
		Category: "еда",
		Limit:    5000,
		Period:   "month",
	}

	invalidBudget := Budget{
		Category: "",
		Limit:    0,
	}

	fmt.Println("Проверка через интерфейс Validatable:")

	if err := CheckValid(validTransaction); err != nil {
		fmt.Printf("Ошибка Transaction: %v\n", err)
	} else {
		fmt.Println("Transaction корректна")
	}

	if err := CheckValid(invalidTransaction); err != nil {
		fmt.Printf("Некорректная Transaction: %v\n", err)
	}

	if err := CheckValid(validBudget); err != nil {
		fmt.Printf("Ошибка Budget: %v\n", err)
	} else {
		fmt.Println("Budget корректен")
	}

	if err := CheckValid(invalidBudget); err != nil {
		fmt.Printf("Некорректный Budget: %v\n", err)
	}

	fmt.Println()

	// Устанавливаем начальные бюджеты.
	if err := ledger.SetBudget(Budget{
		Category: "еда",
		Limit:    5000,
		Period:   "month",
	}); err != nil {
		fmt.Printf("Ошибка установки бюджета: %v\n", err)
	}

	if err := ledger.SetBudget(Budget{
		Category: "такси",
		Limit:    3000,
		Period:   "month",
	}); err != nil {
		fmt.Printf("Ошибка установки бюджета: %v\n", err)
	}

	// Загружаем бюджеты из JSON.
	file, err := os.Open("budgets.json")
	if err != nil {
		fmt.Printf("Не удалось открыть budgets.json: %v\n", err)
	} else {
		defer file.Close()

		reader := bufio.NewReader(file)

		if err := ledger.LoadBudgets(reader); err != nil {
			fmt.Printf("Ошибка загрузки бюджетов: %v\n", err)
		} else {
			fmt.Println("Бюджеты из budgets.json успешно загружены")
		}
	}

	fmt.Println()

	// Валидная транзакция в пределах бюджета.
	tx1 := Transaction{
		Category: "еда",
		Amount:   3000,
		Date:     "2026-09-29",
	}

	if err := ledger.AddTransaction(tx1); err != nil {
		fmt.Printf("Транзакция отклонена: %v\n", err)
	} else {
		fmt.Printf("Транзакция успешно добавлена: %+v\n", tx1)
	}

	// Валидная транзакция, но превышает бюджет.
	tx2 := Transaction{
		Category: "еда",
		Amount:   2500,
		Date:     "2026-09-29",
	}

	if err := ledger.AddTransaction(tx2); err != nil {
		if errors.Is(err, ErrBudgetExceeded) {
			fmt.Printf(
				"Транзакция отклонена из-за превышения бюджета: %v\n",
				err,
			)
		} else {
			fmt.Printf("Ошибка транзакции: %v\n", err)
		}
	}

	// Невалидная транзакция.
	tx3 := Transaction{
		Category: "",
		Amount:   -100,
	}

	if err := ledger.AddTransaction(tx3); err != nil {
		fmt.Printf("Невалидная транзакция отклонена: %v\n", err)
	}

	fmt.Println()
	fmt.Println("Сохранённые транзакции:")

	for _, tx := range ledger.Transactions {
		fmt.Printf(
			"Категория: %s, сумма: %.2f руб., дата: %s\n",
			tx.Category,
			tx.Amount,
			tx.Date,
		)
	}
}
