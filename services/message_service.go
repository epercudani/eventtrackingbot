package services

import (
	"fmt"
	"log"
	"bytes"
	"github.com/kinslayere/eventtrackingbot/persistence"
	"github.com/kinslayere/eventtrackingbot/clients/telegram"
	"github.com/kinslayere/eventtrackingbot/types/telegram"
	"github.com/kinslayere/eventtrackingbot/global"
)

func SendRequestEventNameMessage(chatId, messageId, userId int64, userName string) (err error) {

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
	err = setPendingResponseToMessage(userId, messageSent.MessageId, global.MESSAGE_TYPE_EVENT_PROVIDE_NAME)
	if err != nil {
		return err
	}

	return nil
}

func SendResponseToDeleteEventMessage(chatId, messageId, userId int64, userName string, events []types.Event) (err error) {

	var messageText bytes.Buffer

	messageText.WriteString(fmt.Sprintf("Hi %s! Which event do you want to delete?", userName))
	messageText.WriteString(getSelectEventOptionsText(chatId, events))

	smr := newMessageReplyRequest(chatId, messageId, messageText.String())
	response, err := smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	messageSent := response.Result
	log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
	err = setPendingResponseToMessage(userId, messageSent.MessageId, global.MESSAGE_TYPE_DELETE_EVENT_PROVIDE_INDEX)
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

func SendSelectEventMessage(chatId, messageId, userFromId int64, events []types.Event, noCurrent bool) (err error) {

	var messageText bytes.Buffer

	if noCurrent == true {
		messageText.WriteString("There is no current event in this group. ")
	}
	messageText.WriteString(fmt.Sprint("Please choose one of the following:"))
	messageText.WriteString(getSelectEventOptionsText(chatId, events))

	smr := newMessageReplyRequest(chatId, messageId, messageText.String())
	log.Printf("%v - %v - %v", chatId, messageId, messageText.String())
	response, err := smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	messageSent := response.Result
	log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
	err = setPendingResponseToMessage(userFromId, messageSent.MessageId, global.MESSAGE_TYPE_SELECT_CURRENT_EVENT)
	if err != nil {
		return err
	}

	return nil
}

func SendRequestPropertyMessage(chatId, messageId, userFromId int64, event types.Event, property string) (err error) {

	messageText := fmt.Sprintf("Please, provide \"%s\" %s.", event.Name, property)

	smr := newMessageReplyRequest(chatId, messageId, messageText)
	response, err := smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	messageSent := response.Result
	log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
	err = setPendingResponseToMessage(userFromId, messageSent.MessageId, global.PROPERTY_TO_MESSAGE_TYPE_MAP[property])
	if err != nil {
		return err
	}

	return nil
}

func SendCurrentEventPropertyMessage(chatId int64, currentEvent types.Event, property string) error {

	propertyValue, err := GetEventProperty(currentEvent, property)
	if err != nil {
		log.Printf("Error sending current event property message: %v", err)
		return err
	}

	var messageText string
	if propertyValue != "" {
		messageText = fmt.Sprintf("Current event %s is <b>%s</b>.", property, propertyValue)
	} else {
		messageText = fmt.Sprintf("Current event %s has not been set.", property)
	}

	return sendSimpleMessage(chatId, messageText)
}

func SendCurrentEventMessage(chatId int64, currentEvent types.Event) error {

	messageText := fmt.Sprintf("Your current event is %s.", GetEventDescription(currentEvent))

	return sendSimpleMessage(chatId, messageText)
}

func SendAllEventsMessage(chatId int64, events []types.Event) error {

	var messageText bytes.Buffer

	if len(events) > 0 {
		currentEvent := GetCurrentEvent(chatId)

		messageText.WriteString(fmt.Sprint("Events available in this group are:"))
		for _, event := range events {
			if (currentEvent.Name == event.Name) {
				messageText.WriteString(fmt.Sprintf("\n[<b>%s</b>]", event.Name))
			} else {
				messageText.WriteString(fmt.Sprintf("\n%s", event.Name))
			}
		}
	} else {
		messageText.WriteString(fmt.Sprint("There are no events in this group."))
	}

	return sendSimpleMessage(chatId, messageText.String())
}

func SendAttendanceList(chatId, messageId int64, eventDescription string, participants []types.User) (err error) {

	var messageText bytes.Buffer
	participantWord := "participants"
	if len(participants) == 1 {
		participantWord = "participant"
	}

	messageText.WriteString(fmt.Sprintf("%d %s for %s:", len(participants), participantWord, eventDescription))
	for _, participant := range participants {
		messageText.WriteString(fmt.Sprintf("\n%s %s", participant.FirstName, participant.LastName))
	}

	return sendSimpleMessage(chatId, messageText.String())
}

func SendNoParticipantsYetMessage(chatId, messageId int64, eventDescription string) (err error) {

	messageText := fmt.Sprintf("There are no participants confirmed for %s yet.", eventDescription)

	return sendSimpleMessage(chatId, messageText)
}

func SendAttendanceConfirmationMessage(chatId int64, userName string) error {

	messageText := fmt.Sprintf("Thanks %s! You have been confirmed.", userName)

	return sendSimpleMessage(chatId, messageText)
}

func SendAttendanceRemovalConfirmationMessage(chatId int64, userName string) error {

	messageText := fmt.Sprintf("Thanks %s! You have been removed.", userName)

	return sendSimpleMessage(chatId, messageText)
}

func SendNoEventsInGroupMessage(chatId, messageId int64) (err error) {

	messageText := "There are no events in this group."

	smr := newMessageReplyRequest(chatId, messageId, messageText)
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

	messageText := fmt.Sprintf("Too many parameters. Please use: /%s", commandName)

	return sendSimpleMessage(chatId, messageText)
}

func setPendingResponseToMessage(userId, messageId int64, messageType string) (err error) {

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

func getSelectEventOptionsText(chatId int64, events []types.Event) string {

	var text bytes.Buffer
	currentEvent := GetCurrentEvent(chatId)

	for i, event := range events {

		var eventName string

		if (currentEvent.Name == event.Name) {
			eventName = fmt.Sprintf("\n/%d <b>[</b>%s<b>]</b>", i + 1, GetEventDescription(event))
		} else {
			eventName = fmt.Sprintf("\n/%d %s", i + 1, GetEventDescription(event))
		}
		text.WriteString(eventName)
	}

	return text.String()
}

func sendSimpleMessage(chatId int64, messageText string) (err error) {

	smr := newMessageRequest(chatId, messageText)
	_, err = smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}

	return err
}

func newMessageRequest(chatId int64, messageText string) *telegram.SendMessageRequest {

	smr := telegram.NewSendMessageRequest()
	smr.AddChatId(chatId)
	smr.AddText(messageText)
	smr.AddParseMode(global.MESSAGE_PARSE_MODE_HTML)

	return smr
}

func newMessageReplyRequest(chatId, messageFromId int64, messageText string) *telegram.SendMessageRequest {

	smr := newMessageRequest(chatId, messageText)
	smr.AddReplyToMessageId(messageFromId)
	smr.AddForceReply( types.ForceReply { ForceReply: true, Selective: true } )

	return smr
}