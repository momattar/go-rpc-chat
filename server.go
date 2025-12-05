package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

type ChatMessage struct {
	UserID  string
	Message string
}

type ChatServer struct {
	clients map[string]chan string
	mutex   sync.Mutex
}

func (c *ChatServer) Join(userID string, reply *string) error {
	c.mutex.Lock()
	if _, exists := c.clients[userID]; !exists {
		c.clients[userID] = make(chan string, 10)
	}
	c.mutex.Unlock()

	go c.broadcast(fmt.Sprintf("User %s joined", userID), userID)
	*reply = "Welcome to the chat!"
	return nil
}

func (c *ChatServer) Send(msg ChatMessage, reply *string) error {
	c.broadcast(fmt.Sprintf("[%s]: %s", msg.UserID, msg.Message), msg.UserID)
	*reply = "Message sent!"
	return nil
}

func (c *ChatServer) Receive(userID string, messages *[]string) error {
	c.mutex.Lock()
	ch, ok := c.clients[userID]
	c.mutex.Unlock()
	if !ok {
		return fmt.Errorf("client not registered")
	}

	msgs := []string{}
	done := false
	for !done {
		select {
		case m := <-ch:
			msgs = append(msgs, m)
		default:
			done = true
		}
	}
	*messages = msgs
	return nil
}

func (c *ChatServer) broadcast(msg, senderID string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for id, ch := range c.clients {
		if id != senderID {
			ch <- msg
		}
	}
}

func main() {
	server := &ChatServer{
		clients: make(map[string]chan string),
	}

	rpc.Register(server)
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server listening on port 1234")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
