package models

import (
	"net/http"

	"github.com/mholt/binding"
)

// QuizPass 구조체는 퀴즈에 참여한 User가 제출한 답안에 대한 정보를 담습니다.
type QuizPass struct {
	ID            interface{} `bson:"_id" json:"id"`
	UID           string      `bson:"uid" json:"uid"`
	Pass          bool        `bson:"pass" json:"pass"`
	ParentsQuizID string      `bson:"parentsQuizID" json:"parentsQuizID"`
}

// FieldMap 메서드는 QuizPass 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (qp *QuizPass) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&qp.ID:            "_id",
		&qp.UID:           "uid",
		&qp.Pass:          "pass",
		&qp.ParentsQuizID: "parentsQuizID",
	}
}
