package driverFabric

import (
	"fmt"
	"github.com/lobart/go_geoserver.git/pkg/models"
	"gopkg.in/yaml.v2"
	"os"
	"github.com/lobart/go_geoserver.git/pkg/mongo_driver"
	"github.com/lobart/go_geoserver.git/pkg/mysql_driver"
	"github.com/lobart/go_geoserver.git/pkg/postgresql_driver"
)

type DriverDB interface {
	Connect()
	Close()
	Push(kick models.KickConfig)
}

func CreateDriver() DriverDB {
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
		return postgresqlDriver.PostgreSqlDriver{Cfg:cfg}
	case "MySQL":
		return mysqlDriver.MySqlDriver{Cfg:cfg}
	case "MongoDB":
		return mongoDriver.MongoDriver{Cfg:cfg}
	}
	return nil
}
