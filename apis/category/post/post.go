package post

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
	return "/api/category/posts/:post-id"
}

// Get 메서드는 Post API가 Request 메서드 중 Get을 지원함을 의미합니다
func (postApi *API) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	postID := ps.ByName("post-id")
	log.Println("PostID: ", postID)

	var pointPost models.PointPost
	err := db.MongoDB.DB("gwahangmi").C("posts").FindOne(context.TODO(), bson.M{"postID": postID}).Decode(&pointPost)
	if err != nil {
		log.Println("P 존재하지 않는 Post :", err)
		return api.Response{http.StatusOK, err.Error(), nil}
	}

	var postData models.Post
	err = db.MongoDB.DB("gwahangmi").C("category_"+pointPost.Category).FindOne(context.TODO(), bson.M{"postID": postID}).Decode(&postData)
	if err != nil {
		log.Println("C 존재하지 않는 Post :", err)
		return api.Response{http.StatusOK, err.Error(), nil}
	}
	return api.Response{http.StatusOK, "", postData}
}
