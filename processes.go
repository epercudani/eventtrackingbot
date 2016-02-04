package main

import (
	"fmt"
	"log"
	"strings"
	"github.com/kinslayere/eventtrackingbot/types"
	"github.com/kinslayere/eventtrackingbot/requests"
	"github.com/kinslayere/eventtrackingbot/persistence"
	"github.com/kinslayere/eventtrackingbot/global"
)

func ProcessCreateEvent(update types.Update) {

	fields := strings.Fields(strings.TrimSpace(update.Message.Text))
	smr := requests.NewSendMessageRequest()

	switch len(fields) {
	case 1:
		smr.AddChatIdInt(update.Message.Chat.Id)
		smr.AddText(fmt.Sprintf("Hi %s! I'll help to create your event. Please tell me the event's name.", update.Message.Chat.FirstName))
		smr.AddReplyToMessageId(update.Message.MessageId)
		smr.AddForceReply( types.ForceReply { ForceReply: true, Selective: true } )
		messageSent, err := smr.DoRequest()
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return
		}

		err = persistence.SaveInt(fmt.Sprintf(persistence.KEY_WAITIING_RESPONSE_TO, update.Message.From.Id, messageSent.MessageId), 1)
		if err != nil {
			log.Printf("Error saving waiting response flag: %v", err)
			return
		}

		err = persistence.SaveString(fmt.Sprintf(persistence.KEY_MESSAGE_TYPE, messageSent.MessageId), global.MESSAGE_TYPE_CREATE_EVENT)
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

func ProcessResponseToCreateEvent(update types.Update)  {
	log.Printf("Processing response to create_event")
}

func ProcessCurrentEvent(update types.Update) {
	log.Printf("Processing current event")
}

func processStartOrHelp(update types.Update) {

	text := fmt.Sprintf("Hello %s!\n", update.Message.Chat.FirstName)
	text += "This is an event tracking attendance bot.\nAvailable commands are:"

	smr := requests.NewSendMessageRequest()
	smr.AddChatIdInt(update.Message.Chat.Id)
	smr.AddText(text)
	smr.DoRequest()
}