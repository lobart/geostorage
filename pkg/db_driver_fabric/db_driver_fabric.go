package db_driver_fabric

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/lobart/go_geoserver.git/pkg/models"
	"gopkg.in/yaml.v2"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type driverDB interface {
	Push(kick models.KickConfig)
}

type postgreSqlDriver struct {
	cfg models.DBConfig
}

type mySqlDriver struct {
	cfg models.DBConfig
}

type mongoDriver struct {
	cfg models.DBConfig
}


func (d postgreSqlDriver) Push(kick models.KickConfig) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		d.cfg.Server.Host, d.cfg.Server.Port, d.cfg.Database.Username, d.cfg.Database.Password, d.cfg.Database.DBName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")

	sqlStatement := `
		INSERT INTO kicks (companyname, kickname, longitude, latitude, speed, status)
		VALUES ($1, $2, $3, $4, &5, &6)`
	_, err = db.Exec(sqlStatement, kick.CompanyName, kick.KickName, kick.Longitude, kick.Latitude, kick.Speed, kick.Status )
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully push!")
}

func (d mySqlDriver) Push(kick models.KickConfig) {
	mysqlInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		d.cfg.Database.Username, d.cfg.Database.Password, d.cfg.Server.Host, d.cfg.Server.Port , d.cfg.Database.DBName)


	db, err := sql.Open("mysql", mysqlInfo)

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	sqlStatement := fmt.Sprintf("INSERT INTO kicks (companyname, kickname, longitude, latitude, speed, status) VALUES (%s, %s, %s, %s, %s, %s)", kick.CompanyName, kick.KickName, kick.Longitude, kick.Latitude, kick.Speed, kick.Status)
	// perform a db.Query insert
	insert, err := db.Query(sqlStatement)

	// if there is an error inserting, handle it
	if err != nil {
		panic(err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

	fmt.Println("Successfully push!")
}

func (d mongoDriver) Push(kick models.KickConfig) {
	// Set client options
	mongoInfo:=fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",d.cfg.Database.Username, d.cfg.Database.Password, d.cfg.Server.Host, d.cfg.Server.Port, d.cfg.Database.DBName)

	clientOptions := options.Client().ApplyURI(mongoInfo)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database(d.cfg.Database.DBName).Collection("kicks")


	insertResult, err := collection.InsertOne(context.TODO(), kick)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

}

type DriverCreator struct{}

func (dc DriverCreator) CreateDriver() driverDB {
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
		return postgreSqlDriver{cfg:cfg}
	case "MySQL":
		return mySqlDriver{cfg:cfg}
	case "MongoDB":
		return mongoDriver{cfg:cfg}
	}
	return nil
}
