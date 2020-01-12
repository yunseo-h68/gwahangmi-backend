package account

import (
	"gwahangmi-backend/apis/account/login"
	"gwahangmi-backend/apis/account/profile"
	"gwahangmi-backend/apis/account/signup"
	"gwahangmi-backend/apis/account/user"
	"gwahangmi-backend/apis/account/users"
	"gwahangmi-backend/apis/api"
)

// AccountAPIs 는 Account에 대한 API 리스트입니다
var AccountAPIs []api.API

func init() {
	AccountAPIs = make([]api.API, 0)

	apis := []api.API{
		new(login.API),
		new(signup.API),
		new(profile.API),
		new(users.API),
		new(user.API),
	}

	for i := 0; i < len(apis); i++ {
		AccountAPIs = append(AccountAPIs, apis[i])
	}
}
