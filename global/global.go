package global

import (
	"github.com/mediocregopher/radix.v2/pool"
	"fmt"
)

var RedisPool *pool.Pool
var BOT_NAME string
var BOT_TOKEN string

func GetBaseUrl() string {
	return fmt.Sprintf(TELEGRAM_BASE_URL, BOT_TOKEN)
}