package main

import (
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"fmt"
)

const (
	GET_UPDATES_DEFAULT_TIMEOUT = 0
	GET_UPDATES_DEFAULT_LIMIT = 100
)

type GetUpdatesResponse struct {
	Ok			bool		`json:"ok"`
	Result		[]Update	`json:"result"`
}

func getUpdates(offset uint64, timeout, limit int) (getUpdatesResponse GetUpdatesResponse, err error) {

	url := fmt.Sprintf(BASE_URL + "getUpdates?offset=%d&timeout=%d&limit=%d", offset, timeout, limit)
	log.Printf("url: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("%s", body)

	err = json.Unmarshal(body, &getUpdatesResponse)
	if err != nil {
		log.Fatal(err)
	}

	return
}