package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	userName := os.Getenv("USER")
	if userName == "" {
		userName = os.Getenv("USERNAME")
	}

	fmt.Println("Имя пользователя:", userName)
	fmt.Println("Аргументы CLI:", os.Args[1:])
	fmt.Println("Версия GO:", runtime.Version())
}
