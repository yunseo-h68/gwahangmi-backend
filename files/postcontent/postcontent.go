package postcontent

import (
	"bytes"
	"context"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// FileHandler 구조체는 ProfileImg File Handler에 대한 정보를 담습니다.
type FileHandler struct {
}

// URI 메서드는 Posts API의 URI를 반환합니다
func (postContentHandler *FileHandler) URI() string {
	return "/api/file/post/content/:postID"
}

// Handler 메서드는 ProfileImgFIleHandler의 핸들러 함수입니다.
func (postContentHandler *FileHandler) Handler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	postID := ps.ByName("postID")
	var postContentFile models.PostContent
	err := db.MongoDB.DB("gwahangmi").C("fs.files").FindOne(context.TODO(), bson.M{"filename": postID}).Decode(&postContentFile)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	bucket, _ := gridfs.NewBucket(db.MongoDB.DB("gwahangmi").DB)
	var buf bytes.Buffer
	_, err = bucket.DownloadToStreamByName(
		postContentFile.Filename,
		&buf,
	)
	if err != nil {
		log.Println("Failed to open DownloadStream")
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(buf.String()))
}
