package main

import (
	"fmt"
	"log"
	"time"
	"github.com/lietu/gopusu"
)

func main() {
    pc, err := gopusu.NewPuSuClient("127.0.0.1", 55000)

	if err != nil {
    	log.Println(err)
		log.Fatalf("Failed to create PuSuClient\n")
	}

    defer pc.Close()

    err = pc.Authorize("foo")

	if err != nil {
    	log.Println(err)
		log.Fatalf("Failed to authorize\n")
	}

	/*
    pc.Subscribe("channel.1", readMessage)

	if err != nil {
    	log.Println(err)
		log.Fatalf("Failed to subscribe\n")
	}
	*/

	log.Println("Sending message")

	messages := 10000

	start := time.Now()
	for i := 0; i < messages; i++ {
		pc.Publish("channel.1", fmt.Sprintf("message %d", i))
	}
	since := time.Since(start)
	msec := since / time.Millisecond
	duration := since / time.Duration(messages)
	rate := int64(time.Second / duration)

	log.Printf("Sent %d messages in %d msec", messages, int64(msec))
	log.Printf("%d usec/message", int64(duration / time.Microsecond))
	log.Printf("%d messages/sec", rate)

	log.Println("Waiting to receive")
	time.Sleep(time.Second * 5)
}

func readMessage(message *gopusu.Publish) {
	log.Printf("%s: %s\n", message.Channel, message.Content)
}