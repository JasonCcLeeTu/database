package db

import (
	"database/entities/user"
	repoDB "database/repository/db"
	"fmt"
)

type PracticeDB interface {
	InsertUSR(interface{}) error
	QueryUSR(interface{}) ([]user.User, error)
	InsertUSRTMP(interface{}) error
	QueryUSRListByName(name string) ([]user.User, error)
}

type practiceDB struct {
	ConDB   repoDB.ConnectDB
	RedisDB repoDB.CacheDB
}

func NewPracticeDB(conDB repoDB.ConnectDB, redisDB repoDB.CacheDB) PracticeDB {
	return &practiceDB{
		ConDB:   conDB,
		RedisDB: redisDB,
	}
}

func (p *practiceDB) InsertUSR(usr interface{}) error {
	if err := p.ConDB.InsertData(usr); err != nil {
		return err
	}
	return nil
}

func (p *practiceDB) QueryUSR(usr interface{}) ([]user.User, error) {

	usrList, err := p.ConDB.QueryData(usr)
	if err != nil {
		return nil, err
	}
	list, ok := usrList.([]user.User)
	if !ok {
		return nil, err
	}
	return list, nil
}

func (p *practiceDB) InsertUSRTMP(usr interface{}) error {
	user, ok := usr.(user.User)
	if !ok {
		return fmt.Errorf("assert user error")
	}
	if err := p.RedisDB.InsertUser(user); err != nil {
		return err
	}

	return nil
}

func (p *practiceDB) QueryUSRListByName(name string) ([]user.User, error) {
	list, err := p.RedisDB.QueryUserList()
	var usrList []user.User
	if err != nil {
		return nil, err
	}
	for _, usr := range list {
		if usr.Name == name {
			usrList = append(usrList, usr)
		}
	}
	return usrList, nil
}
