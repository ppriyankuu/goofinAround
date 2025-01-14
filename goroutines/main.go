package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Message struct {
	chats   []string
	friends []string
}

func main() {
	now := time.Now()

	id := getUserId("priyanku")
	println(id)

	wg := &sync.WaitGroup{}
	ch := make(chan *Message, 2)

	wg.Add(2)

	go getUserChats(id, ch, wg)
	go getUserFriends(id, ch, wg)

	wg.Wait()
	close(ch)

	for msg := range ch {
		log.Println(msg)
	}

	log.Println(time.Since(now))
}

func getUserFriends(_ string, ch chan<- *Message, wg *sync.WaitGroup) {
	time.Sleep(time.Second * 2)

	ch <- &Message{
		friends: []string{
			"John",
			"someone",
			"no one",
			"anonymous",
			"that one",
		},
	}

	wg.Done()
}

func getUserChats(_ string, ch chan<- *Message, wg *sync.WaitGroup) {
	time.Sleep(time.Second * 1)

	ch <- &Message{
		chats: []string{
			"John",
			"someone",
			"no one",
		},
	}
	wg.Done()
}

func getUserId(name string) string {
	time.Sleep(time.Second * 1)
	return fmt.Sprintf("%s-2", name)
}
