package users

import (
	"context"
	"fmt"
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/db"
	"gwahangmi-backend/apis/method"
	"gwahangmi-backend/models"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// API 구조체는 Login Api에 대한 정보를 담습니다.
type API struct {
	method.PostNotSupported
	method.PutNotSupported
	method.DeleteNotSupported
}

type response struct {
	Users []string `json:"users"`
}

// URI 메서드는 Post API의 URI를 반환합니다
func (postApi *API) URI() string {
	return "/api/account/users"
}

// Get 메서드는 Post API가 Request 메서드 중 Get을 지원함을 의미합니다
func (postApi *API) Get(w http.ResponseWriter, req *http.Request, ps httprouter.Params) api.Response {
	findOptions := options.Find()

	point, _ := strconv.ParseBool(req.URL.Query().Get("point"))
	post, _ := strconv.ParseBool(req.URL.Query().Get("post"))
	sort, _ := strconv.ParseBool(req.URL.Query().Get("sort"))

	if point {
		if sort {
			findOptions.SetSort(bson.D{primitive.E{Key: "point", Value: 1}})
		} else {
			findOptions.SetSort(bson.D{primitive.E{Key: "point", Value: -1}})
		}
	} else if post {
		if sort {
			findOptions.SetSort(bson.D{primitive.E{Key: "postCnt", Value: 1}})
		} else {
			findOptions.SetSort(bson.D{primitive.E{Key: "postCnt", Value: -1}})
		}
	}
	limit, _ := strconv.Atoi(req.URL.Query().Get("limit"))
	findOptions.SetLimit(int64(limit))

	var users []*models.User
	cur, err := db.MongoDB.DB("gwahangmi").C("users").Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Println("Find Err: ", err)
		return api.Response{http.StatusInternalServerError, err.Error(), response{nil}}
	}
	for cur.Next(context.TODO()) {
		var elem models.User
		err := cur.Decode(&elem)
		fmt.Printf(" document: %+v\n", elem)
		if err != nil {
			log.Println(err)
		}
		users = append(users, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Println(err)
	}

	var uids []string
	for i := 0; i < len(users); i++ {
		uids = append(uids, users[i].UID)
	}
	return api.Response{http.StatusOK, "", response{uids}}
}
