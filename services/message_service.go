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

func SendResponseToCreateEventMessage(chatId, messageId, userId int64, userName string) (err error) {

	// Send acknowledge message and wait for event name
	messageText := fmt.Sprintf("Hi %s! I'll help you to create your event. Please tell me the event's name.", userName)
	smr := newMessageReplyRequest(chatId, messageId, messageText)
	response, err := smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	messageSent := response.Result
	log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
	err = SetPendingResponseToMessage(userId, messageSent.MessageId, global.MESSAGE_TYPE_EVENT_PROVIDE_NAME)
	if err != nil {
		return err
	}

	return nil
}

func SendResponseToDeleteEventMessage(chatId, messageId, userId int64, userName string) (err error) {

	var messageText bytes.Buffer

	events := GetGroupEventNames(chatId)
	if len(events) > 0 {
		messageText.WriteString(fmt.Sprintf("Hi %s! Which event do you want to delete?", userName))
		for i, event := range events {
			messageText.WriteString(fmt.Sprintf("\n/%d %s", i + 1, event))
		}
	} else {
		messageText.WriteString(fmt.Sprint("There are no events in this group."))
	}

	smr := newMessageReplyRequest(chatId, messageId, messageText.String())
	response, err := smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	messageSent := response.Result
	log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
	err = SetPendingResponseToMessage(userId, messageSent.MessageId, global.MESSAGE_TYPE_DELETE_EVENT_PROVIDE_INDEX)
	if err != nil {
		return err
	}

	return nil
}

func SendCreatedEventAcknowledge(chatId int64, userName, eventName string) error {

	messageText := fmt.Sprintf("Congrats %s! Event \"%s\" was created and set as the current event in this group.", userName, eventName)

	return sendSimpleMessage(chatId, messageText)
}

func SendDeletedEventAcknowledge(chatId int64, userName, eventName, selectedOption string, deleted bool) error {

	var messageText string

	if deleted {
		messageText = fmt.Sprintf("%s, you have successfully deleted event \"%s\".", userName, eventName)
	} else {
		messageText = fmt.Sprintf("%s is not a valid option.", selectedOption)
	}

	return sendSimpleMessage(chatId, messageText)
}

func SendSelectedCurrentEventAcknowledge(chatId int64, eventName, selectedOption string, selected bool) error {

	var messageText string

	if selected {
		messageText = fmt.Sprintf("%s is now set as the current event in this group.", eventName)
	} else {
		messageText = fmt.Sprintf("%s is not a valid option.", selectedOption)
	}

	return sendSimpleMessage(chatId, messageText)
}

func SendEventPropertySetAcknowledge(chatId int64, eventName, propertyName, propertyValue string) (err error) {

	messageText := fmt.Sprintf("Ok! Event \"%s\" %s was set to \"%s\".", eventName, propertyName, propertyValue)

	return sendSimpleMessage(chatId, messageText)
}

func SendSelectEventMessage(chatId, replyToMessageId, userFromId int64, events []types.Event) (err error) {

	var text bytes.Buffer

	text.WriteString(fmt.Sprint("Please choose one of the following:"))
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
	smr.AddText(fmt.Sprintf("Your current event is \"%s\".", GetEventDescription(currentEvent)))
	_, err = smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	return nil
}

func SendAllEventsMessage(chatId int64, events []types.Event) error {

	var messageText bytes.Buffer

	if len(events) > 0 {
		currentEvent := GetCurrentEvent(chatId)

		messageText.WriteString(fmt.Sprint("Events available in this group are:"))
		for _, event := range events {
			if (currentEvent.Name == event.Name) {
				messageText.WriteString(fmt.Sprintf("\n%s [C]", GetEventDescription(event)))
			} else {
				messageText.WriteString(fmt.Sprintf("\n%s", GetEventDescription(event)))
			}
		}
	} else {
		messageText.WriteString(fmt.Sprint("There are no events in this group."))
	}

	return sendSimpleMessage(chatId, messageText.String())
}

func SendAttendanceList(chatId, replyToMessageId int64, eventDescription string, participants []types.User) (err error) {

	var text bytes.Buffer
	participantWord := "participants"
	if len(participants) == 1 {
		participantWord = "participant"
	}

	text.WriteString(fmt.Sprintf("%d %s for %s:", len(participants), participantWord, eventDescription))
	for _, participant := range participants {
		text.WriteString(fmt.Sprintf("\n%s %s", participant.FirstName, participant.LastName))
	}

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(chatId)
	smr.AddText(text.String())
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

func SendNoParticipantsYetMessage(chatId, replyToMessageId int64, eventDescription string) (err error) {

	text := fmt.Sprintf("There are no participants confirmed for %s yet.", eventDescription)

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

func SendAttendanceConfirmationMessage(chatId, replyToMessageId int64, userName, eventDescription string) (err error) {

	text := fmt.Sprintf("Thanks %s! You have been confirmed to attend %s", userName, eventDescription)

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

func SendAttendanceRemovalConfirmationMessage(chatId, replyToMessageId int64, userName, eventDescription string) (err error) {

	text := fmt.Sprintf("Thanks %s! You have been removed from %s participants", userName, eventDescription)

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

func SendTooManyParametersMessage(chatId int64, commandName string) (err error) {

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(chatId)
	smr.AddText(fmt.Sprintf("Too many parameters. Please use: /%s", commandName))
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

func sendSimpleMessage(chatId int64, messageText string) (err error) {

	smr := newMessageRequest(chatId, messageText)
	_, err = smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}

	return err
}

func newMessageRequest(chatId int64, messageText string) *requests.SendMessageRequest {

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(chatId)
	smr.AddText(messageText)

	return smr
}

func newMessageReplyRequest(chatId, messageFromId int64, messageText string) *requests.SendMessageRequest {

	smr := newMessageRequest(chatId, messageText)
	smr.AddReplyToMessageId(messageFromId)
	smr.AddForceReply( types.ForceReply { ForceReply: true, Selective: true } )

	return smr
}