package user

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
	ProfileImg string      `bson:"profile_img" json:"profile_img"`
}

// Login 구조체는 로그인에 필요한 정보를 담습니다
type Login struct {
	UID string `bson:"uid" json:"uid"`
	Pw  string `bson:"pw" json:"pw"`
}

// FieldMap 메서드는 User 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (u *User) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&u.UID:  "uid",
		&u.Pw:   "pw",
		&u.Name: "uname",
	}
}

// FieldMap 메서드는 Login 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (l *Login) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&l.UID: "uid",
		&l.Pw:  "pw",
	}
}

// New 함수는 새로운 User 구조체를 생성합니다.
func New(uid, pw, name, profileImg string) (*User, error) {
	u := new(User)
	u.UID = uid
	u.Pw = pw
	u.Name = name
	u.ProfileImg = profileImg
	return u, nil
}
