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
		return WithdrawError{Amount: amount, Balance: acc.Balance}

	}

	acc.Balance -= amount
	return nil
}


type WithdrawError struct {
	Amount int
	Balance int
}

func (e WithdrawError) Error() string {
	smt := fmt.Sprintf("withdraw failed: balance= %d , attempted= %d", e.Balance, e.Amount)
	return smt
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