package login

import (
	"context"
	"encoding/json"
	"gwahangmi-backend/apis/account/user"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/apis/method"
	"gwahangmi-backend/apis/response"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
)

// API 구조체는 Login Api에 대한 정보를 담습니다.
type API struct {
	method.GetNotSupported
	method.PutNotSupported
	method.DeleteNotSupported
}

// Response 는 Login API의 응답값을 정의합니다
type Response struct {
	code      int
	Uname     string `json:"uname"`
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

// Code 메서드는 Login API의 응답 HTTP 상태코드를 반환합니다
func (res Response) Code() int {
	return res.code
}

// Data 메서드는 Login API의 Json Data를 반환합니다
func (res Response) Data() ([]byte, error) {
	return json.Marshal(res)
}

// URI 메서드는 Login API의 URI를 반환합니다
func (api *API) URI() string {
	return "/api/account/login"
}

// Post 메서드는 Login API가 Request 메서드 중 Post을 지원함을 의미합니다
func (api *API) Post(w http.ResponseWriter, req *http.Request, ps httprouter.Params) response.Response {
	ul := new(user.Login)

	if errs := binding.Bind(req, ul); errs != nil {
		log.Println(errs)
		return &Response{http.StatusInternalServerError, "", false, "요청 메시지 파싱에 실패하였습니다"}
	}

	u := new(user.User)
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": ul.UID}).Decode(&u)

	if err != nil {
		log.Println(err)
		return &Response{http.StatusNotFound, "", false, "존재하지 않는 User입니다"}
	}

	if u.UID == ul.UID {
		pwOK, err := user.ComparePw(u.Pw, ul.Pw)
		if pwOK {
			return &Response{http.StatusOK, u.Name, true, "로그인에 성공하셨습니다"}
		}
		log.Println(err)
		return &Response{http.StatusOK, "", false, "잘못된 PW입니다"}
	}
	return &Response{http.StatusOK, "", false, "잘못된 ID입니다"}
}
