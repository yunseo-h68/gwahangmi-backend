package profile

import (
	"context"
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/files"
	"gwahangmi-backend/models"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var defaultProfile string = "profile_default_gwahangmi.jpg"

// API 구조체는 Profile Api에 대한 정보를 담습니다.
type API struct {
}

type response struct {
	ProfileImg string `json:"profileImg"`
	IsSuccess  bool   `json:"isSuccess"`
	Message    string `json:"message"`
}

// URI 메서드는 Profile API의 URI를 반환합니다
func (profileApi *API) URI() string {
	return "/api/account/profile"
}

// Get 메서드는 Profile API가 Request 메서드 중 Get을 지원함을 의미합니다
func (profileApi *API) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	uid := req.URL.Query().Get("uid")

	u := new(models.User)
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&u)

	if err != nil {
		return api.Response{http.StatusNotFound, err.Error(), response{"", false, "존재하지 않는 User"}}
	}

	return api.Response{http.StatusOK, "", response{u.ProfileImg, true, "프로필 조회 성공"}}
}

// Post 메서드는 Profile API가 Request 메서드 중 Post을 지원함을 의미합니다
func (profileApi *API) Post(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	return uploadProfile(w, req, ps)
}

// Put 메서드는 Profile API가 Request 메서드 중 Put을 지원함을 의미합니다
func (profileApi *API) Put(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	return uploadProfile(w, req, ps)
}

// Delete 메서드는 Profile API가 Request 메서드 중 Delete을 지원함을 의미합니다
func (profileApi *API) Delete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	bucket, _ := gridfs.NewBucket(
		db.MongoDB.DB("gwahangmi").DB,
	)
	uid := req.URL.Query().Get("uid")
	u := new(models.User)
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&u)
	if err != nil {
		return api.Response{http.StatusNotFound, err.Error(), response{"", false, "존재하지 않는 User"}}
	}
	if u.ProfileImg == defaultProfile {
		return api.Response{http.StatusNotFound, "", response{"", false, "프로필 이미지가 존재하지 않음"}}
	}
	return deleteProfile(uid, u.ProfileImg, bucket)
}

func uploadProfile(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	req.ParseForm()
	_, fh, err := req.FormFile("profileImg")
	if err != nil {
		return api.Response{http.StatusInternalServerError, err.Error(), response{"", false, "파일을 읽는 중 에러 발생"}}
	}
	uid := req.FormValue("uid")
	timeNow := time.Now().Format("2006-01-02-15:04:05")
	profileImgName := "profile_" + uid + "_gwahangmi_" + timeNow + filepath.Ext(fh.Filename)
	bucket, err := gridfs.NewBucket(
		db.MongoDB.DB("gwahangmi").DB,
	)

	u := new(models.User)
	err = db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&u)
	if err != nil {
		return api.Response{http.StatusNotFound, "", response{"", false, "존재하지 않는 User"}}
	}
	var img *files.ImageFile
	err = db.MongoDB.DB("gwahangmi").C("fs.files").FindOne(context.TODO(), bson.M{"filename": u.ProfileImg}).Decode(&img)
	if err == nil {
		if req.Method == "PUT" {
			log.Println("프로필 이미지 삭제 시도")
			deleteProfile(uid, u.ProfileImg, bucket)
		} else if req.Method == "POST" {
			log.Println("프로필 이미지가 이미 존재")
			return api.Response{http.StatusOK, "", response{profileImgName, false, "프로필 이미지가 이미 존재"}}
		}
	} else {
		log.Println("프로필 이미지가 없음")
	}

	_, err = db.MongoDB.DB("gwahangmi").C("users").UpdateOne(context.TODO(), bson.M{"_id": u.ID}, bson.M{"$set": bson.M{"profileImg": profileImgName}})
	if err != nil {
		return api.Response{http.StatusInternalServerError, err.Error(), response{"", false, "User 프로필이미지 이름 Update 실패"}}
	}

	file, _ := fh.Open()
	fe := filepath.Ext(fh.Filename)
	fileExt := fe[1:]
	opts := options.GridFSUpload().SetMetadata(bson.M{"uid": uid, "ext": fileExt})
	uploadStream, err := bucket.OpenUploadStream(
		profileImgName,
		opts,
	)
	if err != nil {
		return api.Response{http.StatusInternalServerError, err.Error(), response{"", false, "UploadStream Open 실패"}}
	}
	if err := files.WriteToGridFileFile(file, uploadStream); err != nil {
		return api.Response{http.StatusOK, err.Error(), response{"", false, "프로필 이미지 업로드 실패"}}
	}
	return api.Response{http.StatusCreated, "", response{profileImgName, true, "프로필 이미지 업로드 성공"}}
}

func deleteProfile(uid, profileImgName string, bucket *gridfs.Bucket) api.Response {
	if profileImgName == defaultProfile {
		return api.Response{http.StatusNotFound, "", response{"", false, "프로필 이미지가 존재하지 않음"}}
	}
	var img *files.ImageFile
	err := db.MongoDB.DB("gwahangmi").C("fs.files").FindOne(context.TODO(), bson.M{"filename": profileImgName}).Decode(&img)
	if err == nil {
		if err := bucket.Delete(img.ID); err != nil {
			log.Println("프로필 이미지 삭제 실패")
			return api.Response{http.StatusOK, err.Error(), response{"", false, "프로필 이미지 삭제 실패"}}
		}
		_, err = db.MongoDB.DB("gwahangmi").C("users").UpdateOne(context.TODO(), bson.M{"uid": uid}, bson.M{"$set": bson.M{"profileImg": defaultProfile}})
		if err != nil {
			return api.Response{http.StatusInternalServerError, err.Error(), response{"", false, "User 프로필이미지 이름 Update 실패"}}
		}
		log.Println("프로필 이미지 삭제 성공")
		return api.Response{http.StatusOK, "", response{defaultProfile, true, "프로필 이미지 삭제 성공"}}
	}
	log.Println("프로필 이미지가 존재하지 않음")
	return api.Response{http.StatusNotFound, err.Error(), response{"", false, "프로필 이미지가 존재하지 않음"}}
}
