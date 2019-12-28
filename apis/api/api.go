package api

import (
	"gwahangmi-backend/apis/response"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// API 인터페이스는 API에서 구현해야할 메서드를 가집니다
type API interface {
	URI() string
	Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) response.Response
	Post(w http.ResponseWriter, req *http.Request, ps httprouter.Params) response.Response
	Put(w http.ResponseWriter, req *http.Request, ps httprouter.Params) response.Response
	Delete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) response.Response
}
