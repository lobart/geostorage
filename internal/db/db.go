package db

import (
	"errors"
	"fmt"
	"github.com/lobart/go_geoserver.git/internal/db/mongo"
	"github.com/lobart/go_geoserver.git/internal/db/mysql"
	"github.com/lobart/go_geoserver.git/internal/db/postgres"
	"github.com/lobart/go_geoserver.git/internal/models"
	"github.com/lobart/go_geoserver.git/internal/pubsub"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)


type ( DriverDB interface {
		Connect() error
		Close() error
		Push(kick *models.KickConfig) error
	}
)

type Conf struct {
	path string `default: "/home/archi/Golang_example/geostorage/config/config_db.yml"`
}

func Listen(ps *pubsub.Pubs, f func(k *models.KickConfig) error) error{
	mu:=&sync.Mutex{}
	ch := ps.Subscribe("kick")
	var ans string
	fmt.Println("Listening ")
	for msg := range ch{
		fmt.Println("Getting message from channel ")
		err := f(&msg)
		if err!=nil{
			ans += "\n" + err.Error()
		}
	}
	mu.Unlock()
	err := errors.New(ans)
	if err!=nil{
		return err
	}
	return nil
}

func  getConfig(conf Conf) (*models.DBConfig, error) {
	f, err := os.Open(conf.path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg models.DBConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	} else {
		fmt.Print(cfg.Database.Type)
	}
	return &cfg, err
}

func New(ps *pubsub.Pubs, conf Conf) (DriverDB, error) {
	cfg, err := getConfig(conf)
	fmt.Printf("Creating DB driver for %v", cfg.Database.Type)
	switch cfg.Database.Type{
	case "PostgreSQL":
		db:= postgres.PostgreSqlDriver{Cfg: cfg}
		err = db.Connect()
		if err != nil {
			return nil, err
		}
		go Listen(ps, db.Push)
		return &db, nil
	case "MySQL":
		db := mysql.MySqlDriver{Cfg: cfg}
		err =  db.Connect()
		if err != nil {
			return nil, err
		}
		go Listen(ps, db.Push)
		return &db, nil
	case "MongoDB":
		db:= mongo.MongoDriver{Cfg: cfg}
		err =  db.Connect()
		if err != nil {
			return nil, err
		}
		go Listen(ps, db.Push)
		return &db, nil
	default:
		return nil, nil
	}
}
