package processes

import (
	"log"
	"fmt"
	"github.com/kinslayere/eventtrackingbot/requests"
	"github.com/kinslayere/eventtrackingbot/types"
)

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
