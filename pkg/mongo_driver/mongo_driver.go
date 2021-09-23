package mongoDriver

import (
	"context"
	"fmt"
	"github.com/lobart/go_geoserver.git/pkg/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
)

const CONNECTED = "Successfully connected to database: %v"

type MongoDriver struct {
	Cfg models.DBConfig
	db      *mongo.Database
	Session *mongo.Client
}

func (d MongoDriver) Connect() {
	var connectOnce sync.Once
	var db *mongo.Database
	var session *mongo.Client
	connectOnce.Do(func() {
		db, session = d.connectToMongo()
	})
}

func (d MongoDriver) connectToMongo() (a *mongo.Database, b *mongo.Client) {

	mongoInfo:=fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",d.Cfg.Database.Username, d.Cfg.Database.Password, d.Cfg.Server.Host, d.Cfg.Server.Port, d.Cfg.Database.DBName)

	clientOptions := options.Client().ApplyURI(mongoInfo)

	session, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	var DB = session.Database(d.Cfg.Database.DBName)
	return DB, session
}


func (d MongoDriver) Close(){
	d.Session.Disconnect(context.TODO())
}


func (d MongoDriver) Push(kick models.KickConfig) {
	// Set client options
	collection := d.Session.Database(d.Cfg.Database.DBName).Collection("kicks")
	insertResult, err := collection.InsertOne(context.TODO(), kick)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

}