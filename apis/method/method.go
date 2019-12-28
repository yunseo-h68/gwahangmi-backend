package method

import (
	"encoding/json"
	"gwahangmi-backend/apis/response"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Response 는 해당 Request 메서드를 지원하지 않을 때 응답값을 정의합니다
type Response struct {
	code    int
	Message string `json:"message"`
}

// Code 메서드는 해당 Request 메서드를 지원하지 않을 때 응답할 HTTP 상태코드를 반환합니다
func (res Response) Code() int {
	return res.code
}

// Data 메서드는 해당 Request 메서드를 지원하지 않을 때 응답할 Json Data를 반환합니다
func (res Response) Data() ([]byte, error) {
	return json.Marshal(res)
}

type (
	// GetNotSupported 는 해당 Request 메서드를 지원하지 않음을 의미합니다
	GetNotSupported struct{}
	// PostNotSupported 는 해당 Request 메서드를 지원하지 않음을 의미합니다
	PostNotSupported struct{}
	// PutNotSupported 는 해당 Request 메서드를 지원하지 않음을 의미합니다
	PutNotSupported struct{}
	// DeleteNotSupported 는 해당 Request 메서드를 지원하지 않음을 의미합니다
	DeleteNotSupported struct{}
)

// Get 메서드는 해당 Request 메서드를 지원하지 않음을 의미합니다
func (GetNotSupported) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) response.Response {
	return Response{405, "Method is not supported"}
}

// Post 메서드는 해당 Request 메서드를 지원하지 않음을 의미합니다
func (PostNotSupported) Post(w http.ResponseWriter, req *http.Request, ps httprouter.Params) response.Response {
	return Response{405, "Method is not supported"}
}

// Put 메서드는 해당 Request 메서드를 지원하지 않음을 의미합니다
func (PutNotSupported) Put(w http.ResponseWriter, req *http.Request, ps httprouter.Params) response.Response {
	return Response{405, "Method is not supported"}
}

// Delete 메서드는 해당 Request 메서드를 지원하지 않음을 의미합니다
func (DeleteNotSupported) Delete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) response.Response {
	return Response{405, "Method is not supported"}
}
