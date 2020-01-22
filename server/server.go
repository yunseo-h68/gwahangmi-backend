package server

import (
	"flag"
	"gwahangmi-backend/apis"
	"gwahangmi-backend/files"
	"net/http"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
)

// Server 객체는 Server에 대한 정보를 담습니다
type Server struct {
	neg    *negroni.Negroni
	router *httprouter.Router
}

// New 함수는 새로운 Server를 생성합니다
func New() (*Server, error) {
	sv := new(Server)
	sv.router = httprouter.New()
	sv.neg = negroni.Classic()

	directory := flag.String("d", ".", "the directory of static file to host")
	flag.Parse()
	fileServer := http.FileServer(FileSystem{http.Dir(*directory)})

	http.Handle("/", http.StripPrefix(strings.TrimRight("/public/", "/"), fileServer))

	for i := 0; i < len(apis.APIs); i++ {
		apis.AddAPI(sv.router, apis.APIs[i])
	}

	for i := 0; i < len(files.FileHandlers); i++ {
		files.AddFileHandler(sv.router, files.FileHandlers[i])
	}

	sv.neg.UseHandler(sv.router)
	return sv, nil
}

// Run 함수는 Server를 실행합니다
func (s *Server) Run(port string) {
	s.neg.Run(port)
}
