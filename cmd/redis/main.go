package main

import (
	"database/entities/user"
	"database/repository/db"
	usecasedb "database/usecase/db"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var forever = make(chan struct{})

func main() {

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
	if err != nil {
		log.Fatalf("建立MongoDB物件失敗:%s \n", err)
	}

	redisObj, err := db.NewCacheDB(redisOption)
	if err != nil {
		log.Fatalf("建立RedisDB物件失敗:%s \n", err)
	}

	defer func() {
		if err = mongObj.RelCon(); err != nil {
			log.Fatalf(err.Error())
		}

		if err = redisObj.RelCon(); err != nil {
			log.Fatal(err.Error())
		}

		log.Println("db close")

	}()

	prac := usecasedb.NewPracticeDB(mongObj, redisObj) // usecase practice insertdb

	for i := 1; i < 6; i++ {
		user := user.User{
			Name:            fmt.Sprintf("test%d", 7),
			PhoneNum:        9123456780,
			CreateDate:      time.Now(),
			CreateTimeStamp: int64(time.Now().UnixNano()),
		}
		time.Sleep(500 * time.Millisecond)
		prac.InsertUSRTMP(user)
	}

	result, err := prac.QueryUSRListByName("test7")
	if err != nil {
		log.Fatalf("prac.QueryUSRListByName error:%s \n", err.Error())
	}

	for _, val := range result {
		fmt.Printf("%+v \n", val)
	}

	//<-forever
}

func LuaScript(client *redis.Client) {

	for i := 1; i < 2; i++ {
		usr := user.User{
			Name:            fmt.Sprintf("test%d", i),
			PhoneNum:        3979507435 + i,
			CreateTimeStamp: time.Now().UnixMilli(),
			CreateDate:      time.Now(),
		}

		byteUsr, err := json.Marshal(usr)
		if err != nil {
			log.Fatalf("json marshal error: %s", err.Error())
		}

		luaScript := redis.NewScript(`
			local amount = redis.call('llen',KEYS[1])
			if (amount >= 5) then
				local result = redis.call('rpop',KEYS[1])
				if result == false then
					return result
				end
				local result = redis.call('lpush',KEYS[1],ARGV[1]) 
				return result
			else
				local result = redis.call('lpush',KEYS[1],ARGV[1]) 
				return result
			end
		
		`)

		result, err := luaScript.Run(client, []string{"eventLog"}, byteUsr).Result()
		if err != nil {
			log.Fatalf(" luaScript.Run error : %s ", err.Error())
		}
		log.Println("luaScript result:", result)
	}

}

func LLen(client *redis.Client) {
	amount, err := client.LLen("eventLog").Result()
	if err != nil {
		log.Fatalf("LLen error:%s", err.Error())
	}
	fmt.Printf("eventLog amount:%d", amount)
}

func LPush(client *redis.Client) {
	for i := 5; i < 6; i++ {
		usr := user.User{
			Name:            fmt.Sprintf("test%d", i),
			PhoneNum:        3979507435 + i,
			CreateTimeStamp: time.Now().UnixMilli(),
			CreateDate:      time.Now(),
		}
		if val, err := json.Marshal(usr); err != nil {
			log.Fatalf("json marshal error:%s", err.Error())
		} else {
			client.LPush("eventLog", val)

		}

	}

}

func RPop(client *redis.Client) {
	_, err := client.RPop("eventLog").Result()
	if err != nil {
		log.Fatalf("RPop error:%s", err.Error())
	}
}

func Lrange(client *redis.Client) {
	list, err := client.LRange("eventLog", 0, -1).Result()
	if err != nil {
		log.Fatalf("lrange error:%s", err.Error())
	}

	for _, val := range list {
		var usr user.User
		json.Unmarshal([]byte(val), &usr)
		fmt.Printf("%+v \n", usr)
	}
}
