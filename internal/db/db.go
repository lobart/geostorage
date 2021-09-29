package db

import (
	"fmt"
	"github.com/lobart/go_geoserver.git/internal/db/mongo"
	"github.com/lobart/go_geoserver.git/internal/db/mysql"
	"github.com/lobart/go_geoserver.git/internal/db/postgres"
	"github.com/lobart/go_geoserver.git/internal/models"
	"github.com/lobart/go_geoserver.git/internal/pubsub"
	"gopkg.in/yaml.v2"
	"os"
)

type ( DriverDB interface {
		Connect() error
		Close() error
		Push(kick *models.KickConfig) error
	}
)

func listen(ps *pubsub.Pubsub, f func(k *models.KickConfig) error){
	ch := ps.Subscribe("kick")
	for msg := range ch{
		fmt.Println("Getting message from channel ")
		f(&msg)
	}

}

func New(ps *pubsub.Pubsub) (DriverDB, error) {
	f, err := os.Open("/home/archi/Golang_example/geostorage/config/config_db.yml")
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

	fmt.Printf("Creating DB driver for %v", cfg.Database.Type)
	switch cfg.Database.Type{
	case "PostgreSQL":
		db:= postgres.PostgreSqlDriver{Cfg: &cfg}
		err = db.Connect()
		if err != nil {
			return nil, err
		}
		go listen(ps, db.Push)
		return &db, nil
	case "MySQL":
		db := mysql.MySqlDriver{Cfg: &cfg}
		err =  db.Connect()
		if err != nil {
			return nil, err
		}
		go listen(ps, db.Push)
		return &db, nil
	case "MongoDB":
		db:= mongo.MongoDriver{Cfg: &cfg}
		err =  db.Connect()
		if err != nil {
			return nil, err
		}
		go listen(ps, db.Push)
		return &db, nil
	default:
		return nil, nil
	}
}
