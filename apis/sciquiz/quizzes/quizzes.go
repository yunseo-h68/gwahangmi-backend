package quizzes

import (
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/apis/method"
	"gwahangmi-backend/files"
	"gwahangmi-backend/models"

	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// API 구조체는 Posts Api에 대한 정보를 담습니다.
type API struct {
	method.PutNotSupported
	method.DeleteNotSupported
}

type getResponse struct {
	Quizzes []string `json:"quizzes"`
}

type postResponse struct {
	QuizID    string `json:"quizID"`
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

// URI 메서드는 Posts API의 URI를 반환합니다
func (postsApi *API) URI() string {
	return "/api/sci-quiz/quizzes"
}

// Get 메서드는 Posts API가 Request 메서드 중 Get을 지원함을 의미합니다
func (postsApi *API) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	findOptions := options.Find()

	point, _ := strconv.ParseBool(req.URL.Query().Get("point"))
	participantCnt, _ := strconv.ParseBool(req.URL.Query().Get("participantCnt"))
	sort, _ := strconv.ParseBool(req.URL.Query().Get("sort"))
	if popularity, _ := strconv.ParseBool(req.URL.Query().Get("popularity")); popularity {
		if point {
			if sort {
				findOptions.SetSort(bson.D{primitive.E{Key: "point", Value: 1}})
			} else {
				findOptions.SetSort(bson.D{primitive.E{Key: "point", Value: -1}})
			}
		} else if participantCnt {
			if sort {
				findOptions.SetSort(bson.D{primitive.E{Key: "participantCnt", Value: 1}})
			} else {
				findOptions.SetSort(bson.D{primitive.E{Key: "participantCnt", Value: -1}})
			}
		}
	} else {
		if sort {
			findOptions.SetSort(bson.D{primitive.E{Key: "uploadDate", Value: 1}})
		} else {
			findOptions.SetSort(bson.D{primitive.E{Key: "uploadDate", Value: -1}})
		}
	}
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	skip, _ := strconv.Atoi(req.URL.Query().Get("skip"))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(skip))

	cur, err := db.MongoDB.DB("gwahangmi").C("quizzes").Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Println("Find Err: ", err)
		return api.Response{http.StatusInternalServerError, err.Error(), getResponse{nil}}
	}

	var quizzes []string
	for cur.Next(context.TODO()) {
		var elem models.Quiz
		err := cur.Decode(&elem)
		fmt.Printf(" document: %+v\n", elem)
		if err != nil {
			log.Println(err)
		}
		quizzes = append(quizzes, elem.QuizID)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}

	return api.Response{http.StatusOK, "", getResponse{quizzes}}
}

// Post 메서드는 Posts API가 Request 메서드 중 Post을 지원함을 의미합니다
func (postsApi *API) Post(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	return createQuiz(w, req, ps)
}

func createQuiz(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	q, _ := models.NewQuiz()

	if errs := binding.Bind(req, q); errs != nil {
		log.Println("요청 메시지 파싱 실패 : ", errs)
		return api.Response{http.StatusInternalServerError, errs.Error(), postResponse{"", false, "요청 메시지 파싱 실패"}}
	}
	log.Println(q)
	check := models.User{}
	if checkRes, err := idCheck(&check, q.Author); err != nil {
		log.Println(err)
		return checkRes
	}

	timeNow := time.Now()
	q.QuizID = "quiz_" + check.UID + "_gwahangmi_" + timeNow.Format("2006_01_02_15_04_05")
	q.UploadDate.FullDate = timeNow.Format("2006-01-02-15:04:05")
	q.UploadDate.Year = timeNow.Year()
	q.UploadDate.Month = timeNow.Month()
	q.UploadDate.Day = timeNow.Day()
	q.UploadDate.Hour = timeNow.Hour()
	q.UploadDate.Minute = timeNow.Minute()
	q.UploadDate.Second = timeNow.Second()

	if _, err := db.MongoDB.DB("gwahangmi").C("quizzes").InsertOne(context.TODO(), q); err != nil {
		log.Println(err)
		return api.Response{http.StatusInternalServerError, err.Error(), postResponse{"", false, "퀴즈를 DB에 저장하는 데 실패"}}
	}

	_, err := db.MongoDB.DB("gwahangmi").C("users").UpdateOne(context.TODO(), bson.M{"uid": q.Author}, bson.M{"$set": bson.M{"point": check.Point + 5}})
	if err != nil {
		log.Println("포인트 적립 실패 : ", err)
		return api.Response{http.StatusInternalServerError, err.Error(), postResponse{"", false, "포인트 적립 실패"}}
	}

	return api.Response{http.StatusOK, "", postResponse{q.QuizID, true, "Quiz 생성 성공"}}
}

func uploadPostToGridFile(p *models.Post) error {
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

func idCheck(check *models.User, uid string) (api.Response, error) {
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&check)
	if err != nil {
		return api.Response{http.StatusOK, "", postResponse{"", false, "존재하지 않는 User의 접근"}}, errors.New("존재하지 않는 User")
	}
	return api.Response{}, nil
}
