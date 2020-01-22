package point

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
)

// API 구조체는 Login Api에 대한 정보를 담습니다.
type API struct {
	method.PutNotSupported
	method.DeleteNotSupported
}

type response struct {
	IsSuccess bool   `json:"isSuccess"`
	Message   string `json:"message"`
}

type pointRes struct {
	UID           string `bson:"uid" json:"uid"`
	Point         int    `bson:"point" json:"point"`
	ParentsPostID string `bson:"parentsPostID" json:"parentsPostID"`
}

// URI 메서드는 Post API의 URI를 반환합니다
func (pointApi *API) URI() string {
	return "/api/category/posts/:postID/point"
}

// Get 메서드는 Point API가 Request 메서드 중 Get 지원함을 의미합니다
func (pointApi *API) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	postID := ps.ByName("postID")
	uid := req.URL.Query().Get("uid")

	userCheck := models.User{}
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&userCheck)
	if err != nil {
		log.Println("존재하지 않는 User")
		return api.Response{http.StatusOK, err.Error(), "존재하지 않는 User"}
	}
	pointResCheck := pointRes{}
	err = db.MongoDB.DB("gwahangmi").C(postID+"point").FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&pointResCheck)
	if err != nil {
		log.Println("아직 평가하지 않은 Post : ", err)
		return api.Response{http.StatusOK, err.Error(), "아직 평가하지 않은 Post"}
	}
	return api.Response{http.StatusOK, "", pointResCheck}
}

// Post 메서드는 Point API가 Request 메서드 중 Post 지원함을 의미합니다
func (pointApi *API) Post(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	point, _ := models.NewPoint()
	postID := ps.ByName("postID")
	if errs := binding.Bind(req, point); errs != nil {
		log.Println("요청 메시지 파싱 실패 : ", errs)
		return api.Response{http.StatusInternalServerError, errs.Error(), response{false, "요청 메시지 파싱 실패"}}
	}

	userCheck := models.User{}
	err := db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": point.UID}).Decode(&userCheck)
	if err != nil {
		log.Println("존재하지 않는 User")
		return api.Response{http.StatusOK, err.Error(), response{false, "존재하지 않는 User"}}
	}

	point.ParentsPostID = postID
	if _, err := db.MongoDB.DB("gwahangmi").C(postID+"point").InsertOne(context.TODO(), point); err != nil {
		log.Println("DB Insert 실패 : ", err)
		return api.Response{http.StatusInternalServerError, err.Error(), response{false, "DB Insert 실패"}}
	}

	pointPostCheck := new(models.PointPost)
	err = db.MongoDB.DB("gwahangmi").C("posts").FindOne(context.TODO(), bson.M{"postID": postID}).Decode(&pointPostCheck)
	if err != nil {
		log.Println("POST 찾기 실패 : ", err)
		return api.Response{http.StatusOK, err.Error(), response{false, "존재하지 않는 Post ID"}}
	}
	postCheck := new(models.Post)
	err = db.MongoDB.DB("gwahangmi").C("category_"+pointPostCheck.Category).FindOne(context.TODO(), bson.M{"postID": postID}).Decode(&postCheck)
	if err != nil {
		log.Println("POST 찾기 실패 : ", err)
		return api.Response{http.StatusOK, err.Error(), response{false, "존재하지 않는 Post ID"}}
	}

	authorCheck := models.User{}
	err = db.MongoDB.DB("gwahangmi").C("users").FindOne(context.TODO(), bson.M{"uid": postCheck.Author}).Decode(&authorCheck)
	if err != nil {
		log.Println("존재하지 않는 PostAuthor")
	} else {
		_, err = db.MongoDB.DB("gwahangmi").C("users").UpdateOne(context.TODO(), bson.M{"uid": postCheck.Author}, bson.M{"$set": bson.M{"point": authorCheck.Point + 1}})
		if err != nil {
			log.Println("포인트 적립 실패 : ", err)
			return api.Response{http.StatusInternalServerError, err.Error(), response{false, "포인트 적립 실패"}}
		}
	}

	/* 카테고리 포인트 적립 */
	_, err = db.MongoDB.DB("gwahangmi").C("category_"+pointPostCheck.Category).UpdateOne(context.TODO(), bson.M{"postID": postID}, bson.M{"$set": bson.M{"participantCnt": postCheck.ParticipantCnt + 1}})
	if err != nil {
		log.Println("참여자 추가 실패: ", err)
		return api.Response{http.StatusOK, err.Error(), response{false, "참여자 추가 실패"}}
	}
	_, err = db.MongoDB.DB("gwahangmi").C("category_"+pointPostCheck.Category).UpdateOne(context.TODO(), bson.M{"postID": postID}, bson.M{"$set": bson.M{"totalPoint": postCheck.TotalPoint + point.Point}})
	if err != nil {
		log.Println("posts 포인트 적립 실패: ", err)
		return api.Response{http.StatusOK, err.Error(), response{false, "포인트 적립 실패"}}
	}
	_, err = db.MongoDB.DB("gwahangmi").C("category_"+pointPostCheck.Category).UpdateOne(context.TODO(), bson.M{"postID": postID}, bson.M{"$set": bson.M{"averagePoint": float64(float64(postCheck.TotalPoint+point.Point) / (float64(postCheck.ParticipantCnt) + float64(1)))}})
	if err != nil {
		log.Println("posts 포인트 적립 실패: ", err)
		return api.Response{http.StatusOK, err.Error(), response{false, "포인트 적립 실패"}}
	}
	/* 포인트Post 포인트 적립 */
	_, err = db.MongoDB.DB("gwahangmi").C("posts").UpdateOne(context.TODO(), bson.M{"postID": postID}, bson.M{"$set": bson.M{"totalPoint": postCheck.TotalPoint + point.Point}})
	if err != nil {
		log.Println("posts 포인트 적립 실패: ", err)
		return api.Response{http.StatusOK, err.Error(), response{false, "포인트 적립 실패"}}
	}
	_, err = db.MongoDB.DB("gwahangmi").C("posts").UpdateOne(context.TODO(), bson.M{"postID": postID}, bson.M{"$set": bson.M{"averagePoint": float64(float64(postCheck.TotalPoint+point.Point) / (float64(postCheck.ParticipantCnt) + float64(1)))}})
	if err != nil {
		log.Println("posts 포인트 적립 실패: ", err)
		return api.Response{http.StatusOK, err.Error(), response{false, "포인트 적립 실패"}}
	}
	return api.Response{http.StatusOK, "", response{true, "성공"}}
}
