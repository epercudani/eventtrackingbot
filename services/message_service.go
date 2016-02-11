package services

import (
	"github.com/kinslayere/eventtrackingbot/persistence"
	"fmt"
	"log"
)

func SetPendingResponseToMessage(userId, messageId int64, messageType string) (err error) {

	err = persistence.SaveIntWithTTL(fmt.Sprintf(persistence.KEY_WAITIING_RESPONSE_TO, userId, messageId), 1, 10 * 60) // Ten minutes
	if err != nil {
		log.Printf("Error saving waiting response flag: %v", err)
		return
	}

	err = persistence.SaveStringWithTTL(fmt.Sprintf(persistence.KEY_MESSAGE_TYPE, messageId), messageType, 10 * 60) // Ten minutes
	if err != nil {
		log.Printf("Error saving message type: %v", err)
		return
	}

	return
}
