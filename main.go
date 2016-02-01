package main

import (
	"log"
)

const (
	TOKEN = "133198388:AAHHbnm7cNHMEF6hmdehKTCMRDrFMN46n-U"
	BASE_URL = "https://api.telegram.org/bot" + TOKEN + "/"
)

func receiveAndProcessUpdates(nextUpdateId uint64) {

	for {
		updateResponse, err := getUpdates(nextUpdateId, 60, GET_UPDATES_DEFAULT_LIMIT)
		if err != nil {
			log.Fatal(err)
		}

		if updateResponse.Ok {
			if len(updateResponse.Result) > 0 {
				nextUpdateId = processUpdates(updateResponse.Result)
			}
		} else {
			log.Fatal("updateResponse.Ok => false")
		}
	}
}

func main() {

	var nextUpdateId uint64

	receiveAndProcessUpdates(nextUpdateId)
}
