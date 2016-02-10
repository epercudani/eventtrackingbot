package processes

import (
	"strings"
	"fmt"
	"log"
	"github.com/kinslayere/eventtrackingbot/global"
	"github.com/kinslayere/eventtrackingbot/persistence"
	"github.com/kinslayere/eventtrackingbot/requests"
	"github.com/kinslayere/eventtrackingbot/types"
	"github.com/kinslayere/eventtrackingbot/services"
	"bytes"
	"regexp"
	"strconv"
)

func processCreateEventWithoutName(update types.Update) {

	fields := strings.Fields(strings.TrimSpace(update.Message.Text))
	smr := requests.NewSendMessageRequest()

	switch len(fields) {
	case 1:
		smr.AddChatIdInt(update.Message.Chat.Id)
		smr.AddText(fmt.Sprintf("Hi %s! I'll help to create your event. Please tell me the event's name.", update.Message.From.FirstName))
		smr.AddReplyToMessageId(update.Message.MessageId)
		smr.AddForceReply( types.ForceReply { ForceReply: true, Selective: true } )
		messageSent, err := smr.DoRequest()
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return
		}

		err = persistence.SaveIntWithTTL(fmt.Sprintf(persistence.KEY_WAITIING_RESPONSE_TO, update.Message.From.Id, messageSent.MessageId), 1, 10 * 60) // Ten minutes
		if err != nil {
			log.Printf("Error saving waiting response flag: %v", err)
			return
		}

		err = persistence.SaveStringWithTTL(fmt.Sprintf(persistence.KEY_MESSAGE_TYPE, messageSent.MessageId), global.MESSAGE_TYPE_CREATE_EVENT_PROVIDE_NAME, 10 * 60) // Ten minutes
		if err != nil {
			log.Printf("Error saving message type: %v", err)
			return
		}

		log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
	default:
		smr.AddChatIdInt(update.Message.Chat.Id)
		smr.AddText(fmt.Sprintf("Too many parameters. Please use: /create_event"))
	}
}

func processDeleteEventWithoutName(update types.Update)  {

	fields := strings.Fields(strings.TrimSpace(update.Message.Text))
	smr := requests.NewSendMessageRequest()

	switch len(fields) {
	case 1:

		var text bytes.Buffer

		events := services.GetGroupEventNames(update.Message.Chat.Id)
		if len(events) > 0 {
			text.WriteString(fmt.Sprintf("Hi %s! Which event do you want to delete?", update.Message.From.FirstName))
			for i, event := range events {
				text.WriteString(fmt.Sprintf("\n/%d %s", i + 1, event))
			}
		} else {
			text.WriteString(fmt.Sprintf("There are no events for this group."))
		}

		smr.AddChatIdInt(update.Message.Chat.Id)
		smr.AddText(text.String())
		smr.AddReplyToMessageId(update.Message.MessageId)
		smr.AddForceReply( types.ForceReply { ForceReply: true, Selective: true } )
		messageSent, err := smr.DoRequest()
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return
		}

		err = persistence.SaveIntWithTTL(fmt.Sprintf(persistence.KEY_WAITIING_RESPONSE_TO, update.Message.From.Id, messageSent.MessageId), 1, 10 * 60) // Ten minutes
		if err != nil {
			log.Printf("Error saving waiting response flag: %v", err)
			return
		}

		err = persistence.SaveStringWithTTL(fmt.Sprintf(persistence.KEY_MESSAGE_TYPE, messageSent.MessageId), global.MESSAGE_TYPE_DELETE_EVENT_PROVIDE_INDEX, 10 * 60) // Ten minutes
		if err != nil {
			log.Printf("Error saving message type: %v", err)
			return
		}

		log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
	default:
		smr.AddChatIdInt(update.Message.Chat.Id)
		smr.AddText(fmt.Sprintf("Too many parameters. Please use: /delete_event"))
	}
}

func processSetEventName(update types.Update)  {

	log.Printf("Processing event name")

	eventName := update.Message.Text
	if len(strings.TrimSpace(eventName)) > 0 {
		services.CreateEvent(update.Message.Chat.Id, eventName)
	}
}

func processIndexToDeleteEvent(update types.Update) {

	log.Printf("Processing delete index")

	regexpIndex := regexp.MustCompile("^/[0-9]+$")

	deleted := false

	if regexpIndex.MatchString(update.Message.Text) {

		index, err := strconv.Atoi(update.Message.Text[1:])
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		size, err := persistence.GetSortedSetSize(fmt.Sprintf(persistence.KEY_GROUP_EVENTS, update.Message.Chat.Id))
		if err != nil {
			log.Printf("Error getting group events count: %v", err)
		}

		if index >= 1 && index <= size {
			services.DeleteEvent(update.Message.Chat.Id, index - 1)
			deleted = true
		}
	}

	if !deleted {
		smr := requests.NewSendMessageRequest()
		smr.AddChatIdInt(update.Message.Chat.Id)
		smr.AddText(fmt.Sprintf("%s is not a valid option.", update.Message.Text))
		smr.DoRequest()
	}
}

func processIndexToSelectCurrentEvent(update types.Update) {

	log.Printf("Processing select index")

	regexpIndex := regexp.MustCompile("^/[0-9]+$")

	selected := false

	if regexpIndex.MatchString(update.Message.Text) {

		index, err := strconv.Atoi(update.Message.Text[1:])
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		size, err := persistence.GetSortedSetSize(fmt.Sprintf(persistence.KEY_GROUP_EVENTS, update.Message.Chat.Id))
		if err != nil {
			log.Printf("Error getting group events count: %v", err)
		}

		if index >= 1 && index <= size {

			groupId := update.Message.Chat.Id
			eventName := services.GetGroupEventNames(groupId)[index - 1]
			services.SetCurrentEvent(update.Message.Chat.Id, fmt.Sprintf(persistence.KEY_EVENT, groupId, eventName))
			selected = true
		}
	}

	if !selected {
		smr := requests.NewSendMessageRequest()
		smr.AddChatIdInt(update.Message.Chat.Id)
		smr.AddText(fmt.Sprintf("%s is not a valid option.", update.Message.Text))
		smr.DoRequest()
	}
}

func processAllEvents(update types.Update)  {

	events := services.GetGroupEventNames(update.Message.Chat.Id)

	var text bytes.Buffer
	if len(events) > 0 {
		text.WriteString(fmt.Sprintf("Current events in this group are:"))
		for _, event := range events {
			text.WriteString(fmt.Sprintf("\n%s", event))
		}
	} else {
		text.WriteString(fmt.Sprintf("There are no events in this group."))
	}

	smr := requests.NewSendMessageRequest()
	smr.AddChatIdInt(update.Message.Chat.Id)
	smr.AddText(text.String())
	smr.DoRequest()
}

func processCurrentEvent(update types.Update) {

	var text bytes.Buffer
	selectPending := false

	eventName := services.GetCurrentEventName(update.Message.Chat.Id)
	if len(eventName) > 0 {
		text.WriteString(fmt.Sprintf("Your current event is \"%s\".", eventName))
	} else {
		events := services.GetGroupEventNames(update.Message.Chat.Id)
		if len(events) > 0 {
			text.WriteString(fmt.Sprintf("There is no event currently selected. Please chose one of the following:"))
			for i, event := range events {
				text.WriteString(fmt.Sprintf("\n/%d %s", i + 1, event))
			}

			selectPending = true
		} else {
			text.WriteString(fmt.Sprintf("There are no events in this group."))
		}
	}

	smr := requests.NewSendMessageRequest()
	smr.AddChatIdInt(update.Message.Chat.Id)
	smr.AddText(text.String())
	if selectPending {
		smr.AddReplyToMessageId(update.Message.MessageId)
		smr.AddForceReply( types.ForceReply { ForceReply: true, Selective: true } )
		messageSent, err := smr.DoRequest()
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return
		}

		err = persistence.SaveIntWithTTL(fmt.Sprintf(persistence.KEY_WAITIING_RESPONSE_TO, update.Message.From.Id, messageSent.MessageId), 1, 10 * 60) // Ten minutes
		if err != nil {
			log.Printf("Error saving waiting response flag: %v", err)
			return
		}

		err = persistence.SaveStringWithTTL(fmt.Sprintf(persistence.KEY_MESSAGE_TYPE, messageSent.MessageId), global.MESSAGE_TYPE_SELECT_CURRENT_EVENT, 10 * 60) // Ten minutes
		if err != nil {
			log.Printf("Error saving message type: %v", err)
			return
		}

		log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
	} else {
		smr.DoRequest()
	}
}

func processSelectEvent(update types.Update) {

	var text bytes.Buffer

	events := services.GetGroupEventNames(update.Message.Chat.Id)
	if len(events) > 0 {
		text.WriteString(fmt.Sprintf("There is no event currently selected. Please chose one of the following:"))
		for i, event := range events {
			text.WriteString(fmt.Sprintf("\n/%d %s", i + 1, event))
		}
	} else {
		text.WriteString(fmt.Sprintf("There are no events in this group."))
	}

	smr := requests.NewSendMessageRequest()
	smr.AddChatIdInt(update.Message.Chat.Id)
	smr.AddText(text.String())
	smr.AddReplyToMessageId(update.Message.MessageId)
	smr.AddForceReply( types.ForceReply { ForceReply: true, Selective: true } )
	messageSent, err := smr.DoRequest()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return
	}

	err = persistence.SaveIntWithTTL(fmt.Sprintf(persistence.KEY_WAITIING_RESPONSE_TO, update.Message.From.Id, messageSent.MessageId), 1, 10 * 60) // Ten minutes
	if err != nil {
		log.Printf("Error saving waiting response flag: %v", err)
		return
	}

	err = persistence.SaveStringWithTTL(fmt.Sprintf(persistence.KEY_MESSAGE_TYPE, messageSent.MessageId), global.MESSAGE_TYPE_SELECT_CURRENT_EVENT, 10 * 60) // Ten minutes
	if err != nil {
		log.Printf("Error saving message type: %v", err)
		return
	}

	log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
}