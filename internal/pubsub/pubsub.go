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
	return &Pubsub{nBuf: 10}
}


func (ps *Pubsub) Subscribe(topic string) <-chan models.KickConfig {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan models.KickConfig, ps.nBuf)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *Pubsub) Publish(topic string, msg models.KickConfig) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}

	for _, ch := range ps.subs[topic] {
		go func(ch chan models.KickConfig) {
			ch <- msg
		}(ch)
	}
}

func (ps *Pubsub) CheckMessages(ch <-chan models.KickConfig) {
	select {
	case msg := <-ch:
		s.Driver.Push(&msg)
	default:
		fmt.Println("No data in channel!")
	}
}
