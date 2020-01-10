package signup

import (
	"context"
	"gwahangmi-backend/apis/account/user"
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/apis/method"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
)

// API 구조체는 Signup Api에 대한 정보를 담습니다.
type API struct {
	method.GetNotSupported
	method.PutNotSupported
	method.DeleteNotSupported
}

type response struct {
	Uname     string `json:"uname"`
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

// URI 메서드는 Signup API의 URI를 반환합니다
func (signupApi *API) URI() string {
	return "/api/account/signup"
}

// Post 메서드는 Signup API가 Request 메서드 중 Post을 지원함을 의미합니다
func (signupApi *API) Post(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	u, _ := user.New()

	if errs := binding.Bind(req, u); errs != nil {
		log.Println("요청 메시지 파싱 실패 : ", errs)
		return api.Response{http.StatusInternalServerError, errs.Error(), response{"", false, "요청 메시지 파싱에 실패하였습니다"}}
	}

	u.ID = primitive.NewObjectID()
	check := user.User{}
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": u.UID}).Decode(&check)

	if err != nil {
		hashedPw, _ := bcrypt.GenerateFromPassword([]byte(u.Pw), bcrypt.DefaultCost)
		u.Pw = string(hashedPw[:])
		if _, err := db.MongoDB.DB("gwahangmi").C("users").InsertOne(context.TODO(), u); err != nil {
			log.Println("DB Insert 실패 : ", err)
			return api.Response{http.StatusInternalServerError, err.Error(), response{"", false, "DB Insert 실패"}}
		}
		return api.Response{http.StatusCreated, "", response{u.Name, true, "Signup 성공"}}
	}
	return api.Response{http.StatusOK, "", response{"", false, "계정이 이미 존재합니다"}}
}
