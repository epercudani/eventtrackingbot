package persistence

import (
	"github.com/kinslayere/eventtrackingbot/global"
	"strconv"
	"log"
)

func SaveInt(key string, value int) error {

	err := global.RedisClient.Cmd("SET", key, strconv.Itoa(value)).Err
	if err != nil {
		log.Printf("persistence.SaveInt. key=\"%s\" value=\"%d\". Error %v", key, value, err)
		return err
	}

	log.Printf("persistence.SaveInt. Saved key=\"%s\" value=\"%d\"", key, value)

	return nil
}

func SaveUint64(key string, value uint64) error {

	err := global.RedisClient.Cmd("SET", key, strconv.FormatUint(value, 10)).Err
	if err != nil {
		log.Printf("persistence.SaveUint64. key=\"%s\" value=\"%d\". Error %v", key, value, err)
		return err
	}

	log.Printf("persistence.SaveUint64. Saved key=\"%s\" value=\"%d\"", key, value)

	return nil
}

func SaveString(key string, value string) error {

	err := global.RedisClient.Cmd("SET", key, value).Err
	if err != nil {
		log.Printf("persistence.SaveString. key=\"%s\" value=\"%d\". Error %v", key, value, err)
		return err
	}

	log.Printf("persistence.SaveString. Saved key=\"%s\" value=\"%d\"", key, value)

	return nil
}

func GetString(key string) (string, error) {

	result, err := global.RedisClient.Cmd("GET", key).Str()
	if err != nil {
		log.Printf("persistence.GetString. key=\"%s\". Error %v", key, err)
		return "", err
	}

	return result, nil
}

func Exists(key string) (bool, error) {

	result, err := global.RedisClient.Cmd("EXISTS", key).Int()
	if err != nil {
		log.Printf("persistence.Exists. key=\"%s\". Error %v", key, err)
		return false, err
	}

	return result == 1, nil
}