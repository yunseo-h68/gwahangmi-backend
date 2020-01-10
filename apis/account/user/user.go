package user

import (
	"net/http"

	"github.com/mholt/binding"
	"golang.org/x/crypto/bcrypt"
)

// User 구조체는 사용자에 대한 정보를 담습니다
type User struct {
	ID         interface{} `bson:"_id" json:"id"`
	UID        string      `bson:"uid" json:"uid"`
	Pw         string      `bson:"pw" json:"pw"`
	Name       string      `bson:"uname" json:"uname"`
	ProfileImg string      `bson:"profile_img" json:"profile_img"`
	Point      int         `bson:"point" json:"point"`
	PostCnt    int         `bson:"post_cnt" json:"post_cnt"`
}

// Login 구조체는 로그인에 필요한 정보를 담습니다
type Login struct {
	UID string `bson:"uid" json:"uid"`
	Pw  string `bson:"pw" json:"pw"`
}

// FieldMap 메서드는 User 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (u *User) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&u.UID:        "uid",
		&u.Pw:         "pw",
		&u.Name:       "uname",
		&u.ProfileImg: "profile_img",
		&u.Point:      "point",
		&u.PostCnt:    "post_cnt",
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
func New() (*User, error) {
	u := new(User)
	u.UID = ""
	u.Pw = ""
	u.Name = ""
	u.ProfileImg = ""
	u.Point = 0
	u.PostCnt = 0
	return u, nil
}

// ComparePw 함수는 hash화된 Pw와 평문 Pw를 비교하는 함수입니다
func ComparePw(hash, pw string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	if err != nil {
		return false, err
	}
	return true, nil
}
