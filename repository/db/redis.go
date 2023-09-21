package db

import (
	"database/entities/user"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

type CacheDB interface {
	InsertUser(user.User) error
	RelCon() error
	QueryUserList() ([]user.User, error)
}

type cacheDB struct {
	redisClient *redis.Client
}

func NewCacheDB(option redis.Options) (CacheDB, error) {
	client := redis.NewClient(&option)
	if _, err := client.Ping().Result(); err != nil {
		return nil, fmt.Errorf("redis ping err %s", err.Error())
	}
	log.Printf("redis connected successfully \n")
	return &cacheDB{
			redisClient: client,
		},
		nil

}

func (c *cacheDB) InsertUser(usr user.User) error {

	byteUsr, err := json.Marshal(usr)
	if err != nil {
		return fmt.Errorf("json marshal error: %s", err.Error())
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

	result, err := luaScript.Run(c.redisClient, []string{"eventLog"}, byteUsr).Result()
	if err != nil {
		return fmt.Errorf(" luaScript.Run error : %s ", err.Error())
	}

	log.Println("luaScript result:", result)

	return nil
}

func (c *cacheDB) QueryUserList() ([]user.User, error) {
	var userList []user.User
	list, err := c.redisClient.LRange("eventLog", 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("queryUser lrange error:%s", err.Error())
	}

	for _, val := range list {
		var usr user.User
		if err = json.Unmarshal([]byte(val), &usr); err != nil {
			return nil, err
		}
		userList = append(userList, usr)
	}

	return userList, nil
}

func (c *cacheDB) RelCon() error {
	if err := c.redisClient.Close(); err != nil {
		return fmt.Errorf("redis RelCon:%s", err.Error())
	}

	return nil
}
