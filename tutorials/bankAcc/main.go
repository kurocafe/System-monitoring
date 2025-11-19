package main

import "fmt"

type Account struct {
	Owner string
	Balance int
}

func (acc *Account) Deposit(amount int){
	acc.Balance += amount
}

func (acc *Account) Withdraw(amount int) error {
	if acc.Balance - amount < 0{
		return fmt.Errorf("insufficient funds: balance=%d, attempted=%d", acc.Balance, amount)

	}

	acc.Balance -= amount
	return nil
}

func main(){
	acc := Account{
    Owner:   "nachi",
    Balance: 3000,
	}
	acc.Deposit(2000)
	fmt.Println("balance:", acc.Balance)
	err := acc.Withdraw(2000)

	if err != nil {
		 fmt.Println("withdraw error:", err)
		return
	}

	fmt.Printf("balance: %d \n", acc.Balance)

	err = acc.Withdraw(5000)
	if err != nil {
		 fmt.Println("withdraw error:", err)
	}

	fmt.Println("final balance:", acc.Balance)
	
	
}