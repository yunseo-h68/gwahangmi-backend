package post

import (
	"net/http"

	"github.com/mholt/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Point 구조체는 포인트에 대한 정보를 담습니다.
type Point struct {
	ID            interface{} `bson:"_id" json:"id"`
	UID           string      `bson:"uid" json:"uid"`
	Point         int         `bson:"point" json:"point"`
	ParentsPostID string      `bson:"parentsPostID" json:"parentsPostID"`
}

// FieldMap 메서드는 Point 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (p *Point) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&p.ID:            "_id",
		&p.UID:           "uid",
		&p.Point:         "point",
		&p.ParentsPostID: "parentsPostID",
	}
}

// NewPoint 함수는 새로운 Point 구조체를 생성합니다
func NewPoint() (*Point, error) {
	p := new(Point)
	p.ID = primitive.NewObjectID()
	p.Point = 0
	p.ParentsPostID = ""
	return p, nil
}
