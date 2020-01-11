package login

import (
	"context"
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/apis/method"
	"gwahangmi-backend/models"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
)

// API 구조체는 Login Api에 대한 정보를 담습니다.
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

// URI 메서드는 Login API의 URI를 반환합니다
func (loginApi *API) URI() string {
	return "/api/account/login"
}

// Post 메서드는 Login API가 Request 메서드 중 Post을 지원함을 의미합니다
func (loginApi *API) Post(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	ul := new(models.Login)

	if errs := binding.Bind(req, ul); errs != nil {
		log.Println("요청메시지 파싱 실패 : ", errs)
		return api.Response{http.StatusInternalServerError, "", response{"", false, "요청메시지 파싱에 실패하였습니다"}}
	}

	u := new(models.User)
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": ul.UID}).Decode(&u)

	if err != nil {
		log.Println("존재하지 않는 User : ", err)
		return api.Response{http.StatusNotFound, "", response{"", false, "존재하지 않는 User입니다"}}
	}

	if u.UID == ul.UID {
		pwOK, err := comparePw(u.Pw, ul.Pw)
		if pwOK {
			log.Println("로그인 성공")
			return api.Response{http.StatusOK, "", response{u.Name, true, "로그인에 성공하셨습니다"}}
		}
		log.Println("잘못된 PW : ", err)
		return api.Response{http.StatusOK, "", response{"", false, "잘못된 PW입니다"}}
	}
	log.Println("잘못된 ID")
	return api.Response{http.StatusOK, "", response{"", false, "잘못된 ID입니다"}}
}

func comparePw(hash, pw string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
	if err != nil {
		return false, err
	}
	return true, nil
}
