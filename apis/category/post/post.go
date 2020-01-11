package post

import (
	"net/http"

	"github.com/mholt/binding"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Post 구조체는 포스트에 대한 정보를 담습니다.
type Post struct {
	ID             interface{} `bson:"_id" json:"id"`
	PostID         string      `bson:"postID" json:"postID"`
	Author         string      `bson:"author" json:"author"`
	Category       string      `bson:"category" json:"category"`
	Title          string      `bson:"title" json:"title"`
	Content        string      `bson:"content" json:"content"`
	ParticipantCnt int         `bson:"participantCnt" json:"participantCnt"`
	TotalPoint     int         `bson:"totalPoint" json:"totalPoint"`
	AveragePoint   float64     `bson:"averagePoint" json:"averagePoint"`
	UploadDate     date        `bson:"uploadDate" json:"uploadDate"`
}

type date struct {
	Year     interface{} `bson:"year" json:"year"`
	Month    interface{} `bson:"month" json:"month"`
	Day      interface{} `bson:"day" json:"day"`
	Hour     interface{} `bson:"hour" json:"hour"`
	Minute   interface{} `bson:"minute" json:"minute"`
	Second   interface{} `bson:"second" json:"second"`
	FullDate interface{} `bson:"fullDate" json:"fullDate"`
}

// Content 는 Post의 content에 대한 정보를 담습니다.
type Content struct {
	ID       primitive.ObjectID `bson:"_id"`
	Filename string             `bson:"filename"`
	MetaData fileMeta           `bson:"metadata"`
}

// FileMeta 는 Upload할 파일의 메타정보를 담는 구조체입니다.
type fileMeta struct {
	Inode int
	UID   string `bson:"uid" json:"uid"`
}

// PointPost 구조체는 포스트의 포인트에 대한 정보를 담습니다.
type PointPost struct {
	ID         interface{} `bson:"_id" json:"id"`
	PostID     string      `bson:"postID" json:"postID"`
	TotalPoint int         `bson:"totalPoint" json:"totalPoint"`
}

// FieldMap 메서드는 Post 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (p *Post) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&p.ID:             "_id",
		&p.PostID:         "postID",
		&p.Author:         "author",
		&p.Category:       "category",
		&p.Title:          "title",
		&p.Content:        "content",
		&p.UploadDate:     "uploadDate",
		&p.TotalPoint:     "totalPoint",
		&p.AveragePoint:   "averagePoint",
		&p.ParticipantCnt: "participantCnt",
	}
}

// FieldMap 메서드는 PointPost 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (p *PointPost) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&p.ID:         "_id",
		&p.PostID:     "postID",
		&p.TotalPoint: "totalPoint",
	}
}

// NewPost 함수는 새로운 Post 구조체를 생성합니다
func NewPost() (*Post, error) {
	p := new(Post)
	p.ID = primitive.NewObjectID()
	p.PostID = ""
	p.Author = ""
	p.Category = ""
	p.Title = ""
	p.Content = ""

	p.UploadDate.FullDate = ""
	p.UploadDate.Year = ""
	p.UploadDate.Month = ""
	p.UploadDate.Day = ""
	p.UploadDate.Hour = ""
	p.UploadDate.Minute = ""
	p.UploadDate.Second = ""

	p.TotalPoint = 0
	p.AveragePoint = 0.0
	p.ParticipantCnt = 0
	return p, nil
}

// NewPointPost 함수는 새로운 PointPost 구조체를 생성합니다
func NewPointPost() (*PointPost, error) {
	p := new(PointPost)
	p.ID = primitive.NewObjectID()
	p.PostID = ""
	p.TotalPoint = 0
	return p, nil
}
