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
	REGEXP_SET_WHEN = regexp.MustCompile(fmt.Sprintf("^/set_when(@%s)?$", global.BOT_NAME))
	REGEXP_SET_WHERE = regexp.MustCompile(fmt.Sprintf("^/set_where(@%s)?$", global.BOT_NAME))
	REGEXP_DELETE_EVENT = regexp.MustCompile(fmt.Sprintf("^/delete_event(@%s)?$", global.BOT_NAME))
	REGEXP_CURRENT_EVENT = regexp.MustCompile(fmt.Sprintf("^/current_event(@%s)?$", global.BOT_NAME))
	REGEXP_SELECT_EVENT = regexp.MustCompile(fmt.Sprintf("^/select_event(@%s)?$", global.BOT_NAME))
	REGEXP_ALL_EVENTS = regexp.MustCompile(fmt.Sprintf("^/all_events(@%s)?$", global.BOT_NAME))
	REGEXP_PARTICIPANTS = regexp.MustCompile(fmt.Sprintf("^/participants(@%s)?$", global.BOT_NAME))
	REGEXP_WILL_GO = regexp.MustCompile(fmt.Sprintf("^/will_go(@%s)?$", global.BOT_NAME))
	REGEXP_WONT_GO = regexp.MustCompile(fmt.Sprintf("^/wont_go(@%s)?$", global.BOT_NAME))
)

func processUpdate(update types.Update) {

	if update.Message.Chat.ChatType == "private" {
		processPrivateChat(update)
		return
	}

	switch {
	case REGEXP_START.MatchString(strings.TrimSpace(update.Message.Text)):
		processStartOrHelp(update)

	case REGEXP_HELP.MatchString(strings.TrimSpace(update.Message.Text)):
		processStartOrHelp(update)

	case REGEXP_CREATE_EVENT.MatchString(strings.TrimSpace(update.Message.Text)):
		processCreateEventWithoutName(update)

	case REGEXP_SET_WHEN.MatchString(strings.TrimSpace(update.Message.Text)):
		processSetWhen(update)

	case REGEXP_SET_WHERE.MatchString(strings.TrimSpace(update.Message.Text)):
		processSetWhere(update)

	case REGEXP_DELETE_EVENT.MatchString(strings.TrimSpace(update.Message.Text)):
		processDeleteEventWithoutName(update)

	case REGEXP_CURRENT_EVENT.MatchString(strings.TrimSpace(update.Message.Text)):
		processCurrentEvent(update)

	case REGEXP_SELECT_EVENT.MatchString(strings.TrimSpace(update.Message.Text)):
		processSelectEvent(update)

	case REGEXP_ALL_EVENTS.MatchString(strings.TrimSpace(update.Message.Text)):
		processAllEvents(update)

	case REGEXP_PARTICIPANTS.MatchString(strings.TrimSpace(update.Message.Text)):
		processParticipants(update)

	case REGEXP_WILL_GO.MatchString(strings.TrimSpace(update.Message.Text)):
		processWillGo(update)

	case REGEXP_WONT_GO.MatchString(strings.TrimSpace(update.Message.Text)):
		processWontGo(update)

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
				case global.MESSAGE_TYPE_EVENT_PROVIDE_NAME:
					processSetEventName(update)
				case global.MESSAGE_TYPE_EVENT_PROVIDE_DATE:
					processSetEventProperty(update, global.EVENT_PROPERTY_DATE)
				case global.MESSAGE_TYPE_EVENT_PROVIDE_PLACE:
					processSetEventProperty(update, global.EVENT_PROPERTY_PLACE)
				case global.MESSAGE_TYPE_DELETE_EVENT_PROVIDE_INDEX:
					processIndexToDeleteEvent(update)
				case global.MESSAGE_TYPE_SELECT_CURRENT_EVENT:
					processIndexToSelectCurrentEvent(update)
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