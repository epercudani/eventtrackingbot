package main

import (
	"log"
	"github.com/kinslayere/eventtrackingbot/requests"
	"github.com/kinslayere/eventtrackingbot/global"
	"github.com/kinslayere/eventtrackingbot/processes"
	"github.com/mediocregopher/radix.v2/redis"
	"syscall"
	"net/url"
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
	redisUrlRaw, ok := syscall.Getenv("REDIS_URL")
	if !ok {
		log.Fatalf("Failed to obtain redis host address from env variable %s", "REDIS_URL")
	}
	var err error

	redisUrl, err := url.Parse(redisUrlRaw)
	if err != nil {
		log.Fatalf("Failed to parse redis URL: %s. Error is: %s", redisUrlRaw, err.Error())
	}

	global.RedisClient, err = redis.Dial(redisNetwork, redisUrl.Path)
	if err != nil {
		log.Fatalf("Error connecting to redis host %s by %s. Error is: %v", redisUrl, redisNetwork, err)
	}
	defer global.RedisClient.Close()

	if redisUrl.User != nil {
		redisPass, ok := redisUrl.User.Password()
		if ok {
			global.RedisClient.Cmd("AUTH", redisPass)
		}
	}

	receiveAndProcessUpdates()

	log.Printf("Starting bot")
}
