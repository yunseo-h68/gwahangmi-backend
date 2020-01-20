package profileimg

import (
	"bytes"
	"context"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/models"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// FileHandler 구조체는 ProfileImg File Handler에 대한 정보를 담습니다.
type FileHandler struct {
}

// URI 메서드는 Posts API의 URI를 반환합니다
func (profileImgFileHandler *FileHandler) URI() string {
	return "/api/file/profileimg/:id"
}

// Handler 메서드는 ProfileImgFIleHandler의 핸들러 함수입니다.
func (profileImgFileHandler *FileHandler) Handler(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uid := ps.ByName("id")

	u := new(models.User)
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&u)

	name := u.ProfileImg
	var img *models.ImageFile
	err = db.MongoDB.DB("gwahangmi").C("fs.files").FindOne(context.TODO(), bson.M{"filename": name}).Decode(&img)
	if err != nil {
		if name == " " {
			log.Printf("Failed to open %s: %v", name, err)
		} else {
			log.Printf("ProfileImage is emty")
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	bucket, _ := gridfs.NewBucket(
		db.MongoDB.DB("gwahangmi").DB,
	)
	var buf bytes.Buffer
	_, err = bucket.DownloadToStreamByName(
		img.Filename,
		&buf,
	)
	if err != nil {
		log.Println("GetProfileImg : Failed to open DownloadStream")
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
	if err != nil {
		log.Printf("Failed to read downloadStream %s", err)
		http.Error(w, "Failed to read downloadStream", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, req, name, time.Now(), bytes.NewReader(buf.Bytes()))

}
