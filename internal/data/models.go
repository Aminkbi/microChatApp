package data

import (
	"github.com/aminkbi/microChatApp/api/utils"
)

type Models struct {
	Users UserModel
}

func NewModels(mongo *utils.MongoClient) Models {
	return Models{
		Users: UserModel{mongo: mongo},
	}
}
