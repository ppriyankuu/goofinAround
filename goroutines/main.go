package main

import (
	"fmt"
	"go-routines/account"
	"go-routines/transactions"
	"time"
)

func main() {
	initialBalance := 1000
	numTransactions := 100

	acc := account.NewAccount(initialBalance)

	transactions.SimulateTransactions(acc, numTransactions)

	transactions.WaitForTransactions(numTransactions)

	finalBalance := acc.Balance()
	fmt.Printf("Final balance after %d transactions: %d\n", numTransactions, finalBalance)

	time.Sleep(1 * time.Second)
}
