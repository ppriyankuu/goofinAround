package account

import (
	"fmt"
	"sync"
)

type Account struct {
	balance int
	mutex   sync.Mutex
}

func NewAccount(initialBalance int) *Account {
	return &Account{
		balance: initialBalance,
	}
}

func (a *Account) Deposit(amount int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.balance += amount
	fmt.Printf("Deposited %d. New balance: %d\n", amount, a.balance)
}

func (a *Account) Withdraw(amount int) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if amount > a.balance {
		fmt.Printf("Insufficient funds for withdrawal of %d. Current balance: %d\n", amount, a.balance)
		return
	}
	a.balance -= amount
	fmt.Printf("Withdrew %d. New balance: %d\n", amount, a.balance)
}

func (a *Account) Balance() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.balance
}
