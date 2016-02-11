package requests

import (
	"errors"
	"github.com/kinslayere/eventtrackingbot/types"
	"github.com/kinslayere/eventtrackingbot/global"
	"fmt"
)

type SendMessageResponse struct {
	Ok			bool			`json:"ok"`
	Result		types.Message	`json:"result"`
}

type SendMessageRequest struct {
	getRequest *GetRequest
}

func NewSendMessageRequest() *SendMessageRequest {
	getRequest := NewGetRequest()
	getRequest.baseUrl = global.BASE_URL + "sendMessage"
	return &SendMessageRequest{getRequest}
}

func (r *SendMessageRequest) AddChatId(chatId int64) {
	r.getRequest.SetParamInt64("chat_id", chatId)
}

func (r *SendMessageRequest) AddText(text string) {
	r.getRequest.SetParamString("text", text)
}

func (r *SendMessageRequest) AddParseMode(parseMode string) {
	r.getRequest.SetParamString("parse_mode", parseMode)
}

func (r *SendMessageRequest) AddDisableWebPagePreview(disableWebPagePreview bool) {
	r.getRequest.SetParamBoolean("disable_web_page_preview", disableWebPagePreview)
}

func (r *SendMessageRequest) AddReplyToMessageId(replyToMessageId int64) {
	r.getRequest.SetParamInt64("reply_to_message_id", replyToMessageId)
}

func (r *SendMessageRequest) AddForceReply(forceReply types.ForceReply) {
	r.getRequest.SetParamStringer("reply_markup", forceReply)
}

func (r *SendMessageRequest) Execute() (*SendMessageResponse, error) {

	if !r.getRequest.HasParam("chat_id") {
		return nil, errors.New("chat_id is required")
	}

	if !r.getRequest.HasParam("text") {
		return nil, errors.New("text is required")
	}

	var response SendMessageResponse
	err := r.getRequest.Execute(&response)
	if err != nil {
		return nil, err
	}

	if !response.Ok {
		return nil, errors.New(fmt.Sprintf("Error executing request '%v'", r.getRequest.baseUrl))
	}

	return &response, nil
}