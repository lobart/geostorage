package pubsub

import (
	"fmt"
	"github.com/lobart/go_geoserver.git/internal/models"
	"sync"
)

type Pubsub struct {
	mu   sync.RWMutex
	subs map[string][]chan models.KickConfig
	closed bool
	nBuf int
}

func (ps Pubsub) New() *Pubsub{
	return &Pubsub{nBuf: 10, subs: map[string][]chan models.KickConfig{}}
}


func (ps *Pubsub) Subscribe(topic string) <-chan models.KickConfig {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan models.KickConfig, ps.nBuf)
	ps.subs[topic] = append(ps.subs[topic], ch)
	fmt.Println("Success subscribing", ps.subs["kick"])
	return ch
}

func (ps *Pubsub) Publish(topic string, msg models.KickConfig) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}
	fmt.Println("Try publish message to channel ")
	for _, ch := range ps.subs[topic] {
		fmt.Println("In cycle")
		go func(ch chan models.KickConfig) {
			fmt.Println("Sending message to channel ")
			ch <- msg
		}(ch)
	}
}


