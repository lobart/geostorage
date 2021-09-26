package mongo

import (
	"context"
	"fmt"
	"github.com/lobart/go_geoserver.git/internal/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

const CONNECTED = "Successfully connected to database: %v"

type MongoDriver struct {
	Cfg models.DBConfig
	db  *mongo.Database
	Session *mongo.Client
}

func (d *MongoDriver) Connect() {
	var connectOnce sync.Once
	var err error
	connectOnce.Do(func() {

		mongoInfo:=fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",d.Cfg.Database.Username, d.Cfg.Database.Password, d.Cfg.Server.Host, d.Cfg.Server.Port, d.Cfg.Database.DBName)

		clientOptions := options.Client().ApplyURI(mongoInfo)

		d.Session, err = mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			panic(err)
		}
		fmt.Println("\nConnection is success! ")
		d.db = d.Session.Database(d.Cfg.Database.DBName)
	})
}


func (d *MongoDriver) Close() {
	d.Session.Disconnect(context.TODO())
}


func (d *MongoDriver) Push(kick *models.KickConfig) error {
	// Set client options
	collection := d.db.Collection("kicks")
	insertResult, err := collection.InsertOne(context.TODO(), kick)
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return nil
}