package mongoconnection

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//MongoConnection
type MongoConnection struct {
	DB *mongo.Client
}

var once = sync.Once{}

//NewMongoConnection 取得新連線
func NewMongoConnection(username, password, address string) *MongoConnection {
	if username == "" || password == "" || address == "" {
		log.Panic().Msgf("NewMongoConnection fail: username: %s, password: %s, address: %s", username, password, address)
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s",
		username,
		password,
		address)

	//only print one time
	once.Do(func() {
		log.Info().Msgf(mongoURI)
	})

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Panic().Msgf("set mongo config error")
	}

	return &MongoConnection{
		DB: client,
	}
}

func (l *MongoConnection) Ping() {
	ctx := context.Background()
	err := l.DB.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Panic().Msgf("set mongo config error")
	}
}

func (l *MongoConnection) TestConnection(input int) {
	ctx := context.Background()

	collection := l.DB.Database("alertdb").Collection("testTimeout")
	if collection == nil {
		log.Panic().Msgf("get collection  error")
	}

	result, err1 := collection.InsertOne(ctx, bson.D{{Key: "name", Value: "pi"}, {Key: "value", Value: input}})
	if err1 != nil {
		log.Panic().Msgf("set mongo config error: ", err1.Error())
	}

	log.Info().Msgf("%v", result)
}
