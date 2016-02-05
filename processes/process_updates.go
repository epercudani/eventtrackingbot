package processes

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
	REGEXP_START = regexp.MustCompile(fmt.Sprintf("^/start(@%s)?$", global.BOT_NAME))
	REGEXP_HELP = regexp.MustCompile(fmt.Sprintf("^/help(@%s)?$", global.BOT_NAME))
	REGEXP_CREATE_EVENT = regexp.MustCompile(fmt.Sprintf("^/create_event(@%s)?$", global.BOT_NAME))
	REGEXP_DELETE_EVENT = regexp.MustCompile(fmt.Sprintf("^/delete_event(@%s)?$", global.BOT_NAME))
	REGEXP_CURRENT_EVENT = regexp.MustCompile(fmt.Sprintf("^/current_event(@%s)?$", global.BOT_NAME))
	REGEXP_ALL_EVENTS = regexp.MustCompile(fmt.Sprintf("^/all_events(@%s)?$", global.BOT_NAME))
)

func processUpdate(update types.Update) {

	if update.Message.Chat.ChatType == "private" {
		ProcessPrivateChat(update)
		return
	}

	switch {
	case REGEXP_START.MatchString(strings.TrimSpace(update.Message.Text)):
		processStartOrHelp(update)

	case REGEXP_HELP.MatchString(strings.TrimSpace(update.Message.Text)):
		processStartOrHelp(update)

	case REGEXP_CREATE_EVENT.MatchString(strings.TrimSpace(update.Message.Text)):
		processStartEventCreation(update)

	case REGEXP_DELETE_EVENT.MatchString(strings.TrimSpace(update.Message.Text)):
		processStartEventDeletion(update)

	case REGEXP_CURRENT_EVENT.MatchString(strings.TrimSpace(update.Message.Text)):
		processCurrentEvent(update)

	case REGEXP_ALL_EVENTS.MatchString(strings.TrimSpace(update.Message.Text)):
		processAllEvents(update)

	default:

		if update.Message.ReplyToMessage != nil {
			isResponseKey := fmt.Sprintf(persistence.KEY_WAITIING_RESPONSE_TO, update.Message.From.Id, update.Message.ReplyToMessage.MessageId)
			isResponse, err := persistence.Exists(isResponseKey)
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
				case global.MESSAGE_TYPE_CREATE_EVENT_PROVIDE_NAME:
					processEventName(update)
				case global.MESSAGE_TYPE_DELETE_EVENT_PROVIDE_INDEX:
					processDeleteIndex(update)
				}

				_, err = persistence.Delete(isResponseKey)
				if err != nil {
					log.Printf("Error deleting response check key: %v", err)
					return
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

func ProcessUpdates(updateList []types.Update) (nextUpdateId int64) {
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