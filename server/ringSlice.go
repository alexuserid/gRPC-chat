package main

import (
	"sync"
)

const ringSize = 10

type message struct {
	messageId string
	name string
	text string
}

type ringSlice struct {
	history []message
	mutex sync.RWMutex
}

func NewRingSlice() *ringSlice {
	return &ringSlice{
		history: make([]message, ringSize),
	}
}

func (r *ringSlice) AddMessage(msg message) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.history = append(r.history[1:], msg)
}

func (r *ringSlice) GetAllMessages() []message {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	return r.history
}
