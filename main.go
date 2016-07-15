package main

import (
	"log"
	"github.com/kinslayere/eventtrackingbot/requests"
	"github.com/kinslayere/eventtrackingbot/global"
	"github.com/kinslayere/eventtrackingbot/processes"
	"github.com/mediocregopher/radix.v2/redis"
	"syscall"
)

func receiveAndProcessUpdates() {

	var nextUpdateId int64
	log.Printf("Starting to listen updates")
	getUpdatesRequest := requests.NewGetUpdatesRequest()
	getUpdatesRequest.SetTimeout(60)
	getUpdatesRequest.SetLimit(global.GET_UPDATES_DEFAULT_LIMIT)

	for {
		getUpdatesRequest.SetOffset(nextUpdateId)
		updateResponse, err := getUpdatesRequest.Execute()
		if err != nil {
			log.Fatal(err)
		}

		if updateResponse.Ok {
			if len(updateResponse.Result) > 0 {
				nextUpdateId = processes.ProcessUpdates(updateResponse.Result)
			}
		} else {
			log.Fatal("updateResponse.Ok => false")
		}
	}
}

func main() {

	log.Printf("Starting bot")

	redisNetwork := "tcp"
	redisHost, ok := syscall.Getenv("REDIS_HOST")
	if !ok {
		log.Fatalf("Failed to obtain redis host address from env variable %s", "REDIS_HOST")
	}
	var err error

	global.RedisClient, err = redis.Dial(redisNetwork, redisHost)
	if err != nil {
		log.Fatalf("Error connecting to redis host %s by %s: %v", redisHost, redisNetwork, err)
	}
	defer global.RedisClient.Close()

	receiveAndProcessUpdates()

	log.Printf("Starting bot")
}
