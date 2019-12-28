package response

import (
	"log"
	"net/http"
)

// Response 구조체는 응답값을 정의합니다
type Response interface {
	Code() int
	Data() ([]byte, error)
}

// HTTPResponse 함수는 Response 구조체를 HTTP로 응답합니다
func HTTPResponse(w http.ResponseWriter, req *http.Request, res Response) {
	content, err := res.Data()

	log.Println(string(content))
	if err != nil {
		abort(w, 500)
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(res.Code())
	w.Write(content)
}

func abort(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
}
