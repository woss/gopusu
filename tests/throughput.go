package main

import (
	"fmt"
	"log"
	"time"
	"github.com/PuSuEngine/gopusu"
)

func main() {
    pc, err := gopusu.NewClient("127.0.0.1", 55000)

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

	log.Println("Sending messages")

	messages := 100000

	start := time.Now()
	for i := 0; i < messages; i++ {
		if i % 500 == 0 {
			fmt.Print(".")
		}
		pc.Publish("channel.1", fmt.Sprintf("message %d", i))
	}
	fmt.Print("\n")
	since := time.Since(start)
	msec := since / time.Millisecond
	duration := since / time.Duration(messages)
	rate := int64(time.Second / duration)

	log.Printf("Sent %d messages in %d msec", messages, int64(msec))
	log.Printf("%d usec/message", int64(duration / time.Microsecond))
	log.Printf("%d messages/sec", rate)

}
