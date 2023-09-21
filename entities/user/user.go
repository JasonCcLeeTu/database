package user

type User struct {
	Name            string      `bson:"name,omitempty"`
	PhoneNum        int         `bson:"phone_numb,omitempty"`
	CreateDate      interface{} `bson:"createDate,omitempty"`
	CreateTimeStamp int64       `bson:"createTimeStamp,omitempty"`
}
