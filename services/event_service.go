package services
import (
	"fmt"
	"log"
	"bytes"
	"encoding/json"
	"github.com/kinslayere/eventtrackingbot/persistence"
	"github.com/kinslayere/eventtrackingbot/types/telegram"
	"github.com/kinslayere/eventtrackingbot/global"
	"errors"
)

func GetEventKeyFromGroupAndName(chatId int64, eventName string) string {
	return fmt.Sprintf(persistence.KEY_EVENT, chatId, eventName)
}

func GetEvent(eventKey string) (event types.Event) {

	eventJson, err := persistence.GetString(eventKey)
	if err != nil {
		log.Printf("Error getting current event: %v", err)
	}

	if eventJson == "" {
		return
	}

	err = json.Unmarshal([]byte(eventJson), &event)
	if err != nil {
		log.Printf("Error unmarshalling current event: %v", err)
	}

	return
}

func CreateEvent(chatId int64, eventName string) error {

	event := types.Event{ Name: eventName }

	err := SaveEvent(chatId, event)
	if err != nil {
		log.Printf("Error creating event %s: %v", eventName, err)
		return err
	}

	err = SetCurrentEvent(chatId, eventName)
	if err != nil {
		log.Printf("Error setting event %s as current: %v", eventName, err)
		return err
	}

	AddEventToGroup(chatId, eventName)
	if err != nil {
		log.Printf("Error adding event %s to group %d: %v", eventName, chatId, err)
		return err
	}

	return nil
}

func SaveEvent(chatId int64, event types.Event) error {

	eventJson, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return err
	}

	eventKey := GetEventKeyFromGroupAndName(chatId, event.Name)
	err = persistence.SaveString(eventKey, string(eventJson))
	if err != nil {
		log.Printf("Error saving event: %v", err)
		return err
	}

	return nil
}

func DeleteEvent(chatId int64, eventName string) {

	currentEvent := GetCurrentEvent(chatId)
	eventKey := GetEventKeyFromGroupAndName(chatId, eventName)
	_, err := persistence.Delete(eventKey)
	if err != nil {
		log.Printf("Error deleting event: %v", err)
		return
	}

	if currentEvent.Name == eventName {
		UnsetCurrentEvent(chatId)
	}

	RemoveEventFromGroup(chatId, eventName)
}

func GetEventProperty(event types.Event, property string) (string, error) {

	switch property {
	case global.EVENT_PROPERTY_NAME:
		return event.Name, nil
	case global.EVENT_PROPERTY_DATE:
		return event.Date, nil
	case global.EVENT_PROPERTY_PLACE:
		return event.Place, nil
	default:
		return "", errors.New(fmt.Sprintf("GetEventProperty. Invalid property \"%s\"", property))
	}
}

func SetEventProperty(chatId int64, currentEvent types.Event, property, value string) error {

	switch property {
	case global.EVENT_PROPERTY_DATE:
		currentEvent.Date = value
	case global.EVENT_PROPERTY_PLACE:


		currentEvent.Place = value
	}

	err := SaveEvent(chatId, currentEvent)
	if err != nil {
		log.Printf("Error saving message property %s: %v", property, err)
		return err
	}

	return nil
}

func GetCurrentEvent(chatId int64) types.Event {

	eventKey, err := persistence.GetString(fmt.Sprintf(persistence.KEY_CURRENT_EVENT, chatId))
	if err != nil {
		log.Printf("Error getting current event key: %v", err)
	}

	return GetEvent(eventKey)
}

func GetEventDescription(event types.Event) string {

	var description bytes.Buffer

	description.WriteString(fmt.Sprintf("<b>%s</b>", event.Name))

	if event.Date != "" {
		description.WriteString(fmt.Sprintf(" on <b>%s</b>", event.Date))
	}

	if event.Place != "" {
		description.WriteString(fmt.Sprintf(" at <b>%s</b>", event.Place))
	}

	return description.String()
}

func SetCurrentEvent(chatId int64, eventName string) error {

	eventKey := fmt.Sprintf(persistence.KEY_EVENT, chatId, eventName)
	err := persistence.SaveString(fmt.Sprintf(persistence.KEY_CURRENT_EVENT, chatId), eventKey)
	if err != nil {
		log.Printf("Error setting current event: %v", err)
		return err
	}

	return nil
}

func UnsetCurrentEvent(chatId int64) {

	_, err := persistence.Delete(fmt.Sprintf(persistence.KEY_CURRENT_EVENT, chatId))
	if err != nil {
		log.Printf("Error unsetting current event: %v", err)
		return
	}
}

func AddEventToGroup(chatId int64, eventName string) error {

	key := fmt.Sprintf(persistence.KEY_GROUP_EVENTS, chatId)
	index, err := persistence.GetSortedSetSize(key)
	if err != nil {
		log.Printf("Error adding event to group: %v", err)
		return err
	}
	err = persistence.AddStringToSortedSet(key, index, eventName)
	if err != nil {
		log.Printf("Error adding event to group: %v", err)
		return err
	}

	return nil
}

func RemoveEventFromGroup(chatId int64, eventName string) {

	err := persistence.RemoveFromSortedSet(fmt.Sprintf(persistence.KEY_GROUP_EVENTS, chatId), eventName)
	if err != nil {
		log.Printf("Error removing event from group: %v", err)
		return
	}

	_, err = persistence.Delete(fmt.Sprintf(persistence.KEY_EVENT_PARTICIPANTS, chatId, eventName))
	if err != nil {
		log.Printf("Error removing participants data when deleting event: %v", err)
		return
	}
}

func GetGroupEvents(chatId int64) (events []types.Event) {

	eventNames, err := persistence.GetStringsFromSortedSet(fmt.Sprintf(persistence.KEY_GROUP_EVENTS, chatId))
	if err != nil {
		log.Printf("Error getting event names for group %d: %v", chatId, err)
	}

	for _, eventName := range eventNames {
		eventKey := GetEventKeyFromGroupAndName(chatId, eventName)
		events = append(events, GetEvent(eventKey))
	}

	return
}

func GetGroupEventNames(chatId int64) (events []string) {

	events, err := persistence.GetStringsFromSortedSet(fmt.Sprintf(persistence.KEY_GROUP_EVENTS, chatId))
	if err != nil {
		log.Printf("Error getting events for group %d: %v", chatId, err)
	}

	return
}

func GetParticipantsToEvent(chatId int64, eventName string) (participants []types.User) {

	key := fmt.Sprintf(persistence.KEY_EVENT_PARTICIPANTS, chatId, eventName)
	participantsJson, err := persistence.GetStringsFromSet(key)
	if err != nil {
		log.Printf("Error getting participants for group %d: %v", chatId, err)
	}

	var user types.User
	for _, participantJson := range participantsJson {
		err = json.Unmarshal([]byte(participantJson), &user)
		if err != nil {
			log.Printf("Error unmarshalling current event: %v", err)
		}

		participants = append(participants, user)
	}

	return
}

func AddParticipantToEvent(chatId int64, eventName string, participant types.User) (err error) {

	key := fmt.Sprintf(persistence.KEY_EVENT_PARTICIPANTS, chatId, eventName)

	participantJson, err := json.Marshal(participant)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return err
	}

	return persistence.AddStringToSet(key, string(participantJson))
}

func RemoveParticipantFromEvent(chatId int64, eventName string, participant types.User) (err error) {

	key := fmt.Sprintf(persistence.KEY_EVENT_PARTICIPANTS, chatId, eventName)

	participantJson, err := json.Marshal(participant)
	if err != nil {
		log.Printf("Error marshaling event: %v", err)
		return err
	}

	return persistence.RemoveFromSet(key, string(participantJson))
}