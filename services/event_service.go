package services
import (
	"fmt"
	"log"
	"github.com/kinslayere/eventtrackingbot/persistence"
)

func CreateEvent(groupId int64, eventName string) {

	eventKey := fmt.Sprintf(persistence.KEY_EVENT, groupId, eventName)
	err := persistence.AddStringFieldToHash(eventKey, persistence.KEY_EVENT_NAME, eventName)
	if err != nil {
		log.Printf("Error creating event: %v", err)
		return
	}

	SetCurrentEvent(groupId, eventKey)
	AddEventToGroup(groupId, eventName)
}

func DeleteEvent(groupId int64, eventName string) {

	eventKey := fmt.Sprintf(persistence.KEY_EVENT, groupId, eventName)
	_, err := persistence.Delete(eventKey)
	if err != nil {
		log.Printf("Error deleting event: %v", err)
		return
	}

	UnsetCurrentEvent(groupId, eventKey)
	RemoveEventFromGroup(groupId, eventName)
}

func GetCurrentEventName(groupId int64) string {

	eventKey, err := persistence.GetString(fmt.Sprintf(persistence.KEY_CURRENT_EVENT, groupId))
	if err != nil {
		log.Printf("Error getting current event key: %v", err)
		return ""
	}

	eventName, err := persistence.GetStringFieldFromHash(eventKey, persistence.KEY_EVENT_NAME)
	if err != nil {
		log.Printf("Error getting current event name: %v", err)
		return ""
	}

	return eventName
}

func SetCurrentEvent(groupId int64, eventKey string) {

	err := persistence.SaveString(fmt.Sprintf(persistence.KEY_CURRENT_EVENT, groupId), eventKey)
	if err != nil {
		log.Printf("Error setting current event: %v", err)
		return
	}
}

func UnsetCurrentEvent(groupId int64, eventKey string) {

	_, err := persistence.Delete(fmt.Sprintf(persistence.KEY_CURRENT_EVENT, groupId))
	if err != nil {
		log.Printf("Error unsetting current event: %v", err)
		return
	}
}

func AddEventToGroup(groupId int64, eventName string) {

	key := fmt.Sprintf(persistence.KEY_GROUP_EVENTS, groupId)
	index, err := persistence.GetSortedSetSize(key)
	if err != nil {
		return
	}
	err = persistence.AddStringToSortedSet(key, index, eventName)
	if err != nil {
		log.Printf("Error adding event to group: %v", err)
		return
	}
}

func RemoveEventFromGroup(groupId int64, eventName string) {

	err := persistence.RemoveFromSortedSet(fmt.Sprintf(persistence.KEY_GROUP_EVENTS, groupId), eventName)
	if err != nil {
		log.Printf("Error removing event from group: %v", err)
		return
	}
}

func GetGroupEventNames(groupId int64) (events []string) {

	events, err := persistence.GetStringsFromSortedSet(fmt.Sprintf(persistence.KEY_GROUP_EVENTS, groupId))
	if err != nil {
		log.Printf("Error getting events for group %d: %v", groupId, err)
	}

	return
}