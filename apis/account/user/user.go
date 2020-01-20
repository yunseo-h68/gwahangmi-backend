package user

import (
	"context"
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/apis/method"
	"gwahangmi-backend/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mholt/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

// API 구조체는 Login Api에 대한 정보를 담습니다.
type API struct {
	method.PostNotSupported
	method.DeleteNotSupported
}

type user struct {
	UID        string `bson:"uid" json:"uid"`
	Name       string `bson:"uname" json:"uname"`
	ProfileImg string `bson:"profileImg" json:"profileImg"`
	Point      int    `bson:"point" json:"point"`
	PostCnt    int    `bson:"postCnt" json:"postCnt"`
}

type response struct {
	UID       string `bson:"uid" json:"uid"`
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

// URI 메서드는 User API의 URI를 반환합니다
func (userApi *API) URI() string {
	return "/api/account/users/:uid"
}

// Get 메서드는 User API가 Request 메서드 중 Get을 지원함을 의미합니다
func (userApi *API) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	uid := ps.ByName("uid")
	log.Println("UID: ", uid)

	u := user{}
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&u)
	if err != nil {
		log.Println("존재하지 않는 User:", err)
		return api.Response{http.StatusNotFound, err.Error(), nil}
	}
	return api.Response{http.StatusOK, "", u}
}

// Put 메서드는 User API가 Request 메서드 중 Put을 지원함을 의미합니다
func (userApi *API) Put(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	u, _ := models.NewUser()
	uid := ps.ByName("uid")
	log.Println("UID: ", uid)

	if errs := binding.Bind(req, u); errs != nil {
		log.Println("요청 메시지 파싱 실패 : ", errs)
		return api.Response{http.StatusInternalServerError, errs.Error(), nil}
	}

	check := models.User{}
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&check)
	if err == nil {
		_, err = db.MongoDB.DB("gwahangmi").C("users").UpdateOne(context.TODO(), bson.M{"_id": check.ID}, bson.M{"$set": bson.M{"uname": u.Name}})
		if err != nil {
			return api.Response{http.StatusInternalServerError, err.Error(), response{"", false, "User Update 실패"}}
		}
	} else {
		return api.Response{http.StatusOK, "", response{"", false, "존재하지 않는 User"}}
	}
	return api.Response{http.StatusOK, "", response{check.UID, true, "User Update 성공"}}
}

// Delete 메서드는 User API가 Request 메서드 중 Delete 지원함을 의미합니다
func (userApi *API) Delete(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	uid := ps.ByName("uid")
	check := models.User{}
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&check)
	if err == nil {
		res, err := db.MongoDB.DB("gwahangmi").C("users").DeleteOne(context.TODO(), bson.D{primitive.E{Key: "uid", Value: check.UID}})
		if res.DeletedCount == 0 || err != nil {
			log.Println("PostInfo 삭제 실패")
			return api.Response{http.StatusInternalServerError, "", response{uid, false, "User 삭제 실패"}}
		}
		if check.ProfileImg == "profile_default_gwahangmi.jpg" {
			return api.Response{http.StatusOK, "", response{"", true, "User 삭제 성공"}}
		}
		bucket, _ := gridfs.NewBucket(
			db.MongoDB.DB("gwahangmi").DB,
		)
		var img *models.ImageFile
		err = db.MongoDB.DB("gwahangmi").C("fs.files").FindOne(context.TODO(), bson.M{"filename": check.ProfileImg}).Decode(&img)
		if err == nil {
			if err := bucket.Delete(img.ID); err != nil {
				log.Println("프로필 이미지 삭제 실패")
				return api.Response{http.StatusOK, err.Error(), response{uid, false, "User 프로필 이미지 삭제 실패"}}
			}
			return api.Response{http.StatusOK, "", response{"", true, "User 삭제 성공"}}
		}
		return api.Response{http.StatusOK, err.Error(), response{uid, false, "User 프로필 이미지를 찾는 데 실패"}}
	}
	return api.Response{http.StatusOK, "", response{uid, false, "존재하지 않는 User"}}
}
