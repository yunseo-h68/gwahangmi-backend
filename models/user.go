package models

import (
	"net/http"

	"github.com/mholt/binding"
)

// User 구조체는 사용자에 대한 정보를 담습니다
type User struct {
	ID         interface{} `bson:"_id" json:"id"`
	UID        string      `bson:"uid" json:"uid"`
	Pw         string      `bson:"pw" json:"pw"`
	Name       string      `bson:"uname" json:"uname"`
	ProfileImg string      `bson:"profileImg" json:"profileImg"`
	Point      int         `bson:"point" json:"point"`
	PostCnt    int         `bson:"postCnt" json:"postCnt"`
}

// Login 구조체는 로그인에 필요한 정보를 담습니다
type Login struct {
	UID string `bson:"uid" json:"uid"`
	Pw  string `bson:"pw" json:"pw"`
}

// FieldMap 메서드는 User 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (u *User) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&u.ID:         "_id",
		&u.UID:        "uid",
		&u.Pw:         "pw",
		&u.Name:       "uname",
		&u.ProfileImg: "profileImg",
		&u.Point:      "point",
		&u.PostCnt:    "postCnt",
	}
}

// FieldMap 메서드는 Login 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (l *Login) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&l.UID: "uid",
		&l.Pw:  "pw",
	}
}
