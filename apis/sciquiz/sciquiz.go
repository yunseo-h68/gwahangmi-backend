package sciquiz

import (
	"gwahangmi-backend/apis/api"
	"gwahangmi-backend/apis/sciquiz/quiz"
	"gwahangmi-backend/apis/sciquiz/quizzes"
)

// SciQuizAPIs 는 Quiz에 대한 API 리스트입니다
var SciQuizAPIs []api.API

func init() {
	SciQuizAPIs = make([]api.API, 0)

	apis := []api.API{
		new(quiz.API),
		new(quizzes.API),
	}

	for i := 0; i < len(apis); i++ {
		SciQuizAPIs = append(SciQuizAPIs, apis[i])
	}
}
