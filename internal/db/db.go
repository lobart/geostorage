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
		Connect()
		Close()
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

func New(ps *pubsub.Pubsub) DriverDB {
	f, err := os.Open("config/config_db.yml")
	if err != nil {
		fmt.Print("Error ", err)
	}
	defer f.Close()

	var cfg models.DBConfig
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Print("Error ", err)
	} else {
		fmt.Print(cfg.Database.Type)
	}

	fmt.Printf("Creating DB driver for %v", cfg.Database.Type)
	switch cfg.Database.Type{
	case "PostgreSQL":
		db:= postgres.PostgreSqlDriver{Cfg: &cfg}
		db.Connect()
		go listen(ps, db.Push)
		return &db
	case "MySQL":
		db := mysql.MySqlDriver{Cfg: &cfg}
		db.Connect()
		go listen(ps, db.Push)
		return &db
	case "MongoDB":
		db:= mongo.MongoDriver{Cfg: &cfg}
		db.Connect()
		go listen(ps, db.Push)
		return &db
	}
	return nil
}
