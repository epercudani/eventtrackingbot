package requests

import (
	"log"
	"errors"
	"bytes"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strconv"
	"net/url"
	"github.com/kinslayere/eventtrackingbot/types"
	"github.com/kinslayere/eventtrackingbot/global"
)

type SendMessageResponse struct {
	Ok			bool			`json:"ok"`
	Result		types.Message	`json:"result"`
}

type SendMessageRequest struct {
	params map[string]interface{}
}

func NewSendMessageRequest() *SendMessageRequest {
	var smr SendMessageRequest
	smr.params = make(map[string]interface{})
	return &smr
}

func (r *SendMessageRequest) AddChatIdInt(chatId int64) {
	r.params["chat_id"] = chatId
}

func (r *SendMessageRequest) AddText(text string) {
	r.params["text"] = text
}

func (r *SendMessageRequest) AddParseMode(parseMode string) {
	r.params["parse_mode"] = parseMode
}

func (r *SendMessageRequest) AddDisableWebPagePreview(disableWebPagePreview bool) {
	r.params["disable_web_page_preview"] = disableWebPagePreview
}

func (r *SendMessageRequest) AddReplyToMessageId(replyToMessageId int64) {
	r.params["reply_to_message_id"] = replyToMessageId
}

func (r *SendMessageRequest) AddForceReply(forceReply types.ForceReply) {
	r.params["reply_markup"] = forceReply.String()
}

func (r *SendMessageRequest) GetParamsString() string {
	var paramsString bytes.Buffer
	first := true

	for k, v := range r.params {
		var value string
		switch vtype := v.(type) {
		case int64:
			value = strconv.FormatInt(v.(int64), 10)
		case string:
			value = url.QueryEscape(v.(string))
		case types.Stringer:
			value = url.QueryEscape(vtype.String())
		default:
			value = ""
		}

		if first {
			paramsString.WriteString(fmt.Sprintf("%s=%v", k, value))
			first = false
		} else {
			paramsString.WriteString(fmt.Sprintf("&%s=%v", k, value))
		}
	}

	return paramsString.String()
}

func (r *SendMessageRequest) DoRequest() (messageSent types.Message, err error) {

	if _, ok := r.params["chat_id"]; !ok {
		return messageSent, errors.New("chat_id is required")
	}

	if _, ok := r.params["text"]; !ok {
		return messageSent, errors.New("text is required")
	}

	paramsString := r.GetParamsString()
	url := global.BASE_URL + "sendMessage?" + paramsString
	log.Printf("sendMessage:url: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Send Message response body: %s", body)

	var response SendMessageResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
	}

	messageSent = response.Result
	return
}