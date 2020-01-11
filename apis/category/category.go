package category

import (
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/category/post"
	"gwahangmi-backend/apis/category/posts"
)

// CategoryAPIs 는 Category, Post에 대한 API 리스트입니다
var CategoryAPIs []api.API

func init() {
	CategoryAPIs = make([]api.API, 0)

	apis := []api.API{
		new(posts.API),
		new(post.API),
	}

	for i := 0; i < len(apis); i++ {
		CategoryAPIs = append(CategoryAPIs, apis[i])
	}
}
