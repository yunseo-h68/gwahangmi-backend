package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// NewComment 함수는 새로운 Comment 구조체를 생성합니다
func NewComment() (*Comment, error) {
	c := new(Comment)
	c.ID = primitive.NewObjectID()
	c.ParentsPostID = ""
	c.CommentID = ""
	c.Author = ""
	c.Content = ""
	c.UploadDate = ""
	return c, nil
}

// NewPoint 함수는 새로운 Point 구조체를 생성합니다
func NewPoint() (*Point, error) {
	p := new(Point)
	p.ID = primitive.NewObjectID()
	p.Point = 0
	p.ParentsPostID = ""
	return p, nil
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
	p.Category = ""
	p.TotalPoint = 0
	p.AveragePoint = 0.0
	p.UploadDate = ""
	return p, nil
}

// NewUser 함수는 새로운 User 구조체를 생성합니다.
func NewUser() (*User, error) {
	u := new(User)
	u.UID = ""
	u.Pw = ""
	u.Name = ""
	u.ProfileImg = "profile_default_gwahangmi.jpg"
	u.Point = 0
	u.PostCnt = 0
	return u, nil
}

// NewQuiz 함수는 새로운 Quize 구조체를 생성합니다.
func NewQuiz() (*Quiz, error) {
	q := new(Quiz)
	q.ID = primitive.NewObjectID()
	q.QuizID = ""
	q.Author = ""
	q.Title = ""
	q.Explanation = ""
	q.Answers = []string{"", "", "", ""}
	q.RightAnswer = ""

	q.ParticipantCnt = 0
	q.Point = 0

	q.UploadDate.FullDate = ""
	q.UploadDate.Year = ""
	q.UploadDate.Month = ""
	q.UploadDate.Day = ""
	q.UploadDate.Hour = ""
	q.UploadDate.Minute = ""
	q.UploadDate.Second = ""
	return q, nil
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
