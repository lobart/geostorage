package pubsub

import (
	"errors"
	"fmt"
	"github.com/lobart/go_geoserver.git/internal/models"
	"sync"
)

type Pubs struct {
	mu     sync.RWMutex
	Subs   map[string][]chan models.KickConfig
	closed bool
	nBuf int
}

func New() (*Pubs, error){
	ps := Pubs{nBuf: 10, Subs: map[string][]chan models.KickConfig{}}
	if &ps==nil {
		return nil, errors.New("Nil pointer")
	}
	return &ps, nil
}


func (ps *Pubs) Subscribe(topic string) <-chan models.KickConfig {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan models.KickConfig, ps.nBuf)
	ps.Subs[topic] = append(ps.Subs[topic], ch)
	fmt.Println("Success subscribing")
	return ch
}

func (ps *Pubs) Publish(topic string, msg models.KickConfig) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}
	fmt.Println("Try publish message to channel ")
	for _, ch := range ps.Subs[topic] {
		go func(ch chan models.KickConfig) {
			fmt.Println("Sending message to channel ")
			ch <- msg
		}(ch)
	}
}


