package gopusu

import (
	"fmt"
	"time"
	"errors"
	"net/url"
	"github.com/gorilla/websocket"
	"encoding/json"
)

const DEFAULT_TIMEOUT = time.Second * 5

var errConnectionFailed = errors.New("Connection failed")
var errTimeoutExceeded = errors.New("Timeout exceeded when waiting for server to acknowledge message")

type SubscribeCallback func(*Publish)
type Subscribers map[string]SubscribeCallback

type PuSuClient struct {
	Connection  *websocket.Conn
	Subscribers Subscribers
	Timeout     time.Duration
	waiting		string
}

func (pc *PuSuClient) Authorize(authorization string) error {
	err := pc.SendMessage(NewAuthorize(authorization))

	if err != nil {
		return err
	}

	err = pc.wait(TYPE_AUTHORIZATION_OK)

	if err != nil {
		return err
	}

	return nil
}

func (pc *PuSuClient) Subscribe(channel string, callback SubscribeCallback) error {
	pc.Subscribers[channel] = callback

	err := pc.SendMessage(NewSubscribe(channel))

	if err != nil {
		return err
	}

	err = pc.wait(TYPE_SUBSCRIBE_OK)

	if err != nil {
		return err
	}

	return nil
}

func (pc *PuSuClient) Publish(channel string, content string) error {
	return pc.SendMessage(NewPublish(channel, content))
}

func (pc *PuSuClient) Close() {
	pc.Connection.Close()
}

func (pc *PuSuClient) SendMessage(message Message) error {
	data, err := message.ToJson()

	if err != nil {
		return err
	}

//	fmt.Printf("-> %s\n", data)

	pc.Connection.WriteMessage(websocket.TextMessage, data)

	return nil
}

func (pc *PuSuClient) Disconnected() {

}

func (pc *PuSuClient) Receive(data []byte) {
//	fmt.Printf("<- %s\n", data)

	im := IncomingMessage{}
	json.Unmarshal(data, &im)

	if im.Type == pc.waiting {
		pc.waiting = ""
	}

	if im.Type == TYPE_PUBLISH {
		p := Publish{}
		json.Unmarshal(data, &p)

		callback, ok := pc.Subscribers[p.Channel]
		if ok {
			callback(&p)
		}
	}
}

func (pc *PuSuClient) wait(eventType string) error {
	pc.waiting = eventType
	start := time.Now()

//	fmt.Printf("Waiting for %s\n", eventType)

	for {
		if pc.waiting == "" {
			return nil
		}

		time.Sleep(time.Millisecond)

		since := time.Since(start)
		if since > pc.Timeout {
			return errTimeoutExceeded
		}
	}

	return nil
}

func NewPuSuClient(host string, port int) (*PuSuClient, error) {
	pc := PuSuClient{}
	pc.Subscribers = Subscribers{}
	pc.Timeout = DEFAULT_TIMEOUT

	addr := fmt.Sprintf("%s:%d", host, port)
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}

//	fmt.Printf("%s\n", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		return nil, errConnectionFailed
	}

	pc.Connection = c

	go func() {
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				pc.Disconnected()
				return
			}
			pc.Receive(message)
		}
	}()

	pc.wait(TYPE_HELLO)

	return &pc, nil
}

