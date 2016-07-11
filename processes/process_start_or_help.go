package processes

import (
	"log"
	"fmt"
	"github.com/kinslayere/eventtrackingbot/requests"
	"github.com/kinslayere/eventtrackingbot/types"
)

func processStartOrHelp(update types.Update) {

	log.Printf("Processing %s", update.Message.Text)

	text := fmt.Sprintf("Hello %s!", update.Message.From.FirstName)
	text += "\nThis is an event attendance tracking bot.\nAvailable commands are:"
	text += "\n/start - Start this bot"
	text += "\n/help - Info and commands"
	text += "\n/create_event - Create a new event"
	text += "\n/set_when - Set current event's date"
	text += "\n/set_where - Set current event's place"
	text += "\n/delete_event - Delete an existing event"
	text += "\n/current_event - Get current event being tracked"
	text += "\n/select_event - Change the current event"
	text += "\n/all_events - Get all events created in this group"
	text += "\n/participants - Get confirmed participants to current event"
	text += "\n/will_go - Confirm your attendance to current event"
	text += "\n/wont_go - Deny or remove your attendance to current event"

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(update.Message.Chat.Id)
	smr.AddText(text)
	if _, err := smr.Execute(); err != nil {
		log.Printf("Error sending response to %s: %v", update.Message.Text, err)
	}
}
