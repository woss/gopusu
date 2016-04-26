// The Golang client for PuSu Engine. PuSu Engine is a
// (relatively) fast and scalable Pub-Sub message delivery
// system.
//
// The Golang client is a simple synchronous client that
// does little magic internally. For operations that the
// server acknowledges (Connect, Authorize, Subscribe) it
// waits for the appropriate event coming back from the
// server before continuing, to ensure nothing odd happens.
package gopusu

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/PuSuEngine/pusud/messages"
	"net/url"
	"time"
)

const debug = false
const default_timeout = time.Second * 5

// Timeout exceeded when waiting to acknowledge Authorize/Subscribe request
var ErrTimeoutExceeded = errors.New("Timeout exceeded when waiting for server to acknowledge message")

// Callback to call with the published messages in a
// channel we're subscribed to.
type SubscribeCallback func(*Publish)
type subscribers map[string]SubscribeCallback

// The PuSu client. Create one of these with NewClient(),
// and call the Authorize(), Subscribe() and Publish()
// methods to communicate with the PuSu network.
type Client struct {
	connection  *websocket.Conn
	server      string
	subscribers subscribers
	Timeout     time.Duration
	waiting     bool
	waitingCh   chan string
	connected   bool
}

// Published message from the PuSu network. You get these
// to your callback if you subscribe to a channel.
type Publish struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Content string `json:"content"`
}

// Claim you have authorization to access some things.
// The server will determine what those things could be
// based on the configured authenticator and the data you
// give. Expect to get disconnected if this is invalid.
func (pc *Client) Authorize(authorization string) error {
	err := pc.sendMessage(&messages.Authorize{messages.TYPE_AUTHORIZE, authorization})

	if err != nil {
		return err
	}

	err = pc.wait(messages.TYPE_AUTHORIZATION_OK)

	if err != nil {
		return err
	}

	return nil
}

// Ask to subscribe for messages on the given channel. You
// MUST use the Authorize() method before this, even if
// server is configured to use the "None" authenticator.
// Expect to get disconnected if you lack the permissions.
func (pc *Client) Subscribe(channel string, callback SubscribeCallback) error {
	pc.subscribers[channel] = callback

	err := pc.sendMessage(&messages.Subscribe{messages.TYPE_SUBSCRIBE, channel})

	if err != nil {
		return err
	}

	err = pc.wait(messages.TYPE_SUBSCRIBE_OK)

	if err != nil {
		return err
	}

	return nil
}

// Publish a message on the given channel
func (pc *Client) Publish(channel string, content string) error {
	// TODO: Support interface{} for content
	return pc.sendMessage(&messages.Publish{messages.TYPE_PUBLISH, channel, content})
}

// Disconnect from the server
func (pc *Client) Close() {
	pc.connected = false
	pc.connection.Close()
}

func (pc *Client) sendMessage(message messages.Message) (err error) {
	data := message.ToJson()

	if debug {
		fmt.Printf("-> %s\n", data)
	}

	err = pc.connection.WriteMessage(websocket.TextMessage, data)

	return
}

func (pc *Client) disconnected() {
	pc.connected = false
	if debug {
		fmt.Printf("Disconnected from server.\n")
	}
}

func (pc *Client) receive(data []byte) {
	if debug {
		fmt.Printf("<- %s\n", data)
	}

	m := messages.GenericMessage{}
	json.Unmarshal(data, &m)

	if pc.waiting {
		pc.waitingCh <- m.Type
	}

	if m.Type == messages.TYPE_PUBLISH {
		p := Publish{}
		json.Unmarshal(data, &p)

		callback, ok := pc.subscribers[p.Channel]
		if ok {
			callback(&p)
		}
	}
}

func (pc *Client) wait(eventType string) (err error) {
	pc.waiting = true
	defer func() {
		pc.waiting = false
	}()

	if debug {
		fmt.Printf("Waiting for %s\n", eventType)
	}

	select {
	case t := <-pc.waitingCh:
		// Since all the operations that wait are
		// synchronous, we don't have to check channel.
		if t == eventType {
			if debug {
				fmt.Printf("Got %s\n", t)
			}
			return
		}
	case <-time.After(pc.Timeout):
		err = ErrTimeoutExceeded
		return
	}

	return
}

// Connect to the server
func (pc *Client) Connect() (err error) {

	if debug {
		fmt.Printf("Connecting to %s\n", pc.server)
	}

	c, _, err := websocket.DefaultDialer.Dial(pc.server, nil)

	if err != nil {
		return
	}

	pc.connection = c

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				pc.disconnected()
				return
			}
			pc.receive(message)
		}
	}()

	err = pc.wait(messages.TYPE_HELLO)

	if err != nil {
		pc.connected = true
	}

	return
}

// Create a new PuSu client and connect to the given server
func NewClient(host string, port int) (pc *Client, err error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}

	pc = &Client{}
	pc.waitingCh = make(chan string)
	pc.subscribers = subscribers{}
	pc.Timeout = default_timeout
	pc.server = u.String()

	err = pc.Connect()

	return
}
