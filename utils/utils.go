package utils

import (
	historyModel "washit-api/app/history/dto/model"
	orderModel "washit-api/app/order/dto/model"
	userModel "washit-api/app/user/dto/model"
)

var ModelList = []interface{}{&userModel.User{}, &orderModel.Order{}, &historyModel.History{}}
