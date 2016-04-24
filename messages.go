package gopusu

import (
	"encoding/json"
)

const TYPE_HELLO = "hello"
const TYPE_AUTHORIZE = "authorize"
const TYPE_AUTHORIZATION_OK = "authorization_ok"
const TYPE_PUBLISH = "publish"
const TYPE_SUBSCRIBE = "subscribe"
const TYPE_SUBSCRIBE_OK = "subscribe_ok"
const TYPE_UNKNOWN_MESSAGE_RECEIVED = "unknown_message_received"
const TYPE_AUTHORIZATION_FAILED = "authorization_failed"
const TYPE_PERMISSION_DENIED = "permission_denied"


type Message interface {
	ToJson() ([]byte, error)
}


type IncomingMessage struct {
	Type string `json:"type"`
}


type Authorize struct {
	Type string `json:"type"`
	Authorization string `json:"authorization"`
}

func (a *Authorize) ToJson() ([]byte, error) {
	result, err := json.Marshal(&a)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewAuthorize(authorization string) Message {
	a := Authorize{}
	a.Type = TYPE_AUTHORIZE
	a.Authorization = authorization
	return &a
}


type Subscribe struct {
	Type string `json:"type"`
	Channel string `json:"channel"`
}

func (a *Subscribe) ToJson() ([]byte, error) {
	result, err := json.Marshal(&a)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewSubscribe(channel string) Message {
	s := Subscribe{}
	s.Type = TYPE_SUBSCRIBE
	s.Channel = channel
	return &s
}


type Publish struct {
	Type string `json:"type"`
	Channel string `json:"channel"`
	Content string `json:"content"`
}

func (a *Publish) ToJson() ([]byte, error) {
	result, err := json.Marshal(&a)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func NewPublish(channel string, content string) Message {
	s := Publish{}
	s.Type = TYPE_PUBLISH
	s.Channel = channel
	s.Content = content
	return &s
}

