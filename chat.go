package main

import (
	"chat"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type ChatMng struct {
	chat  *chat.Chat
	users map[int64]*User
	sync.RWMutex
}

type User struct {
	send    chan<- string
	receive <-chan string
	abort   chan struct{}
	access  time.Time
}

func (cm *ChatMng) Enter(uid int64, name, login string) (<-chan string, <-chan struct{}) {
	cm.Lock()
	defer cm.Unlock()

	if user, ok := cm.users[uid]; ok {
		log.Println("===========enter")
		close(user.abort)
		user.abort = make(chan struct{})
		return user.receive, user.abort
	}

	var nickname = login
	if name != "" {
		nickname = fmt.Sprintf("%s(%s)", name, login)
	}

	send, receive := cm.chat.AddUser(nickname)
	abort := make(chan struct{})
	cm.users[uid] = &User{send, receive, abort, time.Now()}
	return receive, abort
}
func (cm *ChatMng) Send(uid int64, msg string) error {
	user, ok := cm.users[uid]
	if !ok {
		return errors.New("user not found")
	}
	user.send <- msg
	return nil
}

func (cm *ChatMng) Leave(uid int64) {
	user, ok := cm.users[uid]
	if !ok {
		return
	}
	cm.Lock()
	delete(cm.users, uid)
	cm.Unlock()
	close(user.send)
}
func (cm *ChatMng) GetUserList() string {
	return cm.chat.GetList()
}

var chatmng *ChatMng

func init() {
	chatmng = &ChatMng{chat: chat.New(), users: make(map[int64]*User)}
}
