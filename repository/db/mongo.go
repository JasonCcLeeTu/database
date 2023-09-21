package db

import (
	"context"
	"database/entities/user"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ConnectDB interface {
	InsertData(interface{}) error
	QueryData(interface{}) (interface{}, error)
	RelCon() error
}

type OptionsType string //db操作類型

const (
	INSERTONE OptionsType = "INSERT ONE"
	QUERY     OptionsType = "QUERY"
)

type MongoConfig struct {
	Database   string
	Collection string
	Type       map[OptionsType]interface{}
}

func NewMongoDB(uri string, credential options.Credential) (ConnectDB, error) {

	clientOptions := options.Client().ApplyURI(uri).SetAuth(credential)
	mgClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, fmt.Errorf("momgo connect error: %s", err.Error())
	}
	err = mgClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		return nil, fmt.Errorf("momgoClient ping error: %s", err.Error())
	}

	return &mongoDB{
			connClient: mgClient,
		},
		nil
}

type mongoDB struct {
	connClient *mongo.Client
}

func (mgDB *mongoDB) RelCon() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := mgDB.connClient.Disconnect(ctx); err != nil { // 釋放此連線物件
		return fmt.Errorf("釋放連線失敗:%s \n", err)
	}
	return nil
}

func (mgDB *mongoDB) InsertData(config interface{}) error {

	data, ok := config.(MongoConfig)

	if !ok {
		return fmt.Errorf("斷言 config 失敗")
	}

	collection := mgDB.connClient.Database(data.Database).Collection(data.Collection)

	for key, val := range data.Type {

		switch key {
		case INSERTONE:
			collection.InsertOne(context.Background(), val)
		}
	}

	return nil
}

func (mgDB *mongoDB) QueryData(config interface{}) (interface{}, error) {

	ctx := context.Background()
	var dataList []user.User
	data, ok := config.(MongoConfig)

	if !ok {
		return nil, fmt.Errorf("斷言 config 失敗")
	}

	for key, _ := range data.Type {
		switch key {
		case QUERY:
			option := options.Find().SetSort(bson.D{{"createDate", -1}}).SetLimit(5)
			collection := mgDB.connClient.Database(data.Database).Collection(data.Collection)

			cursor, err := collection.Find(ctx, bson.D{}, option)
			defer cursor.Close(ctx)
			if err != nil {
				return nil, fmt.Errorf("collection find error:%s", err.Error())
			}

			for cursor.Next(ctx) {
				var elem user.User
				err = cursor.Decode(&elem)
				if err != nil {
					return nil, fmt.Errorf(" cursor Decode error:%s", err.Error())
				}

				dataList = append(dataList, elem)
			}

			if err = cursor.Err(); err != nil {
				return nil, err
			}
			if len(dataList) == 0 {
				return nil, nil
			}
		}
	}

	return dataList, nil
}
