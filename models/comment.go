package models

import (
	"net/http"

	"github.com/mholt/binding"
)

// Comment 구조체는 댓글에 대한 정보를 담습니다.
type Comment struct {
	ID            interface{} `bson:"_id" json:"id"`
	ParentsPostID string      `bson:"parentsPostID" json:"parentsPostID"`
	CommentID     string      `bson:"commentID" json:"commentID"`
	Author        string      `bson:"author" json:"author"`
	Content       string      `bson:"content" json:"content"`
	UploadDate    string      `bson:"uploadDate" json:"uploadDate"`
}

// FieldMap 메서드는 Comment 타입을 binding.FieldMapper 인터페이스이도록 하기 위해 만든 메서드입니다.
func (c *Comment) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&c.ID:            "_id",
		&c.ParentsPostID: "parentsPostID",
		&c.CommentID:     "commentID",
		&c.Author:        "author",
		&c.Content:       "content",
		&c.UploadDate:    "uploadDate",
	}
}
