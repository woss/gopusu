package main

import (
	"github.com/PuSuEngine/gopusu"
	"log"
	"time"
)

func main() {
	log.Printf("Connecting to 127.0.0.1:55000")
    pc, err := gopusu.NewClient("127.0.0.1", 55000)

	if err != nil {
		log.Fatalf("error: %s", err)
	}

    defer pc.Close()
	log.Printf("Authorizing with 'foo'")
    pc.Authorize("foo")
	log.Printf("Subscribing to channel.1")
    pc.Subscribe("channel.1", listener)
	log.Printf("Sending message to channel.1")
    pc.Publish("channel.1", "message")
	log.Printf("Waiting for messages")
	time.Sleep(time.Second)
}

func listener(msg *gopusu.Publish) {
	log.Printf("Got message %s on channel %s", msg.Content, msg.Channel)
}
