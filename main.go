package main

import (
	"demo/password/account"
	"demo/password/encrypter"
	"demo/password/files"
	"demo/password/output"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var menu = map[string]func(*account.VaultWithDb){
	"1": createAccount,
	"2": findAccountByUrl,
	"3": findAccountByLogin,
	"4": deleteAccount,
}

func main() {
	fmt.Println("______Passwords manager______")
	err := godotenv.Load()
	if err != nil {
		output.PrintError("Could not find the env file")
	}
	vault := account.NewVault(files.NewJsonDb("data.vault"), *encrypter.NewEncrypter())
Menu:
	for {
		option := promptData(
			"1. Create account",
			"2. Find account by URL",
			"3. Find account by Login",
			"4. Delete account",
			"5. Exit",
			"Choose option",
		)
		menuFunc := menu[option]
		if menuFunc == nil {
			break Menu
		}
		menuFunc(vault)
	}
}

func findAccountByUrl(vault *account.VaultWithDb) {
	url := promptData("Enter the search URL")
	accounts := vault.FindAccounts(url, func(acc account.Account, str string) bool {
		return strings.Contains(acc.Url, str)
	})
	outputRes(&accounts)
}

func findAccountByLogin(vault *account.VaultWithDb) {
	login := promptData("Enter the search Login")
	accounts := vault.FindAccounts(login, func(acc account.Account, str string) bool {
		return strings.Contains(acc.Login, str)
	})
	outputRes(&accounts)
}

func outputRes(accounts *[]account.Account) {
	if len(*accounts) == 0 {
		color.Red("There is no such account")
	}
	for _, account := range *accounts {
		account.Output()
	}
}

func deleteAccount(vault *account.VaultWithDb) {
	url := promptData("Enter the URL to delete")
	isDeleted := vault.DeleteAccountByUrl(url)
	if isDeleted {
		color.Green("Deleted")
	} else {
		output.PrintError("Not found")
	}
}

func createAccount(vault *account.VaultWithDb) {
	login := promptData("Enter login")
	password := promptData("Enter the password")
	url := promptData("Enter URL")
	myAcc, err := account.NewAccount(login, password, url)
	if err != nil {
		output.PrintError("Wrong URL or Login format")
		return
	}
	vault.AddAccount(*myAcc)
}

func promptData(prompt ...any) string {
	for i, line := range prompt {
		if i == len(prompt) - 1 {
			fmt.Printf("%v: ", line)
		} else {
			fmt.Println(line)
		}
	}
	var res string
	fmt.Scanln(&res)
	return res
}