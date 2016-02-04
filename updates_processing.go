package main

import (
	"log"
	"regexp"
	"github.com/kinslayere/eventtrackingbot/types"
	"strings"
	"github.com/kinslayere/eventtrackingbot/persistence"
	"fmt"
	"github.com/kinslayere/eventtrackingbot/global"
)

var (
	REGEXP_START = regexp.MustCompile("^/start$")
	REGEXP_HELP = regexp.MustCompile("^/help$")
	REGEXP_CREATE_EVENT = regexp.MustCompile("^/create_event$")
	REGEXP_CURRENT_EVENT = regexp.MustCompile("^/current_event$")
)

func processUpdate(update types.Update) {

	switch {
	case REGEXP_START.MatchString(strings.TrimSpace(update.Message.Text)):
		processStartOrHelp(update)
	case REGEXP_HELP.MatchString(strings.TrimSpace(update.Message.Text)):
		processStartOrHelp(update)
	case REGEXP_CREATE_EVENT.MatchString(strings.TrimSpace(update.Message.Text)):
		ProcessCreateEvent(update)
	case REGEXP_CURRENT_EVENT.MatchString(strings.TrimSpace(update.Message.Text)):
		ProcessCurrentEvent(update)
	default:

		if update.Message.ReplyToMessage != nil {
			isResponse, err := persistence.Exists(fmt.Sprintf(persistence.KEY_WAITIING_RESPONSE_TO, update.Message.From.Id, update.Message.ReplyToMessage.MessageId))
			if err != nil {
				log.Printf("Error checking if message is response: %v", err)
				return
			}

			if isResponse {

				log.Printf("Response received")

				messageType, err := persistence.GetString(fmt.Sprintf(persistence.KEY_MESSAGE_TYPE, update.Message.ReplyToMessage.MessageId))
				if err != nil {
					log.Printf("Error getting message type: %v", err)
					return
				}

				switch messageType {
				case global.MESSAGE_TYPE_CREATE_EVENT:
					ProcessResponseToCreateEvent(update)
				}
			}
		} else {
			log.Printf("Other message received")
		}
	}

/*
/change_event [name]
/close_event [name]
/will_go [name]
/wont_go [name]
*/
}

func ProcessUpdates(updateList []types.Update) (nextUpdateId uint64) {
	updatesCount := len(updateList)

	for i, update := range updateList {

		processUpdate(update)

		// If it's the last update, return the next id
		if i == (updatesCount - 1) {
			nextUpdateId = update.UpdateId + 1
		}
	}

	return
}