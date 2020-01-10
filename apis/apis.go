package apis

import (
	"gwahangmi-backend/apis/account"
	"gwahangmi-backend/apis/api"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// APIs 변수는 등록된 API들을 담습니다
var APIs []api.API

func init() {
	APIs = make([]api.API, 0)

	apis := [][]api.API{
		account.AccountAPIs,
	}

	for i := 0; i < len(apis); i++ {
		for j := 0; j < len(apis[i]); j++ {
			APIs = append(APIs, apis[i][j])
		}
	}
}

// AddAPI 함수는 API를 등록합니다
func AddAPI(router *httprouter.Router, apiElement api.API) {
	log.Println("\"" + apiElement.URI() + "\" api is registerd")

	router.GET(apiElement.URI(), func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		res := apiElement.Get(w, req, ps)
		api.HTTPResponse(w, req, res)
	})
	router.POST(apiElement.URI(), func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		res := apiElement.Post(w, req, ps)
		api.HTTPResponse(w, req, res)
	})
	router.PUT(apiElement.URI(), func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		res := apiElement.Put(w, req, ps)
		api.HTTPResponse(w, req, res)
	})
	router.DELETE(apiElement.URI(), func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		res := apiElement.Delete(w, req, ps)
		api.HTTPResponse(w, req, res)
	})
}
