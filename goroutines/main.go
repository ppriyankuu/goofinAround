package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type BankAccount struct {
	balance int
	mutex   sync.Mutex
}

func (account *BankAccount) Deposit(amount int, wg *sync.WaitGroup, r *rand.Rand) {
	defer wg.Done() // ensures the wait group is notified when this function exits

	// simulating some processing time for the deposit
	time.Sleep(time.Duration(r.Intn(100)) * time.Millisecond)

	account.mutex.Lock()         // locks the mutex to protect the balance
	defer account.mutex.Unlock() // ensures the mutex is unlocked when the function exits

	fmt.Printf("Depositing %d to account\n", amount)
	account.balance += amount
	fmt.Printf("New balance after deposit: %d\n", account.balance)
}

func (account *BankAccount) Withdraw(amount int, wg *sync.WaitGroup, r *rand.Rand) {
	defer wg.Done() // ensures the wait group is notified when this function exits

	// simulates some processing time for the withdrawal
	time.Sleep(time.Duration(r.Intn(100)) * time.Millisecond)

	account.mutex.Lock()         // locks the mutex to protect the balance
	defer account.mutex.Unlock() // ensures the mutex is unlocked when the function exits

	if account.balance >= amount {
		fmt.Printf("Withdrawing %d from account\n", amount)
		account.balance -= amount
		fmt.Printf("New balance after withdrawal: %d\n", account.balance)
	} else {
		fmt.Printf("Failed to withdraw %d: insufficient balance\n", amount)
	}
}

func main() {
	var wg sync.WaitGroup
	account := BankAccount{balance: 1000} // Initial balance of 1000

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	// simulating multiple customers depositing and withdrawing concurrently
	for i := 1; i <= 5; i++ {
		wg.Add(2) // Add 2 to the wait group for each iteration (1 deposit + 1 withdrawal)

		go func(customerID int) {
			// Simulate a deposit
			account.Deposit(100*customerID, &wg, r)
		}(i)

		go func(customerID int) {
			// Simulate a withdrawal
			account.Withdraw(50*customerID, &wg, r)
		}(i)
	}

	wg.Wait() // Wait for all goroutines to finish

	fmt.Printf("Final account balance: %d\n", account.balance)
}
