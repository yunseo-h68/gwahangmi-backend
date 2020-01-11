package user

import (
	"context"
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/apis/method"
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

type user struct {
	UID        string `bson:"uid" json:"uid"`
	Name       string `bson:"uname" json:"uname"`
	ProfileImg string `bson:"profileImg" json:"profileImg"`
	Point      int    `bson:"point" json:"point"`
	PostCnt    int    `bson:"postCnt" json:"postCnt"`
}

// URI 메서드는 Post API의 URI를 반환합니다
func (postApi *API) URI() string {
	return "/api/account/user/:uid"
}

// Get 메서드는 Post API가 Request 메서드 중 Get을 지원함을 의미합니다
func (postApi *API) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	uid := ps.ByName("uid")
	log.Println("UID: ", uid)

	u := user{}
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&u)
	if err != nil {
		log.Println("존재하지 않는 User:", err)
		return api.Response{http.StatusNotFound, err.Error(), nil}
	}
	return api.Response{http.StatusOK, "", u}
}
