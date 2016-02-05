package processes

import (
	"fmt"
	"github.com/kinslayere/eventtrackingbot/types"
	"github.com/kinslayere/eventtrackingbot/requests"
)

func processStartOrHelp(update types.Update) {

	text := fmt.Sprintf("Hello %s!\n", update.Message.From.FirstName)
	text += "This is an event tracking attendance bot.\nAvailable commands are:"

	smr := requests.NewSendMessageRequest()
	smr.AddChatIdInt(update.Message.Chat.Id)
	smr.AddText(text)
	smr.DoRequest()
}

func ProcessPrivateChat(update types.Update) {

	text := fmt.Sprintf("Hello %s!\n", update.Message.From.FirstName)
	text += "Sorry, but this bot is intended for usage in group chats only."

	smr := requests.NewSendMessageRequest()
	smr.AddChatIdInt(update.Message.Chat.Id)
	smr.AddText(text)
	smr.DoRequest()
}