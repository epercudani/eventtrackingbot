package main

import "log"

func processUpdate(update Update) {

	switch update.Message.Text {
	case "/start":
		sendStartResponse(update)
	default:
		log.Printf("Other message received")
	}
}

func processUpdates(updateList []Update) (nextUpdateId uint64) {
	updatesCount := len(updateList)

	for i, update := range updateList {

		processUpdate(update)

		if i == (updatesCount - 1) {
			nextUpdateId = update.UpdateId + 1
		}
	}

	return
}