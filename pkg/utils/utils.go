package utils

import (
	historyModel "washit-api/internal/history/dto/model"
	orderModel "washit-api/internal/order/dto/model"
	userModel "washit-api/internal/user/dto/model"
)

var ModelList = []interface{}{&userModel.User{}, &orderModel.Order{}, &historyModel.History{}}
