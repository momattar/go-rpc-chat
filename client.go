package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your user ID: ")
	userID, _ := reader.ReadString('\n')
	userID = strings.TrimSpace(userID)

	var reply string
	err = client.Call("ChatServer.Join", userID, &reply)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply)

	go func() {
		for {
			var messages []string
			err := client.Call("ChatServer.Receive", userID, &messages)
			if err != nil {
				log.Println(err)
			}
			for _, m := range messages {
				fmt.Println(m)
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	for {
		fmt.Print("> ")
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}

		var sendReply string
		err := client.Call("ChatServer.Send", struct {
			UserID  string
			Message string
		}{UserID: userID, Message: msg}, &sendReply)
		if err != nil {
			log.Println(err)
		}
	}
}
