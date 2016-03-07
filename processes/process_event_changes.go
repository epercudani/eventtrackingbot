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
	"errors"
)

func processCreateEventWithoutName(update types.Update) (err error) {

	fields := strings.Fields(strings.TrimSpace(update.Message.Text))
	smr := requests.NewSendMessageRequest()

	switch len(fields) {
	case 1:

		// Send acknowledge message and wait for event name
		smr.AddChatId(update.Message.Chat.Id)
		smr.AddText(fmt.Sprintf("Hi %s! I'll help you to create your event. Please tell me the event's name.", update.Message.From.FirstName))
		smr.AddReplyToMessageId(update.Message.MessageId)
		smr.AddForceReply( types.ForceReply { ForceReply: true, Selective: true } )
		response, err := smr.Execute()
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return err
		}

		messageSent := response.Result
		log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
		err = services.SetPendingResponseToMessage(update.Message.From.Id, messageSent.MessageId, global.MESSAGE_TYPE_EVENT_PROVIDE_NAME)
		if err != nil {
			return err
		}

	default:
		smr.AddChatId(update.Message.Chat.Id)
		smr.AddText(fmt.Sprintf("Too many parameters. Please use: /create_event"))
		_, err = smr.Execute()
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return err
		}
	}

	return nil
}

func processSetWhen(update types.Update) (err error) {

	fields := strings.Fields(strings.TrimSpace(update.Message.Text))

	switch len(fields) {
	case 1:

		currentEvent := services.GetCurrentEvent(update.Message.Chat.Id)
		if currentEvent.Name != "" {
			// Send acknowledge message and wait for event name
			services.SendRequestPropertyMessage(update.Message.Chat.Id, update.Message.MessageId, update.Message.From.Id, currentEvent, global.EVENT_PROPERTY_DATE)
		} else {
			return processCurrentEventNotSet(update)
		}

	default:
		services.SendTooManyParametersMessage(update.Message.Chat.Id)
	}

	return nil
}

func processSetWhere(update types.Update) (err error) {

	fields := strings.Fields(strings.TrimSpace(update.Message.Text))

	switch len(fields) {
	case 1:

		currentEvent := services.GetCurrentEvent(update.Message.Chat.Id)
		if currentEvent.Name != "" {
			// Send acknowledge message and wait for event name
			services.SendRequestPropertyMessage(update.Message.Chat.Id, update.Message.MessageId, update.Message.From.Id, currentEvent, global.EVENT_PROPERTY_PLACE)
		} else {
			return processCurrentEventNotSet(update)
		}

	default:
		smr := requests.NewSendMessageRequest()
		smr.AddChatId(update.Message.Chat.Id)
		smr.AddText(fmt.Sprintf("Too many parameters. Please use: /set_when"))
		_, err = smr.Execute()
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return err
		}
	}

	return nil
}

func processDeleteEventWithoutName(update types.Update) (err error) {

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
			text.WriteString(fmt.Sprintf("There are no events in this group."))
		}

		smr.AddChatId(update.Message.Chat.Id)
		smr.AddText(text.String())
		smr.AddReplyToMessageId(update.Message.MessageId)
		smr.AddForceReply( types.ForceReply { ForceReply: true, Selective: true } )
		response, err := smr.Execute()
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return err
		}

		messageSent := response.Result
		log.Printf("Message sent id: %d in response to: %s", messageSent.MessageId, messageSent.ReplyToMessage.Text)
		err = services.SetPendingResponseToMessage(update.Message.From.Id, messageSent.MessageId, global.MESSAGE_TYPE_DELETE_EVENT_PROVIDE_INDEX)
		if err != nil {
			return err
		}

	default:
		smr.AddChatId(update.Message.Chat.Id)
		smr.AddText(fmt.Sprintf("Too many parameters. Please use: /delete_event"))
		_, err = smr.Execute()
		if err != nil {
			log.Printf("Error sending message: %v", err)
			return err
		}
	}

	return nil
}

func processSetEventName(update types.Update) (err error) {

	log.Printf("Processing event name")

	eventName := update.Message.Text
	if len(strings.TrimSpace(eventName)) > 0 {
		services.CreateEvent(update.Message.Chat.Id, eventName)
	}

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(update.Message.Chat.Id)
	smr.AddText(fmt.Sprintf("Congrats %s! Event \"%s\" was created and set as the current event in this group.", update.Message.From.FirstName, eventName))
	_, err = smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}

	return
}

func processSetEventProperty(update types.Update, property string) (err error) {

	log.Printf(fmt.Sprintf("Processing event %s with value \"%s\"", property, update.Message.Text))
	if len(strings.TrimSpace(update.Message.Text)) == 0 {
		errors.New(fmt.Sprintf("Error setting property %s with empty value", property))
	}

	currentEvent := services.GetCurrentEvent(update.Message.Chat.Id)

	switch property {
	case global.EVENT_PROPERTY_DATE:
		currentEvent.Date = update.Message.Text
	case global.EVENT_PROPERTY_PLACE:
		currentEvent.Place = update.Message.Text
	}

	err = services.SaveEvent(update.Message.Chat.Id, currentEvent)
	if err != nil {
		log.Printf("Error saving message property %s: %v", property, err)
		return
	}

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(update.Message.Chat.Id)
	smr.AddText(fmt.Sprintf("Ok! Event \"%s\" %s was set to \"%s\".", currentEvent.Name, property, currentEvent.Date))
	_, err = smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}

	return
}

func processIndexToDeleteEvent(update types.Update) (err error) {

	log.Printf("Processing delete index")

	var eventName string

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
			eventName = services.GetGroupEventNames(update.Message.Chat.Id)[index - 1]
			services.DeleteEvent(update.Message.Chat.Id, eventName)
			deleted = true
		}
	}

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(update.Message.Chat.Id)

	if deleted {
		smr.AddText(fmt.Sprintf("%s, you have successfully deleted event \"%s\".", update.Message.From.FirstName, eventName))
	} else {
		smr.AddText(fmt.Sprintf("%s is not a valid option.", update.Message.Text))
	}

	_, err = smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
	}

	return
}

func processIndexToSelectCurrentEvent(update types.Update) (err error) {

	log.Printf("Processing select index")

	var eventName string

	selected := false

	regexpIndex := regexp.MustCompile("^/[0-9]+$")

	if regexpIndex.MatchString(update.Message.Text) {

		index, err := strconv.Atoi(update.Message.Text[1:])
		if err != nil {
			log.Fatalf("Error: %v", err)
			return err
		}

		size, err := persistence.GetSortedSetSize(fmt.Sprintf(persistence.KEY_GROUP_EVENTS, update.Message.Chat.Id))
		if err != nil {
			log.Printf("Error getting group events count: %v", err)
			return err
		}

		if index >= 1 && index <= size {
			groupId := update.Message.Chat.Id
			eventName = services.GetGroupEventNames(groupId)[index - 1]
			services.SetCurrentEvent(update.Message.Chat.Id, eventName)
			selected = true
		}
	}

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(update.Message.Chat.Id)

	if selected {
		smr.AddText(fmt.Sprintf("%s is now set as the current event in this group.", eventName))
	} else {
		smr.AddText(fmt.Sprintf("%s is not a valid option.", update.Message.Text))
	}

	_, err = smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	return nil
}

func processAllEvents(update types.Update) (err error) {

	events := services.GetGroupEvents(update.Message.Chat.Id)

	var text bytes.Buffer
	if len(events) > 0 {
		currentEvent := services.GetCurrentEvent(update.Message.Chat.Id)

		text.WriteString(fmt.Sprintf("Events available in this group are:"))
		for _, event := range events {
			if (currentEvent.Name == event.Name) {
				text.WriteString(fmt.Sprintf("\n%s [C]", services.GetEventDescription(event)))
			} else {
				text.WriteString(fmt.Sprintf("\n%s", services.GetEventDescription(event)))
			}
		}
	} else {
		text.WriteString(fmt.Sprintf("There are no events in this group."))
	}

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(update.Message.Chat.Id)
	smr.AddText(text.String())
	_, err = smr.Execute()
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	return nil
}

func processCurrentEvent(update types.Update) (err error) {

	currentEvent := services.GetCurrentEvent(update.Message.Chat.Id)
	if currentEvent.Name != "" {
		services.SendCurrentEventMessage(update.Message.Chat.Id, currentEvent)
	} else {
		return processCurrentEventNotSet(update)
	}

	return nil
}

func processSelectEvent(update types.Update) (err error) {

	events := services.GetGroupEvents(update.Message.Chat.Id)
	if len(events) > 0 {
		return services.SendSelectEventMessage(update.Message.Chat.Id, update.Message.MessageId, update.Message.From.Id, events)
	} else {
		return services.SendNoEventsInGroupMessage(update.Message.Chat.Id, update.Message.MessageId)
	}
}

func processCurrentEventNotSet(update types.Update) (err error) {

	err = services.SendCurrentEventNotSetMessage(update.Message.Chat.Id, update.Message.MessageId)
	if err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	processSelectEvent(update)

	return nil
}