package global

const (
	BOT_NAME = "EventTrackingBot"
	TOKEN = "133198388:AAHHbnm7cNHMEF6hmdehKTCMRDrFMN46n-U"
	BASE_URL = "https://api.telegram.org/bot" + TOKEN + "/"

	MESSAGE_TYPE_CREATE_EVENT_PROVIDE_NAME = "create_event_provide_name"
	MESSAGE_TYPE_DELETE_EVENT_PROVIDE_INDEX = "delete_event_provide_index"
	MESSAGE_TYPE_SELECT_CURRENT_EVENT = "select_current_event"

	GET_UPDATES_DEFAULT_TIMEOUT = 0
	GET_UPDATES_DEFAULT_LIMIT = 100
)