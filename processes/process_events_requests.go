package processes

import (
	"strings"
	"fmt"
	"log"
	"github.com/kinslayere/eventtrackingbot/global"
	"github.com/kinslayere/eventtrackingbot/persistence"
	"github.com/kinslayere/eventtrackingbot/types"
	"github.com/kinslayere/eventtrackingbot/services"
	"regexp"
	"strconv"
)

func processCreateEvent(update types.Update) error {

	fields := strings.Fields(strings.TrimSpace(update.Message.Text))

	switch len(fields) {
	case 1:
		return services.SendRequestEventNameMessage(
			update.Message.Chat.Id, update.Message.MessageId,
			update.Message.From.Id, update.Message.From.FirstName)

	default:
		eventName := strings.Join(fields[1:], " ")
		err := services.CreateEvent(update.Message.Chat.Id, eventName)
		if err != nil {
			log.Printf("Error processing create event: %v", err)
			return err
		}

		err = services.SendCreatedEventAcknowledge(update.Message.Chat.Id, update.Message.From.FirstName, eventName)
		if err != nil {
			log.Printf("Error processing create event: %v", err)
			return err
		}
	}

	return nil
}

func processSetWhen(update types.Update) error {

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
		services.SendTooManyParametersMessage(update.Message.Chat.Id, global.COMMAND_NAME_SET_WHEN)
	}

	return nil
}

func processSetWhere(update types.Update) error {

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
		services.SendTooManyParametersMessage(update.Message.Chat.Id, global.COMMAND_NAME_SET_WHERE)
	}

	return nil
}

func processDeleteEventWithoutName(update types.Update) error {

	fields := strings.Fields(strings.TrimSpace(update.Message.Text))

	switch len(fields) {
	case 1:
		events := services.GetGroupEvents(update.Message.Chat.Id)
		if len(events) > 0 {
			return services.SendResponseToDeleteEventMessage(
				update.Message.Chat.Id, update.Message.MessageId,
				update.Message.From.Id, update.Message.From.FirstName, events)
		} else {
			return services.SendNoEventsInGroupMessage(update.Message.Chat.Id, update.Message.MessageId)
		}

	default:
		return services.SendTooManyParametersMessage(update.Message.Chat.Id, global.COMMAND_NAME_DELETE_EVENT)
	}

	return nil
}

func processSetEventName(update types.Update) error {

	eventName := update.Message.Text
	if len(strings.TrimSpace(eventName)) > 0 {
		services.CreateEvent(update.Message.Chat.Id, eventName)
	}

	return services.SendCreatedEventAcknowledge(update.Message.Chat.Id, update.Message.From.FirstName, eventName)
}

func processSetEventProperty(update types.Update, property string) (err error) {

	if len(strings.TrimSpace(update.Message.Text)) == 0 {
		return fmt.Errorf("Error setting property %s with empty value", property)
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

	return services.SendEventPropertySetAcknowledge(update.Message.Chat.Id, currentEvent.Name, property, update.Message.Text)
}

func processIndexToDeleteEvent(update types.Update) error {

	log.Print("Processing delete index")

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

	return services.SendDeletedEventAcknowledge(update.Message.Chat.Id, update.Message.From.FirstName, eventName, update.Message.Text, deleted)
}

func processIndexToSelectCurrentEvent(update types.Update) (err error) {

	log.Print("Processing select index")

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

	return services.SendSelectedCurrentEventAcknowledge(update.Message.Chat.Id, eventName, update.Message.Text, selected)
}

func processAllEvents(update types.Update) (err error) {

	events := services.GetGroupEvents(update.Message.Chat.Id)

	return services.SendAllEventsMessage(update.Message.Chat.Id, events)
}

func processParticipants(update types.Update) error {

	currentEvent := services.GetCurrentEvent(update.Message.Chat.Id)
	if currentEvent.Name != "" {
		participants := services.GetParticipantsToEvent(update.Message.Chat.Id, currentEvent.Name)
		if len(participants) > 0 {
			return services.SendAttendanceList(update.Message.Chat.Id, update.Message.MessageId, services.GetEventDescription(currentEvent), participants)
		} else {
			return services.SendNoParticipantsYetMessage(update.Message.Chat.Id, update.Message.MessageId, services.GetEventDescription(currentEvent))
		}
	} else {
		return processCurrentEventNotSet(update)
	}
}

func processWillGo(update types.Update) error {

	currentEvent := services.GetCurrentEvent(update.Message.Chat.Id)
	if currentEvent.Name != "" {
		err := services.AddParticipantToEvent(update.Message.Chat.Id, currentEvent.Name, update.Message.From)
		if err != nil {
			return err
		}

		return services.SendAttendanceConfirmationMessage(update.Message.Chat.Id, update.Message.MessageId,
			update.Message.From.FirstName, services.GetEventDescription(currentEvent))
	} else {
		return processCurrentEventNotSet(update)
	}
}

func processWontGo(update types.Update) error {

	currentEvent := services.GetCurrentEvent(update.Message.Chat.Id)
	if currentEvent.Name != "" {
		err := services.RemoveParticipantFromEvent(update.Message.Chat.Id, currentEvent.Name, update.Message.From)
		if err != nil {
			return err
		}

		return services.SendAttendanceRemovalConfirmationMessage(update.Message.Chat.Id, update.Message.MessageId,
			update.Message.From.FirstName, services.GetEventDescription(currentEvent))
	} else {
		return processCurrentEventNotSet(update)
	}
}

func processCurrentEvent(update types.Update) error {

	currentEvent := services.GetCurrentEvent(update.Message.Chat.Id)
	if currentEvent.Name != "" {
		services.SendCurrentEventMessage(update.Message.Chat.Id, currentEvent)
	} else {
		return processCurrentEventNotSet(update)
	}

	return nil
}

func processSelectEvent(update types.Update) error {

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