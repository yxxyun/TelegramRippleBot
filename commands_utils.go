package main

import (
	"net/http"
	"log"
	"io/ioutil"
	cmc "github.com/coincircle/go-coinmarketcap"
	"strings"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

func getRippleStats() string {
	coin, err := cmc.GetCoinData("Ripple")
	if err != nil {
		log.Print(err)
	}
	return "Price: " + float64ToString(coin.PriceUSD) + " USD\nVolume: " +
		float64ToString(coin.USD24HVolume) + " USD\nCapitalization: " +
		float64ToString(coin.MarketCapUSD) + " USD"
}

func checkAddress(a string) bool {
	if string(a[0]) != "r" {
		return false
	}
	if len(a) > 35 || len(a) < 25 {
		return false
	}
	return true
}

func getCurrency(name string) string {
	for _, listing := range listings {
		if strings.ToLower(listing.Name) == strings.ToLower(name) ||
			strings.ToLower(listing.Symbol) == strings.ToLower(name) {
			return listing.Name
		}
	}
	return ""
}

func getBalance(address string) float64 {
	resp, err := http.Get(config.RippleUrlBase + address + "/balances")
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}

	xrp := json.Get(bodyBytes, "balances", 0, "value").ToString()
	return stringToFloat64(xrp)
}

func rememberPost(message *tgbotapi.Message) {
	currPost = *message
	currState = "waitingForDelay"
	sendMessage(message.Chat.ID, phrases[10], nil)
}

func rememberDelay(message *tgbotapi.Message) {
	var posts []PendingPost
	readJson(&posts, "posts.json")

	post := PendingPost{Message: currPost,
		DelayHours: stringToFloat64(message.Text),
		PostTime: time.Now(),
	}
	posts = append(posts, post)

	writeJson(&posts, "posts.json")
	currPost = tgbotapi.Message{}
	currState = ""
	sendMessage(message.Chat.ID, phrases[11], nil)
}
