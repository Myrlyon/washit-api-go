package utils

import (
	"strconv"
	historyModel "washit-api/internal/history/dto/model"
	orderModel "washit-api/internal/order/dto/model"
	userModel "washit-api/internal/user/dto/model"
)

var ModelList = []interface{}{&userModel.User{}, &orderModel.Order{}, &historyModel.History{}}

func StringToInt64(s string) (int64, error) {
    i, err := strconv.ParseInt(s, 10, 64)
    if err != nil {
        return 0, err
    }
    return i, nil
}
