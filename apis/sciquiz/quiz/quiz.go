package quiz

import (
	"context"
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/apis/method"
	"gwahangmi-backend/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
)

// API 구조체는 Login Api에 대한 정보를 담습니다.
type API struct {
	method.PostNotSupported
	method.PutNotSupported
	method.DeleteNotSupported
}

// URI 메서드는 Post API의 URI를 반환합니다
func (postApi *API) URI() string {
	return "/api/sci-quiz/quizzes/:quizID"
}

// Get 메서드는 Post API가 Request 메서드 중 Get을 지원함을 의미합니다
func (postApi *API) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	quizID := ps.ByName("quizID")
	log.Println("QuizID: ", quizID)

	var quiz models.Quiz
	err := db.MongoDB.DB("gwahangmi").C("quizzes").FindOne(context.TODO(), bson.M{"quizID": quizID}).Decode(&quiz)
	if err != nil {
		log.Println("P 존재하지 않는 Quiz :", err)
		return api.Response{http.StatusOK, err.Error(), nil}
	}

	return api.Response{http.StatusOK, "", quiz}
}
