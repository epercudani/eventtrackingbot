package main

import (
	"log"
	"github.com/kinslayere/eventtrackingbot/requests"
	"github.com/kinslayere/eventtrackingbot/global"
	"github.com/kinslayere/eventtrackingbot/processes"
	"github.com/mediocregopher/radix.v2/redis"
)

func receiveAndProcessUpdates() {

	var nextUpdateId int64
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

	redisNetwork := "tcp"
	redisHost := "192.168.99.100:32768"
	var err error

	global.RedisClient, err = redis.Dial(redisNetwork, redisHost)
	if err != nil {
		log.Fatalf("Error connecting to redis host %s by %s: %v", redisHost, redisNetwork, err)
	}
	defer global.RedisClient.Close()

	receiveAndProcessUpdates()
}
