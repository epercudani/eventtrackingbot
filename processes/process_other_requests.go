package processes

import (
    "log"
    "fmt"
    "github.com/kinslayere/eventtrackingbot/clients/telegram"
    "github.com/kinslayere/eventtrackingbot/types/telegram"
)

func processStartOrHelp(update types.Update) {

    log.Printf("Processing %s", update.Message.Text)

    text := fmt.Sprintf("Hello %s!", update.Message.From.FirstName)
    text += "\nThis is an event attendance tracking bot.\nAvailable commands are:"
    text += "\n/start - Start this bot"
    text += "\n/help - Info and commands"
    text += "\n/create_event [name] - Create a new event. Name may be provided inline"
    text += "\n/when [date] - View or set current event's date"
    text += "\n/where [place] - View or set current event's place"
    text += "\n/delete_event - Delete an existing event"
    text += "\n/current_event - Get current event being tracked"
    text += "\n/select_event - Change the current event"
    text += "\n/all_events - Get all events created in this group"
    text += "\n/participants - Get confirmed participants to current event"
    text += "\n/will_go - Confirm your attendance to current event"
    text += "\n/wont_go - Deny or remove your attendance to current event"

    smr := telegram.NewSendMessageRequest()
    smr.AddChatId(update.Message.Chat.Id)
    smr.AddText(text)
    if _, err := smr.Execute(); err != nil {
        log.Printf("Error sending response to %s: %v", update.Message.Text, err)
    }
}

func processPrivateChat(update types.Update) {

    text := fmt.Sprintf("Hello %s!\n", update.Message.From.FirstName)
    text += "Sorry, but this bot is intended for usage in group chats only."

    smr := telegram.NewSendMessageRequest()
    smr.AddChatId(update.Message.Chat.Id)
    smr.AddText(text)
    if _, err := smr.Execute(); err != nil {
        log.Printf("Error sending response to %s: %v", update.Message.Text, err)
    }
}
