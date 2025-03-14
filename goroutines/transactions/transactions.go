package transactions

import (
	"go-routines/account"
	"math/rand"
	"time"
)

func SimulateTransactions(acc *account.Account, numTransactions int) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < numTransactions; i++ {
		go func(id int) {
			amount := rand.Intn(100) + 1
			if rand.Intn(2) == 0 {
				acc.Deposit(amount)
			} else {
				acc.Withdraw(amount)
			}
		}(i)
	}
}

// waits for all the transactions to complete
func WaitForTransactions(numTransactions int) {
	time.Sleep(time.Duration(numTransactions*50) * time.Millisecond)
}
