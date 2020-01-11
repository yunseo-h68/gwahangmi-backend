package posts

import (
	"context"
	"errors"
	"fmt"
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/category/post"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/apis/method"
	"gwahangmi-backend/files"
	"log"
	"net/http"
	"strconv"
	"time"

	"gwahangmi-backend/apis/account/user"

	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// API 구조체는 Posts Api에 대한 정보를 담습니다.
type API struct {
	method.PutNotSupported
	method.DeleteNotSupported
}

type getResponse struct {
	Posts []string `json:"posts"`
}

type postResponse struct {
	PostID    string `json:"postID"`
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

// URI 메서드는 Posts API의 URI를 반환합니다
func (signupApi *API) URI() string {
	return "/api/category/posts"
}

// Get 메서드는 Posts API가 Request 메서드 중 Get을 지원함을 의미합니다
func (signupApi *API) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	findOptions := options.Find()
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	findOptions.SetLimit(int64(limit))
	var results []*post.PointPost
	cur, err := db.MongoDB.DB("gwahangmi").C("posts").Find(context.TODO(), bson.D{{}}, findOptions)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("completed find")

	for cur.Next(context.TODO()) {
		var elem post.PointPost
		err := cur.Decode(&elem)
		fmt.Printf(" document: %+v\n", elem)
		if err != nil {
			log.Println(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}

	var posts []string
	for i := 0; i < len(results); i++ {
		posts = append(posts, results[i].PostID)
	}
	return api.Response{http.StatusOK, "", getResponse{posts}}
}

// Post 메서드는 Posts API가 Request 메서드 중 Post을 지원함을 의미합니다
func (signupApi *API) Post(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	p, _ := post.NewPost()

	if errs := binding.Bind(req, p); errs != nil {
		log.Println("요청 메시지 파싱 실패 : ", errs)
		return api.Response{http.StatusInternalServerError, errs.Error(), postResponse{"", false, "요청 메시지 파싱 실패"}}
	}
	log.Println(p)
	check := user.User{}
	if checkRes, err := idCheck(&check, p.Author); err != nil {
		log.Println(err)
		return checkRes
	}
	timeNow := time.Now()
	p.PostID = "post_" + check.UID + "_gwahangmi_" + timeNow.Format("2006-01-02-15:04:05")
	p.UploadDate.FullDate = timeNow.Format("2006-01-02-15:04:05")
	p.UploadDate.Year = timeNow.Year()
	p.UploadDate.Month = timeNow.Month()
	p.UploadDate.Day = timeNow.Day()
	p.UploadDate.Hour = timeNow.Hour()
	p.UploadDate.Minute = timeNow.Minute()
	p.UploadDate.Second = timeNow.Second()

	if err := uploadPost(p); err != nil {
		log.Println(err)
		return api.Response{http.StatusInternalServerError, err.Error(), postResponse{"", false, "글을 DB에 저장하는 데 실패"}}
	}
	p.Content = p.PostID
	if _, err := db.MongoDB.DB("gwahangmi").C("category_"+p.Category).InsertOne(context.TODO(), p); err != nil {
		log.Println(err)
		return api.Response{http.StatusInternalServerError, err.Error(), postResponse{"", false, "글을 DB에 저장하는 데 실패"}}
	}
	ppoint, _ := post.NewPointPost()
	ppoint.PostID = p.PostID
	ppoint.TotalPoint = p.TotalPoint
	if _, err := db.MongoDB.DB("gwahangmi").C("posts").InsertOne(context.TODO(), ppoint); err != nil {
		log.Println(err)
		return api.Response{http.StatusInternalServerError, err.Error(), postResponse{"", false, "글을 DB에 저장하는 데 실패"}}
	}

	_, err := db.MongoDB.DB("gwahangmi").C("users").UpdateOne(context.TODO(), bson.M{"uid": p.Author}, bson.M{"$set": bson.M{"point": check.Point + 5}})
	if err != nil {
		log.Println("포인트 적립 실패 : ", err)
		return api.Response{http.StatusInternalServerError, err.Error(), postResponse{"", false, "포인트 적립 실패"}}
	}
	_, err = db.MongoDB.DB("gwahangmi").C("users").UpdateOne(context.TODO(), bson.M{"uid": p.Author}, bson.M{"$set": bson.M{"postCnt": check.PostCnt + 1}})
	log.Println(check.PostCnt + 1)
	if err != nil {
		log.Println("PostCnt 업데이트 실패 : ", err)
		return api.Response{http.StatusInternalServerError, err.Error(), postResponse{"", false, "PostCnt 업데이트 실패"}}
	}

	return api.Response{http.StatusOK, "", postResponse{p.PostID, true, "Post 업로드 성공"}}
}

func uploadPost(p *post.Post) error {
	bucket, err := gridfs.NewBucket(
		db.MongoDB.DB("gwahangmi").DB,
	)
	if err != nil {
		return err
	}

	opts := options.GridFSUpload().SetMetadata(bson.M{"uid": p.Author})
	uploadStream, err := bucket.OpenUploadStream(
		p.PostID,
		opts,
	)
	if err := files.WriteToGridFileString(p.Content, uploadStream); err != nil {
		return err
	}
	return nil
}

func idCheck(check *user.User, uid string) (api.Response, error) {
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&check)
	if err != nil {
		return api.Response{http.StatusNotFound, "", postResponse{"", false, "존재하지 않는 User의 접근"}}, errors.New("존재하지 않는 User")
	}
	return api.Response{}, nil
}
