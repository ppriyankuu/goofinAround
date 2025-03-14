# Goroutines
A program that simulates a `bank account` where multiple `goroutines` will `deposit` and `withdraw` money concurrently. This might help us understand how to handle shared resources safely. A simple yet complex application in Go that demonstrates concurrency using goroutines and mutexes.

## Explanation 
- **Account Struct** : The Account struct has a balance field to store the account balance and a mutex field to synchronize access to the balance.
- **Deposit and Withdraw Methods** : These methods lock the mutex before modifying the balance to ensure that only one goroutine can modify the balance at a time.
- **SimulateTransactions Function** : This function spawns goroutines to simulate concurrent deposits and withdrawals. Each goroutine randomly chooses to deposit or withdraw a random amount.
- **WaitForTransactions Function** : This function waits for all transactions to complete by sleeping for a duration based on the number of transactions.
- **Main Function** : Initializes the account, starts the transactions, waits for them to complete, and prints the final balance.
