package posts

import (
	"context"
	"errors"
	"fmt"
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/category/post"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/files"
	"log"
	"net/http"
	"strconv"
	"time"

	"gwahangmi-backend/apis/account/user"

	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// API 구조체는 Posts Api에 대한 정보를 담습니다.
type API struct {
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
func (postsApi *API) URI() string {
	return "/api/category/posts"
}

// Get 메서드는 Posts API가 Request 메서드 중 Get을 지원함을 의미합니다
func (postsApi *API) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
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
func (postsApi *API) Post(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	return uploadPost(w, req, ps)
}

// Put 메서드는 Posts API가 Request 메서드 중 Put을 지원함을 의미합니다
func (postsApi *API) Put(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	return uploadPost(w, req, ps)
}

// Delete 메서드는 Posts API가 Request 메서드 중 Delete을 지원함을 의미합니다
func (postsApi *API) Delete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	postID := req.URL.Query().Get("post_id")
	category := req.URL.Query().Get("category")
	p := new(post.Post)
	err := db.MongoDB.DB("gwahangmi").C("category_"+category).FindOne(context.TODO(), bson.M{"postID": postID}).Decode(&p)
	if err != nil {
		return api.Response{http.StatusNotFound, "", postResponse{postID, false, "해당 글이 존재하지 않음"}}
	}
	log.Println("글 삭제 시도")
	return deletePost(*p)
}

func uploadPost(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
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
	if p.PostID != "" {
		if req.Method == "PUT" {
			log.Println("Post 삭제 시도")
			deletePost(*p)
		} else if req.Method == "POST" {
			log.Println("Post가 이미 존재")
			return api.Response{http.StatusOK, "", postResponse{"", false, "해당 글ID가 이미 존재"}}
		}
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

	if err := uploadPostToGridFile(p); err != nil {
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

func uploadPostToGridFile(p *post.Post) error {
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

func deletePost(p post.Post) api.Response {
	check := user.User{}
	if res, err := idCheck(&check, p.Author); err != nil {
		return res
	}

	bucket, _ := gridfs.NewBucket(
		db.MongoDB.DB("gwahangmi").DB,
	)
	var content *post.Content
	err := db.MongoDB.DB("gwahangmi").C("fs.files").FindOne(context.TODO(), bson.M{"filename": p.PostID}).Decode(&content)
	if err == nil {
		// 글 content 삭제
		if err := bucket.Delete(content.ID); err != nil {
			log.Println("글 삭제 실패")
			return api.Response{http.StatusOK, err.Error(), postResponse{p.PostID, false, "글 삭제 실패"}}
		}
		// 유저의 포인트 감소
		_, err := db.MongoDB.DB("gwahangmi").C("users").UpdateOne(context.TODO(), bson.M{"uid": p.Author}, bson.M{"$set": bson.M{"point": check.Point - 5}})
		if err != nil {
			log.Println("포인트 적립 실패 : ", err)
			return api.Response{http.StatusInternalServerError, err.Error(), postResponse{p.PostID, false, "포인트 적립 실패"}}
		}
		// 유저의 작성글 개수 감소
		_, err = db.MongoDB.DB("gwahangmi").C("users").UpdateOne(context.TODO(), bson.M{"uid": p.Author}, bson.M{"$set": bson.M{"postCnt": check.PostCnt - 1}})
		log.Println(check.PostCnt + 1)
		if err != nil {
			log.Println("PostCnt 업데이트 실패 : ", err)
			return api.Response{http.StatusInternalServerError, err.Error(), postResponse{p.PostID, false, "PostCnt 업데이트 실패"}}
		}
		// Category에 있는 글 정보 삭제
		res, err := db.MongoDB.DB("gwahangmi").C("category_"+p.Category).DeleteOne(context.TODO(), bson.D{primitive.E{Key: "postID", Value: p.PostID}})
		if res.DeletedCount == 0 || err != nil {
			log.Println("PostInfo 삭제 실패")
			return api.Response{http.StatusInternalServerError, "", postResponse{p.PostID, false, "PostInfo 삭제 실패"}}
		}
		// posts에 있는 글 포인트 정보 삭제
		res, err = db.MongoDB.DB("gwahangmi").C("posts").DeleteOne(context.TODO(), bson.D{primitive.E{Key: "postID", Value: p.PostID}})
		if res.DeletedCount == 0 || err != nil {
			log.Println("PointPost 삭제 실패")
			return api.Response{http.StatusInternalServerError, "", postResponse{p.PostID, false, "PointPost 삭제 실패"}}
		}
		log.Println("글 삭제 성공")
		return api.Response{http.StatusOK, "", postResponse{"", true, "글 삭제 성공"}}
	}
	log.Println("해당 글이 존재하지 않음 : ", err)
	return api.Response{http.StatusOK, err.Error(), postResponse{p.PostID, false, "해당 글이 존재하지 않음"}}
}
