package apis

import (
	"gwahangmi-backend/apis/account"
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/response"
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
func AddAPI(router *httprouter.Router, api api.API) {
	log.Println("\"" + api.URI() + "\" api is registerd")

	router.GET(api.URI(), func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		res := api.Get(w, req, ps)
		response.HTTPResponse(w, req, res)
	})
	router.POST(api.URI(), func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		res := api.Post(w, req, ps)
		response.HTTPResponse(w, req, res)
	})
	router.PUT(api.URI(), func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		res := api.Put(w, req, ps)
		response.HTTPResponse(w, req, res)
	})
	router.DELETE(api.URI(), func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		res := api.Delete(w, req, ps)
		response.HTTPResponse(w, req, res)
	})
}
