package persistence

import (
	"strconv"
	"log"
	"github.com/kinslayere/eventtrackingbot/global"
	"github.com/mediocregopher/radix.v2/redis"
)

func doWithRetry(cmd string, args ...interface{}) *redis.Resp {

	var resp *redis.Resp
	for i := 0; i < global.PERSISTENCE_RETRIES; i++ {

		resp = global.RedisPool.Cmd(cmd, args)

		if resp.Err != nil {
			log.Printf("persistence.doWithRetry. cmd=\"%s\" args=\"%v\". Attempt %d Error %v", cmd, args, i+1, resp.Err)
			continue
		}

		break
	}

	return resp
}

func SaveInt(key string, value int) error {

	err := doWithRetry("SET", key, strconv.Itoa(value)).Err
	if err != nil {
		log.Printf("persistence.SaveInt. key=\"%s\" value=\"%d\". Error %v", key, value, err)
		return err
	}

	log.Printf("persistence.SaveInt. Saved key=\"%s\" value=\"%d\"", key, value)

	return nil
}

func SaveIntWithTTL(key string, value, ttl int) error {

	err := doWithRetry("SET", key, strconv.Itoa(value), "EX", ttl).Err
	if err != nil {
		log.Printf("persistence.SaveInt. key=\"%s\" value=\"%d\". Error %v", key, value, err)
		return err
	}

	log.Printf("persistence.SaveInt. Saved key=\"%s\" value=\"%d\"", key, value)

	return nil
}

func SaveUint64(key string, value uint64) error {

	err := doWithRetry("SET", key, strconv.FormatUint(value, 10)).Err
	if err != nil {
		log.Printf("persistence.SaveUint64. key=\"%s\" value=\"%d\". Error %v", key, value, err)
		return err
	}

	log.Printf("persistence.SaveUint64. Saved key=\"%s\" value=\"%d\"", key, value)

	return nil
}

func SaveString(key string, value string) error {

	err := doWithRetry("SET", key, value).Err
	if err != nil {
		log.Printf("persistence.SaveString. key=\"%s\" value=\"%s\". Error %v", key, value, err)
		return err
	}

	log.Printf("persistence.SaveString. Saved key=\"%s\" value=\"%s\"", key, value)

	return nil
}

func SaveStringWithTTL(key, value string, ttl int) error {

	err := doWithRetry("SET", key, value, "EX", ttl).Err
	if err != nil {
		log.Printf("persistence.SaveString. key=\"%s\" value=\"%d\". Error %v", key, value, err)
		return err
	}

	log.Printf("persistence.SaveString. Saved key=\"%s\" value=\"%s\"", key, value)

	return nil
}

func AddStringToSet(key, value string) error {

	err := doWithRetry("SADD", key, value).Err
	if err != nil {
		log.Printf("persistence.AddStringToSet. set=\"%s\" value=\"%s\". Error %v", key, value, err)
		return err
	}

	log.Printf("persistence.AddStringToSet. Saved set=\"%s\" value=\"%s\"", key, value)

	return nil
}

func AddStringToList(key, value string) error {

	err := doWithRetry("RPUSH", key, value).Err
	if err != nil {
		log.Printf("persistence.AddStringToList. list=\"%s\" value=\"%s\". Error %v", key, value, err)
		return err
	}

	log.Printf("persistence.AddStringToList. Saved list=\"%s\" value=\"%s\"", key, value)

	return nil
}

func AddStringToSortedSet(key string, score int, value string) error {

	err := doWithRetry("ZADD", key, score, value).Err
	if err != nil {
		log.Printf("persistence.AddStringToSortedSet. set=\"%s\" score=\"%d\" value=\"%s\". Error %v", key, score, value, err)
		return err
	}

	log.Printf("persistence.AddStringToSortedSet. Saved set=\"%s\" score=\"%d\" value=\"%s\"", key, score, value)

	return nil
}

func AddStringFieldToHash(hashKey, key, value string) error {

	err := doWithRetry("HSET", hashKey, key, value).Err
	if err != nil {
		log.Printf("persistence.AddStringFieldToHash. hash=\"%s\" key=\"%s\" value=\"%s\". Error %v", hashKey, key, value, err)
		return err
	}

	log.Printf("persistence.AddStringFieldToHash. Saved hash=\"%s\" key=\"%s\" value=\"%s\"", hashKey, key, value)

	return nil
}

func RemoveStringFromList(key, value string) error {

	err := doWithRetry("LREM", key, 0, value).Err
	if err != nil {
		log.Printf("persistence.RemoveStringFromList. list=\"%s\" value=\"%s\". Error %v", key, value, err)
		return err
	}

	log.Printf("persistence.RemoveStringFromList. Removed list=\"%s\" value=\"%s\"", key, value)

	return nil
}

func RemoveFromSet(setKey, key string) error {

	err := doWithRetry("SREM", setKey, key).Err
	if err != nil {
		log.Printf("persistence.RemoveFromSet. set=\"%s\" key=\"%s\". Error %v", setKey, key, err)
		return err
	}

	log.Printf("persistence.RemoveFromSet. Removed set=\"%s\" key=\"%s\"", setKey, key)

	return nil
}

func RemoveFromSortedSet(setKey, key string) error {

	err := doWithRetry("ZREM", setKey, key).Err
	if err != nil {
		log.Printf("persistence.RemoveFromSortedSet. sset=\"%s\" key=\"%s\". Error %v", setKey, key, err)
		return err
	}

	log.Printf("persistence.RemoveFromSortedSet. Saved set=\"%s\" key=\"%s\"", setKey, key)

	return nil
}

func RemoveFromSortedSetByScore(setKey string, scoreMin, scoreMax int) error {

	err := doWithRetry("ZREMRANGEBYSCORE", setKey, scoreMin, scoreMax).Err
	if err != nil {
		log.Printf("persistence.RemoveFromSortedSetByScore. sset=\"%s\" scoreMin=\"%d\" scoreMax=\"%d\" value=\"%s\". Error %v", setKey, scoreMin, scoreMax, err)
		return err
	}

	log.Printf("persistence.AddStringToOrderedSet. Saved set=\"%s\" scoreMin=\"%d\" scoreMax=\"%d\" value=\"%s\"", setKey, scoreMin, scoreMax)

	return nil
}

func GetFullHash(hashKey string) ([]string, error) {

	result, err := doWithRetry("HGETALL", hashKey).List()
	if err != nil {
		log.Printf("persistence.getFullHash. hash=\"%s\". Error %v", hashKey, err)
		return nil, err
	}

	return result, nil
}

func GetStringFieldFromHash(hashKey, key string) (string, error) {

	result, err := doWithRetry("HGET", hashKey, key).Str()
	if err != nil {
		log.Printf("persistence.GetStringFieldFromHash. hash=\"%s\" key=\"%s\". Error %v", hashKey, key, err)
		return "", err
	}

	return result, nil
}

func GetString(key string) (string, error) {

	result, err := doWithRetry("GET", key).Str()
	if err != nil {
		log.Printf("persistence.GetString. key=\"%s\". Error %v", key, err)
		return "", err
	}

	return result, nil
}

func GetInt(key string) (int, error) {

	result, err := doWithRetry("GET", key).Int()
	if err != nil {
		log.Printf("persistence.GetString. key=\"%s\". Error %v", key, err)
		return "", err
	}

	return result, nil
}

func GetStringsFromSet(key string) ([]string, error) {

	result, err := doWithRetry("SMEMBERS", key).List()
	if err != nil {
		log.Printf("persistence.GetStringsFromSet. set=\"%s\". Error %v", key, err)
		return nil, err
	}

	return result, nil
}

func GetStringsFromSortedSet(key string) ([]string, error) {

	result, err := doWithRetry("ZRANGE", key, 0, -1).List()
	if err != nil {
		log.Printf("persistence.GetStringsFromSortedSet. sset=\"%s\". Error %v", key, err)
		return nil, err
	}

	return result, nil
}

func GetStringsFromList(key string) ([]string, error) {

	result, err := doWithRetry("LRANGE", key, 0, -1).List()
	if err != nil {
		log.Printf("persistence.GetStringsFromList. list=\"%s\". Error %v", key, err)
		return nil, err
	}

	return result, nil
}

func GetListSize(key string) (int, error) {

	result, err := doWithRetry("LLEN", key).Int()
	if err != nil {
		log.Printf("persistence.GetListSize. list=\"%s\". Error %v", key, err)
		return -1, err
	}

	return result, nil
}

func GetSortedSetSize(key string) (int, error) {

	result, err := doWithRetry("ZCARD", key).Int()
	if err != nil {
		log.Printf("persistence.GetSortedSetSize. sset=\"%s\". Error %v", key, err)
		return -1, err
	}

	return result, nil
}

func Exists(key string) (bool, error) {

	result, err := doWithRetry("EXISTS", key).Int()
	if err != nil {
		log.Printf("persistence.Exists. key=\"%s\". Error %v", key, err)
		return false, err
	}

	return result == 1, nil
}

func Delete(key string) (bool, error) {

	result, err := doWithRetry("DEL", key).Int()
	if err != nil {
		log.Printf("persistence.Delete. key=\"%s\". Error %v", key, err)
		return false, err
	}

	return result == 1, nil
}

