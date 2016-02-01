package main

import (
	"net/http"
	"log"
	"fmt"
	"net/url"
)

func sendStartResponse(update Update) {

	text := fmt.Sprintf("Hello %s!\n", update.Message.Chat.FirstName)
	text += "This is an event tracking attendance bot.\nAvailable commands are:"

	url := fmt.Sprintf(BASE_URL + "sendMessage?chat_id=%d&text=%s", update.Message.Chat.Id, url.QueryEscape(text))
	log.Printf("url: %s", url)

	_, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}

	return
}
