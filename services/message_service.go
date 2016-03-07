package services

import (
	"fmt"
	"log"
	"bytes"
	"github.com/kinslayere/eventtrackingbot/persistence"
	"github.com/kinslayere/eventtrackingbot/requests"
	"github.com/kinslayere/eventtrackingbot/types"
	"github.com/kinslayere/eventtrackingbot/global"
)

func SendSelectEventMessage(chatId, replyToMessageId, userFromId int64, events []types.Event) (err error) {

	var text bytes.Buffer

	text.WriteString(fmt.Sprintf("Please choose one of the following:"))
	for i, event := range events {
		text.WriteString(fmt.Sprintf("\n/%d %s", i + 1, GetEventDescription(event)))
	}

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(chatId)
	smr.AddText(text.String())
	smr.AddReplyToMessageId(replyToMessageId)
	smr.AddForceReply( types.ForceReply { ForceReply: true, Selective: true } )
	response, err := smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	messageSent := response.Result
	log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
	err = SetPendingResponseToMessage(userFromId, messageSent.MessageId, global.MESSAGE_TYPE_SELECT_CURRENT_EVENT)
	if err != nil {
		return err
	}

	return nil
}

func SendRequestPropertyMessage(chatId, replyToMessageId, userFromId int64, event types.Event, property string) (err error) {

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(chatId)
	smr.AddText(fmt.Sprintf("Please, provide \"%s\" %s.", event.Name, property))
	smr.AddReplyToMessageId(replyToMessageId)
	smr.AddForceReply( types.ForceReply { ForceReply: true, Selective: true } )
	response, err := smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	messageSent := response.Result
	log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
	err = SetPendingResponseToMessage(userFromId, messageSent.MessageId, global.PROPERTY_TO_MESSAGE_TYPE_MAP[property])
	if err != nil {
		return err
	}

	return nil
}

func SendCurrentEventMessage(chatId int64, currentEvent types.Event) (err error) {

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(chatId)
	smr.AddText(fmt.Sprintf("Your current event is \"%s\".", GetEventDetails(currentEvent)))
	_, err = smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	return nil
}

func SendCurrentEventNotSetMessage(chatId, replyToMessageId int64) (err error) {

	text := "There is no current event in this group."

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(chatId)
	smr.AddText(text)
	smr.AddReplyToMessageId(replyToMessageId)
	response, err := smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	messageSent := response.Result
	log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)

	return nil
}

func SendNoEventsInGroupMessage(chatId, replyToMessageId int64) (err error) {

	text := "There are no events in this group."

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(chatId)
	smr.AddText(text)
	smr.AddReplyToMessageId(replyToMessageId)
	response, err := smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	messageSent := response.Result
	log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)

	return nil
}

func SendTooManyParametersMessage(chatId int64) (err error) {

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(chatId)
	smr.AddText(fmt.Sprintf("Too many parameters. Please use: /set_when"))
	_, err = smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	return nil
}

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
