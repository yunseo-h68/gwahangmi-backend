package models

import (
	"net/http"

	"github.com/mholt/binding"
)

// Quiz 구조체는 퀴즈에 대한 정보를 담습니다.
type Quiz struct {
	ID             interface{} `bson:"_id" json:"id"`
	QuizID         string      `bson:"quizID" json:"quizID"`
	Author         string      `bson:"author" json:"author"`
	Title          string      `bson:"title" json:"title"`
	Explanation    string      `bson:"explanation" json:"explanation"`
	Answers        []string    `bson:"answers" json:"answers"`
	RightAnswer    string      `bson:"rightAnswer" json:"rightAnswer"`
	ParticipantCnt int         `bson:"participantCnt" json:"participantCnt"`
	Point          int         `bson:"point" json:"point"`
	UploadDate     date        `bson:"uploadDate" json:"uploadDate"`
}

// QuizParticipant 구조체는 퀴즈를 맞춘 User(퀴즈 참여자)에 대한 정보를 담습니다
type QuizParticipant struct {
	ID            interface{} `bson:"_id" json:"id"`
	UID           string      `bson:"uid" json:"uid"`
	ParentsQuizID string      `bson:"parentsQuizID" json:"parentsQuizID"`
}

// FieldMap 메서드는 Point 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (q *Quiz) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&q.ID:             "_id",
		&q.QuizID:         "quizID",
		&q.Author:         "author",
		&q.Title:          "title",
		&q.Explanation:    "explanation",
		&q.Answers:        "answers",
		&q.RightAnswer:    "rightAnswer",
		&q.ParticipantCnt: "participantCnt",
		&q.Point:          "point",
		&q.UploadDate:     "uploadDate",
	}
}

// FieldMap 메서드는 Point 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (q *QuizParticipant) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&q.ID:            "_id",
		&q.UID:           "uid",
		&q.ParentsQuizID: "parentsQuizID",
	}
}
