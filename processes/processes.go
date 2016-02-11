package processes

import (
	"fmt"
	"github.com/kinslayere/eventtrackingbot/types"
	"github.com/kinslayere/eventtrackingbot/requests"
	"log"
)

func processStartOrHelp(update types.Update) {

	log.Printf("Processing %s", update.Message.Text)

	text := fmt.Sprintf("Hello %s!", update.Message.From.FirstName)
	text += "\nThis is an event attendance tracking bot.\nAvailable commands are:"
	text += "\n/start - Start this bot"
	text += "\n/help - Info and commands"
	text += "\n/create_event - Create a new event"
	text += "\n/delete_event - Delete an existing event"
	text += "\n/current_event - Get current event being tracked"
	text += "\n/select_event - Change the current event"
	text += "\n/all_events - Get all events created in this group"

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(update.Message.Chat.Id)
	smr.AddText(text)
	if _, err := smr.Execute(); err != nil {
		log.Printf("Error sending response to %s: %v", update.Message.Text, err)
	}
}

func processPrivateChat(update types.Update) {

	text := fmt.Sprintf("Hello %s!\n", update.Message.From.FirstName)
	text += "Sorry, but this bot is intended for usage in group chats only."

	smr := requests.NewSendMessageRequest()
	smr.AddChatId(update.Message.Chat.Id)
	smr.AddText(text)
	if _, err := smr.Execute(); err != nil {
		log.Printf("Error sending response to %s: %v", update.Message.Text, err)
	}
}