package main

import (
	"log"
	"github.com/kinslayere/eventtrackingbot/requests"
	"github.com/kinslayere/eventtrackingbot/global"
	"github.com/kinslayere/eventtrackingbot/processes"
	"syscall"
	"net/url"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
	"time"
)

func receiveAndProcessUpdates() {

	var nextUpdateId int64
	log.Println("Starting to listen updates")
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

	log.Println("Starting bot")

	redisNetwork := "tcp"
	redisUrlRaw, ok := syscall.Getenv("REDIS_URL")
	if !ok {
		log.Fatalf("Failed to obtain redis host address from env variable %s", "REDIS_URL")
	}
	var err error

	log.Printf("Redis raw URL: %s\n", redisUrlRaw)
	redisUrl, err := url.Parse(redisUrlRaw)
	if err != nil {
		log.Fatalf("Failed to parse redis URL: %s. Error is: %s", redisUrlRaw, err.Error())
	}

	log.Printf("Redis URL host: %s\n", redisUrl.Host)

	var redisPass string
	var passSet bool
	if redisUrl.User != nil {
		redisPass, passSet = redisUrl.User.Password()
	}

	connectFunc := pool.DialFunc(func(network, addr string) (*redis.Client, error) {

		client, err := redis.DialTimeout(network, addr, time.Duration(5 * time.Second))
		if err != nil {
			return nil, err

		}

		if passSet {
			if err = client.Cmd("AUTH", redisPass).Err; err != nil {
				client.Close()
				return nil, err
			}
		}

		return client, nil
	})

	global.RedisPool, err = pool.NewCustom(redisNetwork, redisUrl.Host, 1, connectFunc)
	if err != nil {
		log.Fatalf("Error connecting to redis host %s by %s. Error is: %v", redisUrl.Host, redisNetwork, err)
	}

	log.Println("Bot started")

	receiveAndProcessUpdates()

	log.Println("Finishing bot")
}
