package db

import (
	"fmt"
	"github.com/lobart/go_geoserver.git/internal/models"
	"github.com/lobart/go_geoserver.git/internal/pubsub"
	"testing"
	"time"
	"errors"
)

func TestNew(t *testing.T) {
	ps, _ := pubsub.New()
	configs := []Conf{
		Conf{path: "/home/archi/Golang_example/geostorage/config/test_mongo.yml"},
		Conf{path: "/home/archi/Golang_example/geostorage/config/test_mysql.yml"},
		Conf{path: "/home/archi/Golang_example/geostorage/config/test_postgres.yml"},
	}
	for _, c:= range(configs){
		fmt.Println("Config is ",c)
		_, err := New(ps, c)
		if err!=nil{
			t.Errorf("Error in creating DB : %v",err)
		}
	}
}

type TestsForListen struct {
	f []func()
}



func TestListen(t *testing.T) {
	tests := []func(t *testing.T){func(t *testing.T) {
		ps, _ := pubsub.New()
		data := models.KickConfig{KickName: "test"}
		go Listen(ps, func(k *models.KickConfig) error {
			if k.KickName != "test" {
				return errors.New("Another test name " + k.KickName)
			}
			fmt.Println("Goodd")
			return nil
		})
		time.Sleep(1)
		go ps.Publish("kick", data)
		err := error(nil)
		if err != nil {
			t.Errorf("Error in Listen function %v", err)
		}
	},
		func(t *testing.T) {
			ps, _ := pubsub.New()
			data := models.KickConfig{KickName: "test"}
			go Listen(ps, func(k *models.KickConfig) error {
				if k.KickName != "test1" {
					return errors.New("Another test name " + k.KickName)
				}
				fmt.Println("Goodd")
				return nil
			})
			time.Sleep(1)
			go ps.Publish("kick", data)
			err := error(nil)
			if err != nil {
				t.Errorf("Error in Listen function %v", err)
			}
		},
	}
	for _, test:=range(tests){
		t.Run("test",test)
	}

}







