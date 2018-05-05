package main

import (
	"log"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

func addUserIfAbsent(user *tgbotapi.User) {
	log.Print("Started adding user")
	var u User
	db.Find(&u, "id = ?", user.ID)
	if u.ID == (User{}.ID) {
		userInsert := User{ID: int64(user.ID), FirstName: user.FirstName, LastName: user.LastName, UserName: user.UserName}
		db.Create(&userInsert)
		log.Printf("Added new user with id = %v", userInsert.ID)
	}
}

func addWalletDB(message *tgbotapi.Message) {
	log.Print("Started adding wallet")
	var addr, name string
	fields := strings.Fields(message.Text)
	addr = fields[1]
	if len(fields) == 3{
		name = fields[2]
	}

	u := getUser(message.From.ID)
	var wallet Wallet
	db.Find(&wallet, "address = ?", addr)
	if wallet.ID == 0{
		db.Model(&u).Association("Wallets").Append(Wallet{Address:addr})
		db.First(&wallet, "address = ?", addr)
	}else {
		db.Model(&u).Association("Wallets").Append(wallet)
	}

	if name != ""{
		db.Model(&UserWallet{}).Where("user_id = ? AND wallet_id = ?", u.ID, wallet.ID).Update("name", name)
	}
}

func removeWalletDB(message *tgbotapi.Message) {
	log.Print("Started removing wallet")
	addr := strings.Fields(message.Text)[1]
	u := getUser(message.From.ID)
	var wallet Wallet
	db.Find(&wallet, "address = ?", addr)
	db.Model(&u).Association("Wallets").Delete(wallet)
	if db.Model(&wallet).Association("Users").Count() == 0{
		db.Delete(&wallet)
	}
}

func getUser(id int) User {
	log.Print("Started get user")
	var u User
	db.Find(&u, "id = ?", id)
	return u
}
