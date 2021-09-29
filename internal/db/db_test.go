package db

import (
	"fmt"
	"github.com/lobart/go_geoserver.git/internal/models"
	"github.com/lobart/go_geoserver.git/internal/pubsub"
	"sync"
	"testing"
	"errors"
)

func TestListen(t *testing.T) {
	ps, _ := pubsub.Pubs{}.New()
	ch := make(chan models.KickConfig, 10)
	topic:="kick"
	ps.Subs[topic] = append(ps.Subs[topic], ch)
	wg:=&sync.WaitGroup{}
	wg.Add(1)
	for _, ch := range ps.Subs[topic] {
		go func(wg *sync.WaitGroup, ch chan models.KickConfig) {
			defer wg.Done()
			fmt.Println("Sending message to channel ")
			ch <- models.KickConfig{KickName: "test"}
		}(wg, ch)
	}
	go Listen(ps, func(k *models.KickConfig) error {
		if k.KickName != "test" {
			return errors.New("Another test name " + k.KickName)
		}
		fmt.Println("Goodd")
		return nil
	})
	wg.Wait()
	err := error(nil)
	if err!=nil{
		t.Errorf("Error in Listen function %v",err)
	}
}

func TestNew(t *testing.T) {
	ps, _ := pubsub.Pubs{}.New()
	_, err := New(ps)
	if err!=nil{
		t.Errorf("Error in creating DB : %v",err)
	}
}


