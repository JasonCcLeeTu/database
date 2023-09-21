package main

import (
	"database/entities/user"
	"database/repository/db"
	usecasedb "database/usecase/db"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	initEngine()
	//QueryData()
	fmt.Println("完成")

}

func initEngine() {

	cred := options.Credential{
		Username:   "jason",
		Password:   "123456",
		AuthSource: "admin",
	}
	uri := "mongodb://127.0.0.1:27017"

	redisOption := redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	}

	mongObj, err := db.NewMongoDB(uri, cred) //mongodb repository
	redisObj, err := db.NewCacheDB(redisOption)

	if err != nil {
		log.Fatalf("建立MongoDB物件失敗:%s \n", err)
	}

	prac := usecasedb.NewPracticeDB(mongObj, redisObj) // usecase practice insertdb

	for i := 1; i < 6; i++ {
		user := user.User{
			Name:            fmt.Sprintf("test%d", i),
			PhoneNum:        9123456780,
			CreateDate:      primitive.NewDateTimeFromTime(time.Now()),
			CreateTimeStamp: int64(time.Now().UnixNano()),
		}

		config := db.MongoConfig{
			Database:   "golang_testdb",
			Collection: "data2",
			Type: map[db.OptionsType]interface{}{
				db.INSERTONE: user,
			},
		}

		if err = prac.InsertUSR(config); err != nil {
			log.Fatalf("Insert User Data error:%s", err)
		}
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("完成")

	defer func() {
		if err = mongObj.RelCon(); err != nil {
			log.Fatalf(err.Error())
		}

	}()
}

func QueryData() {
	cred := options.Credential{
		Username:   "jason",
		Password:   "123456",
		AuthSource: "admin",
	}
	uri := "mongodb://127.0.0.1:27017"

	redisOption := redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0,
	}

	mongObj, err := db.NewMongoDB(uri, cred) //mongodb repository
	redisObj, err := db.NewCacheDB(redisOption)
	defer func() {
		if err = mongObj.RelCon(); err != nil {
			log.Fatalf(err.Error())
		}

	}()

	if err != nil {
		log.Fatalf("建立MongoDB物件失敗:%s \n", err)
	}

	prac := usecasedb.NewPracticeDB(mongObj, redisObj) // usecase practice insertdb

	config := db.MongoConfig{
		Database:   "golang_testdb",
		Collection: "data2",
		Type: map[db.OptionsType]interface{}{
			db.QUERY: nil,
		},
	}

	list, err := prac.QueryUSR(config)
	if err != nil {
		log.Fatalf("Insert User Data error:%s", err)
	}

	for _, val := range list {
		fmt.Printf("%+v \n", val)
	}

}
